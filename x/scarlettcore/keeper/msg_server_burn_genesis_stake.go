package keeper

import (
	"context"
	"fmt"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	"scarlett-core/x/scarlettcore/types"
)

func (k msgServer) BurnGenesisStake(goCtx context.Context, msg *types.MsgBurnGenesisStake) (*types.MsgBurnGenesisStakeResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	ctx.Logger().Info("DEBUG: Starting BurnGenesisStake handler", "amount", msg.Amount)

	// 1. Parse and validate the burn amount
	burnCoin, err := sdk.ParseCoinNormalized(msg.Amount)
	if err != nil {
		ctx.Logger().Error("DEBUG: Failed to parse amount", "amount", msg.Amount, "error", err)
		return nil, errorsmod.Wrapf(types.ErrInvalidAmount, "failed to parse amount: %s", err)
	}

	ctx.Logger().Info("DEBUG: Parsed burn amount", "amount", burnCoin.String())

	// 2. Get module parameters
	params, err := k.Params.Get(ctx)
	if err != nil {
		ctx.Logger().Error("DEBUG: Failed to get params", "error", err)
		return nil, errorsmod.Wrap(err, "failed to get module parameters")
	}

	ctx.Logger().Info("DEBUG: Got params", "genesis_address", params.GenesisAddress)

	// 3. Validate genesis address is configured
	if params.GenesisAddress == "" {
		ctx.Logger().Error("DEBUG: Genesis address not configured")
		return nil, errorsmod.Wrap(types.ErrNotGenesisAddress, "genesis address not configured in module parameters")
	}

	// 4. Validate caller is the genesis address
	if msg.Creator != params.GenesisAddress {
		ctx.Logger().Error("DEBUG: Caller is not genesis address", "caller", msg.Creator, "genesis_address", params.GenesisAddress)
		return nil, errorsmod.Wrapf(types.ErrNotGenesisAddress, "caller is not the configured genesis address")
	}

	ctx.Logger().Info("DEBUG: Caller validation passed")

	// 5. Get all validators for safety check
	validators, err := k.stakingKeeper.GetAllValidators(ctx)
	if err != nil {
		ctx.Logger().Error("DEBUG: Failed to get validators", "error", err)
		return nil, errorsmod.Wrap(err, "failed to get validators")
	}

	ctx.Logger().Info("DEBUG: Got validators", "count", len(validators))

	// 6. Check minimum validator count for safety (need at least 1 for testing)
	activeValidators := 0
	for _, val := range validators {
		if val.Status == stakingtypes.Bonded {
			activeValidators++
		}
	}

	ctx.Logger().Info("DEBUG: Found active validators", "count", activeValidators)

	if activeValidators < 1 {
		return nil, errorsmod.Wrap(types.ErrInsufficientValidators, "need at least 1 active validator for safe execution")
	}

	// 7. Find ONLY Alice's own validator and unbond ONLY the specified amount from her self-delegation
	genesisAddr, err := k.addressCodec.StringToBytes(params.GenesisAddress)
	if err != nil {
		ctx.Logger().Error("DEBUG: Failed to decode genesis address", "error", err)
		return nil, errorsmod.Wrap(err, "failed to decode genesis address")
	}

	ctx.Logger().Info("DEBUG: Decoded genesis address")

	// Find Alice's validator (where she is the operator)
	var aliceValidator *stakingtypes.Validator
	for _, validator := range validators {
		// Convert validator operator address to account address for comparison
		valAddr, err := sdk.ValAddressFromBech32(validator.OperatorAddress)
		if err != nil {
			continue
		}

		// Convert validator address to account address
		valAccAddr := sdk.AccAddress(valAddr)

		// Check if this validator's operator is Alice
		if valAccAddr.String() == sdk.AccAddress(genesisAddr).String() {
			aliceValidator = &validator
			break
		}
	}

	if aliceValidator == nil {
		ctx.Logger().Error("DEBUG: Alice is not a validator operator")
		return nil, errorsmod.Wrap(types.ErrNoStakeToBurn, "genesis address is not a validator operator")
	}

	ctx.Logger().Info("DEBUG: Found Alice's validator", "operator", aliceValidator.OperatorAddress)

	// 8. Get Alice's self-delegation to her own validator
	valAddr, err := sdk.ValAddressFromBech32(aliceValidator.OperatorAddress)
	if err != nil {
		ctx.Logger().Error("DEBUG: Failed to decode validator address", "error", err)
		return nil, errorsmod.Wrap(err, "failed to decode validator address")
	}

	// Check Alice's delegation to her own validator
	delegation, err := k.stakingKeeper.GetDelegation(ctx, sdk.AccAddress(genesisAddr), valAddr)
	if err != nil {
		ctx.Logger().Error("DEBUG: Alice has no self-delegation", "error", err)
		return nil, errorsmod.Wrap(types.ErrNoStakeToBurn, "genesis address has no self-delegation to unbond")
	}

	if delegation.Shares.IsZero() {
		ctx.Logger().Error("DEBUG: Alice has zero shares")
		return nil, errorsmod.Wrap(types.ErrNoStakeToBurn, "genesis address has no stake to burn")
	}

	ctx.Logger().Info("DEBUG: Found Alice's self-delegation", "shares", delegation.Shares.String())

	// 9. Calculate shares to unbond based on requested amount
	validator, err := k.stakingKeeper.GetValidator(ctx, valAddr)
	if err != nil {
		ctx.Logger().Error("DEBUG: Failed to get validator", "error", err)
		return nil, errorsmod.Wrap(err, "failed to get validator")
	}

	// Convert burn amount to shares - check if requested amount exceeds available tokens
	availableTokens := validator.TokensFromShares(delegation.Shares)
	if burnCoin.Amount.GT(availableTokens.TruncateInt()) {
		ctx.Logger().Error("DEBUG: Requested amount exceeds available stake", "requested", burnCoin.Amount.String(), "available", availableTokens.TruncateInt().String())
		return nil, errorsmod.Wrapf(types.ErrInvalidAmount, "requested amount (%s) exceeds available stake (%s)", burnCoin.String(), availableTokens.TruncateInt().String())
	}

	// Calculate actual shares to unbond (convert amount to shares)
	actualSharesToUnbond, err := validator.SharesFromTokens(burnCoin.Amount)
	if err != nil {
		ctx.Logger().Error("DEBUG: Failed to convert amount to shares", "error", err)
		return nil, errorsmod.Wrap(err, "failed to convert amount to shares")
	}

	ctx.Logger().Info("DEBUG: Calculated shares to unbond", "requested_amount", burnCoin.String(), "shares_to_unbond", actualSharesToUnbond.String(), "total_shares", delegation.Shares.String())

	// 10. Unbond the specified amount from Alice's self-delegation
	ctx.Logger().Info("DEBUG: About to unbond specified amount from Alice's self-delegation")

	unbondTime, unbondedAmount, err := k.stakingKeeper.Undelegate(ctx, sdk.AccAddress(genesisAddr), valAddr, actualSharesToUnbond)
	if err != nil {
		ctx.Logger().Error("DEBUG: Undelegate failed", "error", err)
		return nil, errorsmod.Wrapf(err, "failed to undelegate genesis self-delegation")
	}

	ctx.Logger().Info("DEBUG: Successfully unbonded amount from Alice's self-delegation", "unbonded", unbondedAmount.String(), "completion_time", unbondTime.String())

	// 11. Store pending burn info (we can now have multiple pending burns)
	completionTime := unbondTime.Unix()
	// Use completion time as key to allow multiple pending burns
	if err := k.PendingGenesisBurns.Set(ctx, uint64(completionTime), completionTime); err != nil {
		ctx.Logger().Error("DEBUG: Failed to set pending burn", "error", err)
		return nil, errorsmod.Wrap(err, "failed to store pending burn information")
	}

	ctx.Logger().Info("DEBUG: Stored pending burn info")

	// 12. Emit event
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeBurnGenesisStake,
			sdk.NewAttribute(types.AttributeKeyGenesisAddress, params.GenesisAddress),
			sdk.NewAttribute(types.AttributeKeyUnbondedAmount, unbondedAmount.String()),
			sdk.NewAttribute(types.AttributeKeyCompletionTime, fmt.Sprintf("%d", completionTime)),
			sdk.NewAttribute("requested_amount", burnCoin.String()),
		),
	)

	ctx.Logger().Info("DEBUG: Emitted event, returning response")

	return &types.MsgBurnGenesisStakeResponse{
		UnbondedAmount:          unbondedAmount.String(),
		UnbondingCompletionTime: completionTime,
	}, nil
}

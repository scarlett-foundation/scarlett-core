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

	ctx.Logger().Info("DEBUG: Starting BurnGenesisStake handler")

	// 1. Get module parameters
	params, err := k.Params.Get(ctx)
	if err != nil {
		ctx.Logger().Error("DEBUG: Failed to get params", "error", err)
		return nil, errorsmod.Wrap(err, "failed to get module parameters")
	}

	ctx.Logger().Info("DEBUG: Got params", "genesis_address", params.GenesisAddress)

	// 2. Validate genesis address is configured
	if params.GenesisAddress == "" {
		ctx.Logger().Error("DEBUG: Genesis address not configured")
		return nil, errorsmod.Wrap(types.ErrNotGenesisAddress, "genesis address not configured in module parameters")
	}

	// 3. Validate caller is the genesis address
	if msg.Creator != params.GenesisAddress {
		ctx.Logger().Error("DEBUG: Caller is not genesis address", "caller", msg.Creator, "genesis_address", params.GenesisAddress)
		return nil, errorsmod.Wrapf(types.ErrNotGenesisAddress, "caller is not the configured genesis address")
	}

	ctx.Logger().Info("DEBUG: Caller validation passed")

	// 4. Check if already executed by looking for existing pending burns for this address
	// Use a simple hash of genesis address as key
	addressKey := uint64(0) // Simple key for now - in production would hash the address
	alreadyExecuted, err := k.PendingGenesisBurns.Get(ctx, addressKey)
	if err == nil && alreadyExecuted > 0 {
		ctx.Logger().Error("DEBUG: Already executed")
		return nil, errorsmod.Wrap(types.ErrAlreadyExecuted, "genesis stake burn already executed")
	}

	ctx.Logger().Info("DEBUG: One-time check passed")

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

	// 7. Find ONLY Alice's own validator and unbond ONLY her self-delegation
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

	// 9. Unbond ONLY Alice's self-delegation
	ctx.Logger().Info("DEBUG: About to unbond Alice's self-delegation")

	unbondTime, unbondedAmount, err := k.stakingKeeper.Undelegate(ctx, sdk.AccAddress(genesisAddr), valAddr, delegation.Shares)
	if err != nil {
		ctx.Logger().Error("DEBUG: Undelegate failed", "error", err)
		return nil, errorsmod.Wrapf(err, "failed to undelegate genesis self-delegation")
	}

	ctx.Logger().Info("DEBUG: Successfully unbonded Alice's self-delegation", "unbonded", unbondedAmount.String(), "completion_time", unbondTime.String())

	// 10. Mark as executed to prevent re-execution
	completionTime := unbondTime.Unix()
	if err := k.PendingGenesisBurns.Set(ctx, addressKey, completionTime); err != nil {
		ctx.Logger().Error("DEBUG: Failed to set execution flag", "error", err)
		return nil, errorsmod.Wrap(err, "failed to mark genesis burn as executed")
	}

	ctx.Logger().Info("DEBUG: Set execution flag")

	// 11. Emit event
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeBurnGenesisStake,
			sdk.NewAttribute(types.AttributeKeyGenesisAddress, params.GenesisAddress),
			sdk.NewAttribute(types.AttributeKeyUnbondedAmount, unbondedAmount.String()),
			sdk.NewAttribute(types.AttributeKeyCompletionTime, fmt.Sprintf("%d", completionTime)),
		),
	)

	ctx.Logger().Info("DEBUG: Emitted event, returning response")

	return &types.MsgBurnGenesisStakeResponse{
		UnbondedAmount:          unbondedAmount.String(),
		UnbondingCompletionTime: completionTime,
	}, nil
}

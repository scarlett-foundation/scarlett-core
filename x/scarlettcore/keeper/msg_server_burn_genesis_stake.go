package keeper

import (
	"context"
	"strconv"

	"scarlett-core/x/scarlettcore/types"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
)

func (k msgServer) BurnGenesisStake(ctx context.Context, msg *types.MsgBurnGenesisStake) (*types.MsgBurnGenesisStakeResponse, error) {
	// Check if staking keeper is available
	if k.stakingKeeper == nil {
		return nil, errorsmod.Wrap(types.ErrNoStakeToBurn, "staking functionality not available")
	}

	// 1. Validate caller address
	callerAddr, err := k.addressCodec.StringToBytes(msg.Creator)
	if err != nil {
		return nil, errorsmod.Wrap(err, "invalid caller address")
	}

	// 2. Get module params to check genesis address
	params, err := k.Params.Get(ctx)
	if err != nil {
		return nil, errorsmod.Wrap(err, "failed to get module params")
	}

	// 3. Validate caller is genesis address
	if params.GenesisAddress == "" {
		return nil, errorsmod.Wrap(types.ErrNotGenesisAddress, "genesis address not configured")
	}

	genesisAddr, err := k.addressCodec.StringToBytes(params.GenesisAddress)
	if err != nil {
		return nil, errorsmod.Wrap(err, "invalid configured genesis address")
	}

	if !sdk.AccAddress(callerAddr).Equals(sdk.AccAddress(genesisAddr)) {
		return nil, errorsmod.Wrap(types.ErrNotGenesisAddress, "caller is not the configured genesis address")
	}

	// 4. Check if genesis burn already executed (check for any pending burns)
	hasPendingBurns := false
	err = k.PendingGenesisBurns.Walk(ctx, nil, func(key uint64, value int64) (bool, error) {
		hasPendingBurns = true
		return true, nil // Stop iteration, we found at least one
	})
	if err != nil {
		return nil, errorsmod.Wrap(err, "failed to check pending burns")
	}
	if hasPendingBurns {
		return nil, errorsmod.Wrap(types.ErrAlreadyExecuted, "genesis burn already in progress")
	}

	// 5. Get all validators to find genesis validator
	validators, err := k.stakingKeeper.GetAllValidators(ctx)
	if err != nil {
		return nil, errorsmod.Wrap(err, "failed to get validators")
	}

	// 6. Check minimum validator count for safety (need at least 2 others besides genesis)
	activeValidators := 0
	var genesisValidator *stakingtypes.Validator
	for _, val := range validators {
		if val.Status == stakingtypes.Bonded {
			activeValidators++
			// Check if this validator belongs to genesis address
			valAddr, err := k.addressCodec.StringToBytes(val.OperatorAddress)
			if err == nil {
				// Convert validator operator address to delegator address for comparison
				// In Cosmos SDK, the validator operator and delegator are often the same for self-delegation
				if sdk.AccAddress(genesisAddr).Equals(sdk.AccAddress(valAddr)) {
					genesisValidator = &val
				}
			}
		}
	}

	if activeValidators < 3 {
		return nil, errorsmod.Wrap(types.ErrInsufficientValidators, "need at least 3 active validators for safe execution")
	}

	// 7. If no direct match, find validator by self-delegation
	if genesisValidator == nil {
		for _, val := range validators {
			if val.Status == stakingtypes.Bonded {
				valAddr, _ := sdk.ValAddressFromBech32(val.OperatorAddress)
				delegation, err := k.stakingKeeper.GetDelegation(ctx, sdk.AccAddress(genesisAddr), valAddr)
				if err == nil && !delegation.Shares.IsZero() {
					// Found a delegation from genesis address
					genesisValidator = &val
					break
				}
			}
		}
	}

	if genesisValidator == nil {
		return nil, errorsmod.Wrap(types.ErrNoStakeToBurn, "genesis address has no validator or stake")
	}

	// 8. Get the delegation amount
	valAddr, err := sdk.ValAddressFromBech32(genesisValidator.OperatorAddress)
	if err != nil {
		return nil, errorsmod.Wrap(err, "invalid validator address")
	}

	delegation, err := k.stakingKeeper.GetDelegation(ctx, sdk.AccAddress(genesisAddr), valAddr)
	if err != nil {
		return nil, errorsmod.Wrap(types.ErrNoStakeToBurn, "no delegation found")
	}

	if delegation.Shares.IsZero() {
		return nil, errorsmod.Wrap(types.ErrNoStakeToBurn, "no shares to unbond")
	}

	// 9. Unbond all shares
	completionTime, unbondedAmount, err := k.stakingKeeper.Undelegate(ctx, sdk.AccAddress(genesisAddr), valAddr, delegation.Shares)
	if err != nil {
		return nil, errorsmod.Wrap(err, "failed to unbond stake")
	}

	// 10. Track for burning when unbonding completes
	sdkCtx := sdk.UnwrapSDKContext(ctx)

	// Store the unbonding ID -> completion time mapping
	// Use a simple incrementing ID since we don't have access to the actual unbonding delegation
	unbondingID := uint64(sdkCtx.BlockHeight()) // Use block height as unique ID
	err = k.PendingGenesisBurns.Set(ctx, unbondingID, completionTime.Unix())
	if err != nil {
		return nil, errorsmod.Wrap(err, "failed to track pending burn")
	}

	// 11. Emit genesis burn event
	sdkCtx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeBurnGenesisStake,
			sdk.NewAttribute(types.AttributeKeyGenesisAddress, msg.Creator),
			sdk.NewAttribute(types.AttributeKeyUnbondedAmount, unbondedAmount.String()),
			sdk.NewAttribute(types.AttributeKeyCompletionTime, strconv.FormatInt(completionTime.Unix(), 10)),
			sdk.NewAttribute(types.AttributeKeyUnbondingID, strconv.FormatUint(unbondingID, 10)),
		),
	)

	return &types.MsgBurnGenesisStakeResponse{
		UnbondedAmount:          unbondedAmount.String(),
		UnbondingCompletionTime: completionTime.Unix(),
	}, nil
}

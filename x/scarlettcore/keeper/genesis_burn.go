package keeper

import (
	"context"
	"strconv"

	"scarlett-core/x/scarlettcore/types"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// ProcessCompletedGenesisBurns checks for completed unbonding and burns the tokens
func (k Keeper) ProcessCompletedGenesisBurns(ctx context.Context) error {
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	currentTime := sdkCtx.BlockTime().Unix()

	// Get module params to get genesis address
	params, err := k.Params.Get(ctx)
	if err != nil {
		return errorsmod.Wrap(err, "failed to get module params")
	}

	if params.GenesisAddress == "" {
		// No genesis address configured, nothing to process
		return nil
	}

	genesisAddr, err := k.addressCodec.StringToBytes(params.GenesisAddress)
	if err != nil {
		return errorsmod.Wrap(err, "invalid genesis address")
	}

	// Check all pending burns
	var completedBurns []uint64
	err = k.PendingGenesisBurns.Walk(ctx, nil, func(unbondingID uint64, completionTime int64) (bool, error) {
		if currentTime >= completionTime {
			// This unbonding has completed, add to burn list
			completedBurns = append(completedBurns, unbondingID)
		}
		return false, nil // Continue iteration
	})
	if err != nil {
		return errorsmod.Wrap(err, "failed to walk pending burns")
	}

	// Process completed burns
	for _, unbondingID := range completedBurns {
		if err := k.burnCompletedUnbonding(ctx, sdk.AccAddress(genesisAddr), unbondingID); err != nil {
			// Log error but continue processing other burns
			sdkCtx.Logger().Error("failed to burn completed unbonding", "unbonding_id", unbondingID, "error", err)
			continue
		}

		// Remove from pending burns
		if err := k.PendingGenesisBurns.Remove(ctx, unbondingID); err != nil {
			sdkCtx.Logger().Error("failed to remove pending burn", "unbonding_id", unbondingID, "error", err)
		}
	}

	return nil
}

// burnCompletedUnbonding burns tokens from a completed unbonding
func (k Keeper) burnCompletedUnbonding(ctx context.Context, genesisAddr sdk.AccAddress, unbondingID uint64) error {
	sdkCtx := sdk.UnwrapSDKContext(ctx)

	// Get the balance of genesis address to see what can be burned
	balance := k.bankKeeper.GetAllBalances(ctx, genesisAddr)
	if balance.IsZero() {
		// No balance to burn
		return nil
	}

	// For now, burn all available balance as this represents the unbonded tokens
	// In a production system, you might want to track the exact unbonding amounts
	if err := k.bankKeeper.SendCoinsFromAccountToModule(ctx, genesisAddr, types.ModuleName, balance); err != nil {
		return errorsmod.Wrap(err, "failed to transfer tokens to module for burning")
	}

	if err := k.bankKeeper.BurnCoins(ctx, types.ModuleName, balance); err != nil {
		return errorsmod.Wrap(err, "failed to burn tokens")
	}

	// Emit burn completion event
	sdkCtx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeGenesisDecentralized,
			sdk.NewAttribute(types.AttributeKeyGenesisAddress, genesisAddr.String()),
			sdk.NewAttribute(types.AttributeKeyBurnedAmount, balance.String()),
			sdk.NewAttribute(types.AttributeKeyUnbondingID, strconv.FormatUint(unbondingID, 10)),
		),
	)

	return nil
}

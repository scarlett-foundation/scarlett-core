package app

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	mintkeeper "github.com/cosmos/cosmos-sdk/x/mint/keeper"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"
)

// MinerEmissionsSplitMintFn implements custom 50/50 emission splitting.
// It's designed to be used as a custom MintFn with the x/mint module keeper.
func MinerEmissionsSplitMintFn(ctx sdk.Context, k *mintkeeper.Keeper, bk bankkeeper.Keeper) error {
	// 1. Get current minter state and parameters
	minter, err := k.Minter.Get(ctx)
	if err != nil {
		return err
	}
	params, err := k.Params.Get(ctx)
	if err != nil {
		return err
	}

	// 2. Calculate total supply for provisions calculation
	stakingTokenSupply, err := k.StakingTokenSupply(ctx)
	if err != nil {
		return err
	}
	bondedRatio, err := k.BondedRatio(ctx)
	if err != nil {
		return err
	}

	// 3. Calculate inflation using standard mint logic
	minter.Inflation = minter.NextInflationRate(params, bondedRatio)
	minter.AnnualProvisions = minter.NextAnnualProvisions(params, stakingTokenSupply)

	// 4. Calculate this block's provision
	blockProvisionCoin := minter.BlockProvision(params)
	blockProvisionAmount := blockProvisionCoin.Amount

	// 5. Mint total amount to mint module first
	// The MintCoins method increases the supply and sends the coins to the mint module's account.
	if err := k.MintCoins(ctx, sdk.NewCoins(blockProvisionCoin)); err != nil {
		return err
	}

	// 6. Split emissions 50/50
	// sdk.Int.QuoRaw performs integer division, a subsequent Sub ensures the total is conserved.
	stakingAmount := blockProvisionAmount.QuoRaw(2)
	minerAmount := blockProvisionAmount.Sub(stakingAmount) // Remaining (handles odd amounts correctly)

	stakingCoins := sdk.NewCoins(sdk.NewCoin(params.MintDenom, stakingAmount))
	minerCoins := sdk.NewCoins(sdk.NewCoin(params.MintDenom, minerAmount))

	// 7. Transfer to FeeCollector (staking rewards)
	// These coins are sent from the mint module's account to the FeeCollector module account.
	if err := bk.SendCoinsFromModuleToModule(ctx, minttypes.ModuleName, authtypes.FeeCollectorName, stakingCoins); err != nil {
		return err
	}

	// 8. Transfer to inference rewards module account
	// These coins are sent from the mint module's account to the "inferencerewards" module account.
	// Ensure "inferencerewards" module account exists and has permissions to receive coins.
	if err := bk.SendCoinsFromModuleToModule(ctx, minttypes.ModuleName, "inferencerewards", minerCoins); err != nil {
		return err
	}

	// 9. Update minter state in the store
	k.Minter.Set(ctx, minter)

	// 10. Emit events
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			minttypes.EventTypeMint,
			sdk.NewAttribute(minttypes.AttributeKeyBondedRatio, bondedRatio.String()),
			sdk.NewAttribute(minttypes.AttributeKeyInflation, minter.Inflation.String()),
			sdk.NewAttribute(minttypes.AttributeKeyAnnualProvisions, minter.AnnualProvisions.String()),
			sdk.NewAttribute(sdk.AttributeKeyAmount, blockProvisionCoin.Amount.String()), // Total minted amount
			sdk.NewAttribute("staking_amount", stakingAmount.String()),
			sdk.NewAttribute("miner_amount", minerAmount.String()),
		),
	)

	return nil
}

package app

import (
	"fmt"

	"cosmossdk.io/log"

	sdk "github.com/cosmos/cosmos-sdk/types"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	mintkeeper "github.com/cosmos/cosmos-sdk/x/mint/keeper"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"

	"scarlett-core/app/emissions"
	emissionsmodulekeeper "scarlett-core/x/emissions/keeper"
)

// ProvideDynamicEmissionsMintFn provides a dynamic mint function using the emissions keeper
func ProvideDynamicEmissionsMintFn(emissionsKeeper emissionsmodulekeeper.Keeper, logger log.Logger) mintkeeper.MintFn {
	// CRITICAL DEBUG: Log that dependency injection is working
	logger.Info("🔧🔧🔧 DEPENDENCY INJECTION: ProvideDynamicEmissionsMintFn CALLED 🔧🔧🔧")

	return func(ctx sdk.Context, mintKeeper *mintkeeper.Keeper) error {
		// CRITICAL DEBUG: This should appear if our function is being called
		ctx.Logger().Info("🚨🚨🚨 CUSTOM MINT FUNCTION CALLED 🚨🚨🚨", "height", ctx.BlockHeight())
		return emissionsKeeper.DynamicEmissionsMintFn(ctx, mintKeeper)
	}
}

// DEPRECATED: Legacy functions kept for backward compatibility during transition

// ProvideMinerEmissionsMintFn is a depinject provider for our custom mint function (LEGACY)
// This is kept for backward compatibility but should be replaced with ProvideDynamicEmissionsMintFn
func ProvideMinerEmissionsMintFn(bankKeeper bankkeeper.Keeper) mintkeeper.MintFn {
	return MinerEmissionsSplitMintFnFactory(bankKeeper)
}

// MinerEmissionsSplitMintFnFactory creates a custom mint function with the required dependencies (LEGACY)
// This follows the depinject pattern for Cosmos SDK v0.53.0 but uses hardcoded parameters.
// DEPRECATED: Use the dynamic governance-controlled version instead.
func MinerEmissionsSplitMintFnFactory(bk bankkeeper.Keeper) mintkeeper.MintFn {
	return func(ctx sdk.Context, k *mintkeeper.Keeper) error {
		return MinerEmissionsSplitMintFn(ctx, k, bk)
	}
}

// MinerEmissionsSplitMintFn implements custom emission splitting using the modular emissions package (LEGACY)
// It's designed to be used as a custom MintFn with the x/mint module keeper.
// DEPRECATED: This function uses hardcoded emission parameters. Use the dynamic governance-controlled
// version from x/emissions module instead.
func MinerEmissionsSplitMintFn(ctx sdk.Context, k *mintkeeper.Keeper, bk bankkeeper.Keeper) error {
	// 1. Initialize modular emissions system with hardcoded config
	config := emissions.DefaultEmissionsConfig()
	splitter, err := emissions.NewEmissionSplitter(config, bk)
	if err != nil {
		return fmt.Errorf("failed to create emission splitter: %w", err)
	}

	// 2. Get current minter state and parameters using keeper methods
	minter, err := getMinterState(ctx, k)
	if err != nil {
		return fmt.Errorf("failed to get minter state: %w", err)
	}

	params, err := k.Params.Get(ctx)
	if err != nil {
		return fmt.Errorf("failed to get mint params: %w", err)
	}

	// 3. Calculate inflation and provisions using standard mint logic
	if err := updateInflation(ctx, k, &minter, params); err != nil {
		return fmt.Errorf("failed to update inflation: %w", err)
	}

	// 4. Calculate and mint this block's provision
	blockProvisionCoin := minter.BlockProvision(params)
	if err := mintTokens(ctx, k, blockProvisionCoin); err != nil {
		return fmt.Errorf("failed to mint tokens: %w", err)
	}

	// 5. Distribute minted tokens using modular splitter
	if err := splitter.DistributeTokens(ctx, params.MintDenom, blockProvisionCoin.Amount); err != nil {
		return fmt.Errorf("failed to distribute tokens: %w", err)
	}

	// 6. Update minter state in the store
	if err := k.Minter.Set(ctx, minter); err != nil {
		return fmt.Errorf("failed to update minter state: %w", err)
	}

	// 7. Emit comprehensive events (legacy format)
	distributions, _ := splitter.CalculateDistribution(blockProvisionCoin.Amount)
	bondedRatio, _ := k.BondedRatio(ctx)

	emissions.EmitEmissionSplitEvent(
		ctx,
		blockProvisionCoin.Amount,
		distributions,
		config,
		minter,
		bondedRatio,
	)

	return nil
}

// LEGACY HELPER FUNCTIONS (kept for backward compatibility)

// getMinterState retrieves and validates the current minter state
func getMinterState(ctx sdk.Context, k *mintkeeper.Keeper) (minttypes.Minter, error) {
	minter, err := k.Minter.Get(ctx)
	if err != nil {
		return minttypes.Minter{}, fmt.Errorf("failed to get minter: %w", err)
	}
	return minter, nil
}

// updateInflation calculates and updates inflation rates and annual provisions
func updateInflation(ctx sdk.Context, k *mintkeeper.Keeper, minter *minttypes.Minter, params minttypes.Params) error {
	// Get staking metrics
	stakingTokenSupply, err := k.StakingTokenSupply(ctx)
	if err != nil {
		return fmt.Errorf("failed to get staking token supply: %w", err)
	}

	bondedRatio, err := k.BondedRatio(ctx)
	if err != nil {
		return fmt.Errorf("failed to get bonded ratio: %w", err)
	}

	// Calculate inflation using standard mint logic
	minter.Inflation = minter.NextInflationRate(params, bondedRatio)
	minter.AnnualProvisions = minter.NextAnnualProvisions(params, stakingTokenSupply)

	return nil
}

// mintTokens mints the specified amount of tokens to the mint module account
func mintTokens(ctx sdk.Context, k *mintkeeper.Keeper, coin sdk.Coin) error {
	if coin.Amount.IsZero() {
		return nil // Nothing to mint
	}

	coins := sdk.NewCoins(coin)
	if err := k.MintCoins(ctx, coins); err != nil {
		return fmt.Errorf("failed to mint %s: %w", coins.String(), err)
	}

	return nil
}

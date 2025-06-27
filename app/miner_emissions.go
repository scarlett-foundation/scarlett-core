package app

import (
	"encoding/json"
	"fmt"

	"cosmossdk.io/log"
	"cosmossdk.io/math"

	sdk "github.com/cosmos/cosmos-sdk/types"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	mintkeeper "github.com/cosmos/cosmos-sdk/x/mint/keeper"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"

	"scarlett-core/app/emissions"
	emissionsmodulekeeper "scarlett-core/x/emissions/keeper"
)

// KeeperProvider defines a function that returns the necessary keepers.
// This is used to break a dependency cycle during initialization.
type KeeperProvider func() (bankkeeper.Keeper, emissionsmodulekeeper.Keeper)

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

// ProvideMinerEmissionsMintFn provides a dynamic mint function that uses a provider
// to lazily load keepers, preventing initialization dependency cycles.
func ProvideMinerEmissionsMintFn(keeperProvider KeeperProvider) mintkeeper.MintFn {
	return func(ctx sdk.Context, mk *mintkeeper.Keeper) error {
		// Lazily load the keepers only when the mint function is actually called.
		bankKeeper, emissionsKeeper := keeperProvider()
		return GovernanceControlledMintFn(ctx, mk, bankKeeper, emissionsKeeper)
	}
}

// GovernanceControlledMintFn is the core minting logic that uses governance parameters.
func GovernanceControlledMintFn(
	ctx sdk.Context,
	mintKeeper *mintkeeper.Keeper,
	bankKeeper bankkeeper.Keeper,
	emissionsKeeper emissionsmodulekeeper.Keeper,
) error {
	ctx.Logger().Info("🔥🔥🔥 GOVERNANCE-CONTROLLED MINT FUNCTION CALLED 🔥🔥🔥", "height", ctx.BlockHeight())

	var config emissions.EmissionsConfig

	// 1. Try to get parameters from governance
	emissionParams, err := emissionsKeeper.Params.Get(ctx)
	// If parameters don't exist or are explicitly disabled, use the default 50/50 split.
	if err != nil || !emissionParams.Enabled {
		if err != nil {
			ctx.Logger().Info("Could not get governance parameters, falling back to default 50/50 split", "error", err)
		} else {
			ctx.Logger().Info("Governance emissions are explicitly disabled, falling back to default 50/50 split")
		}
		config = emissions.DefaultEmissionsConfig()
	} else {
		// Governance parameters exist and are enabled, so we parse and use them.
		ctx.Logger().Info("Using governance-defined emission parameters")
		var destinations []struct {
			ModuleName  string `json:"module_name"`
			Weight      string `json:"weight"`
			Description string `json:"description"`
			Enabled     bool   `json:"enabled"`
		}

		if err := json.Unmarshal([]byte(emissionParams.EmissionDestinations), &destinations); err != nil {
			ctx.Logger().Error("Failed to parse governance emission destinations JSON, falling back to default", "error", err)
			config = emissions.DefaultEmissionsConfig()
		} else {
			parsedConfig := emissions.EmissionsConfig{
				Enabled:      true,
				Destinations: make([]emissions.EmissionDestination, 0, len(destinations)),
			}
			for _, dest := range destinations {
				if !dest.Enabled {
					continue // Skip disabled destinations in the governance config
				}
				weight, err := math.LegacyNewDecFromStr(dest.Weight)
				if err != nil {
					ctx.Logger().Error("Invalid weight in governance destination, falling back to default", "module", dest.ModuleName, "error", err)
					config = emissions.DefaultEmissionsConfig()
					goto end_parse // Exit loop and use default config
				}
				parsedConfig.Destinations = append(parsedConfig.Destinations, emissions.EmissionDestination{
					ModuleName:  dest.ModuleName,
					Weight:      weight,
					Description: dest.Description,
				})
			}
			config = parsedConfig
		}
	}
end_parse:

	ctx.Logger().Info("📊 Using emission configuration",
		"enabled", config.Enabled,
		"destinations", len(config.Destinations),
	)

	// 2. Initialize modular emissions system with governance config
	splitter, err := emissions.NewEmissionSplitter(config, bankKeeper)
	if err != nil {
		return fmt.Errorf("failed to create emission splitter: %w", err)
	}

	// 3. Get current minter state and parameters using keeper methods
	minter, err := getMinterState(ctx, mintKeeper)
	if err != nil {
		return fmt.Errorf("failed to get minter state: %w", err)
	}

	mintParams, err := mintKeeper.Params.Get(ctx)
	if err != nil {
		return fmt.Errorf("failed to get mint params: %w", err)
	}

	// 4. Calculate inflation and provisions using standard mint logic
	if err := updateInflation(ctx, mintKeeper, &minter, mintParams); err != nil {
		return fmt.Errorf("failed to update inflation: %w", err)
	}

	// 5. Calculate and mint this block's provision
	blockProvisionCoin := minter.BlockProvision(mintParams)
	if err := mintTokens(ctx, mintKeeper, blockProvisionCoin); err != nil {
		return fmt.Errorf("failed to mint tokens: %w", err)
	}

	ctx.Logger().Info("💰 Minted tokens for distribution",
		"amount", blockProvisionCoin.Amount.String(),
		"denom", blockProvisionCoin.Denom,
	)

	// 6. Distribute minted tokens using governance-controlled splitter
	if err := splitter.DistributeTokens(ctx, mintParams.MintDenom, blockProvisionCoin.Amount); err != nil {
		return fmt.Errorf("failed to distribute tokens: %w", err)
	}

	// 7. Update minter state in the store
	if err := mintKeeper.Minter.Set(ctx, minter); err != nil {
		return fmt.Errorf("failed to update minter state: %w", err)
	}

	// 8. Emit comprehensive events
	distributions, _ := splitter.CalculateDistribution(blockProvisionCoin.Amount)
	bondedRatio, _ := mintKeeper.BondedRatio(ctx)

	emissions.EmitEmissionSplitEvent(
		ctx,
		blockProvisionCoin.Amount,
		distributions,
		config,
		minter,
		bondedRatio,
	)

	ctx.Logger().Info("✅ Governance-controlled emission distribution complete")
	return nil
}

// getGovernanceEmissionConfig reads governance parameters and converts them to EmissionsConfig
func getGovernanceEmissionConfig(ctx sdk.Context, ek emissionsmodulekeeper.Keeper) (emissions.EmissionsConfig, error) {
	// Get governance parameters from emissions module
	params, err := ek.Params.Get(ctx)
	if err != nil {
		return emissions.EmissionsConfig{}, fmt.Errorf("failed to get emissions params: %w", err)
	}

	// If emissions not enabled, return disabled config
	if !params.Enabled {
		return emissions.EmissionsConfig{Enabled: false}, nil
	}

	// Parse emission destinations from JSON
	var destinations []emissionDestination
	if err := json.Unmarshal([]byte(params.EmissionDestinations), &destinations); err != nil {
		return emissions.EmissionsConfig{}, fmt.Errorf("failed to parse emission destinations: %w", err)
	}

	// Convert to emissions package format
	config := emissions.EmissionsConfig{
		Enabled:      true,
		Destinations: make([]emissions.EmissionDestination, len(destinations)),
	}

	for i, dest := range destinations {
		if !dest.Enabled {
			continue // Skip disabled destinations
		}

		weight, err := math.LegacyNewDecFromStr(dest.Weight)
		if err != nil {
			return emissions.EmissionsConfig{}, fmt.Errorf("invalid weight for %s: %w", dest.ModuleName, err)
		}

		config.Destinations[i] = emissions.EmissionDestination{
			ModuleName:  dest.ModuleName,
			Weight:      weight,
			Description: dest.Description,
		}
	}

	return config, nil
}

// emissionDestination represents the JSON structure from governance parameters
type emissionDestination struct {
	ModuleName  string `json:"module_name"`
	Weight      string `json:"weight"`
	Description string `json:"description"`
	Enabled     bool   `json:"enabled"`
}

// LEGACY FUNCTIONS (kept for backward compatibility during transition)

// MinerEmissionsSplitMintFnFactory creates a custom mint function with the required dependencies (LEGACY)
// This follows the depinject pattern for Cosmos SDK v0.53.0 but uses hardcoded parameters.
// DEPRECATED: Use the governance-controlled version instead.
func MinerEmissionsSplitMintFnFactory(bk bankkeeper.Keeper) mintkeeper.MintFn {
	return func(ctx sdk.Context, k *mintkeeper.Keeper) error {
		return MinerEmissionsSplitMintFn(ctx, k, bk)
	}
}

// MinerEmissionsSplitMintFn implements custom emission splitting using the modular emissions package (LEGACY)
// It's designed to be used as a custom MintFn with the x/mint module keeper.
// DEPRECATED: This function uses hardcoded emission parameters. Use the governance-controlled
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

// HELPER FUNCTIONS

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

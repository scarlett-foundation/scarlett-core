package app

import (
	"cosmossdk.io/log"
	sdk "github.com/cosmos/cosmos-sdk/types"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	mintkeeper "github.com/cosmos/cosmos-sdk/x/mint/keeper"

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
		_, emissionsKeeper := keeperProvider()
		// NEW: Call consolidated autonomous method instead of complex logic
		return emissionsKeeper.MintAndDistributeEmissions(ctx, mk)
	}
}

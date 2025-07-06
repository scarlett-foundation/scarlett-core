package keeper

import (
	"context"

	"scarlett-core/x/contracts/types"

	wasmtypes "github.com/CosmWasm/wasmd/x/wasm/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// InitGenesis initializes the module's state from a provided genesis state.
func (k Keeper) InitGenesis(ctx context.Context, genState types.GenesisState) {
	// Get the wasm keeper with error handling
	wasmKeeper, err := k.getWasmKeeper()
	if err != nil {
		panic(err) // Genesis initialization should not fail
	}

	wasmKeeper.Logger(sdk.UnwrapSDKContext(ctx)).Info("🚀 Initializing contracts module genesis")

	// Set permissionless parameters for the wasm module
	// This allows anyone to upload and instantiate contracts
	wasmParams := wasmtypes.Params{
		CodeUploadAccess: wasmtypes.AccessConfig{
			Permission: wasmtypes.AccessTypeEverybody,
		},
		InstantiateDefaultPermission: wasmtypes.AccessTypeEverybody,
	}

	// Set the permissionless parameters using the lazy-loaded keeper
	if err := wasmKeeper.SetParams(ctx, wasmParams); err != nil {
		panic(err)
	}

	wasmKeeper.Logger(sdk.UnwrapSDKContext(ctx)).Info("✅ Contracts module genesis initialized with permissionless parameters")
}

// ExportGenesis returns the module's exported genesis state.
func (k Keeper) ExportGenesis(ctx context.Context) *types.GenesisState {
	genesis := types.DefaultGenesis()

	// Get the wasm keeper with error handling
	wasmKeeper, err := k.getWasmKeeper()
	if err != nil {
		panic(err) // Genesis export should not fail
	}

	// Export any wasm-specific state
	// This would typically include code storage, contract instances, etc.
	// For now, we'll just export the basic params
	_ = wasmKeeper

	return genesis
}

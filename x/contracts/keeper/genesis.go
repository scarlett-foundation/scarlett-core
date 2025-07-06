package keeper

import (
	"context"

	"scarlett-core/x/contracts/types"

	wasmtypes "github.com/CosmWasm/wasmd/x/wasm/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// InitGenesis initializes the module's state from a provided genesis state.
func (k Keeper) InitGenesis(ctx context.Context, genState types.GenesisState) {
	// Check if wasm keeper is initialized
	if k.wasmKeeper == nil {
		panic("wasm keeper not initialized during genesis")
	}

	k.wasmKeeper.Logger(sdk.UnwrapSDKContext(ctx)).Info("🚀 Initializing contracts module genesis")

	// Set permissionless parameters for the wasm module
	// This allows anyone to upload and instantiate contracts
	wasmParams := wasmtypes.Params{
		CodeUploadAccess: wasmtypes.AccessConfig{
			Permission: wasmtypes.AccessTypeEverybody,
		},
		InstantiateDefaultPermission: wasmtypes.AccessTypeEverybody,
	}

	// Set the permissionless parameters using the wasm keeper
	if err := k.wasmKeeper.SetParams(ctx, wasmParams); err != nil {
		panic(err)
	}

	k.wasmKeeper.Logger(sdk.UnwrapSDKContext(ctx)).Info("✅ Contracts module genesis initialized with permissionless parameters")
}

// ExportGenesis returns the module's exported genesis state.
func (k Keeper) ExportGenesis(ctx context.Context) *types.GenesisState {
	genesis := types.DefaultGenesis()

	// Check if wasm keeper is initialized
	if k.wasmKeeper == nil {
		panic("wasm keeper not initialized during genesis export")
	}

	// Export any wasm-specific state
	// This would typically include code storage, contract instances, etc.
	// For now, we'll just export the basic params
	_ = k.wasmKeeper

	return genesis
}

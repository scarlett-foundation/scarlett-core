package keeper

import (
	"context"

	wasmtypes "github.com/CosmWasm/wasmd/x/wasm/types"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"scarlett-core/x/contracts/types"
)

// InitGenesis initializes the module's state from a provided genesis state.
// This sets up permissionless parameters for the wasm module.
func (k Keeper) InitGenesis(ctx context.Context, genState types.GenesisState) error {
	// Convert context to SDK context for compatibility
	sdkCtx := sdk.UnwrapSDKContext(ctx)

	k.Keeper.Logger(sdkCtx).Info("🚀 Initializing contracts module genesis")

	// Set permissionless parameters for the wasm module
	// This makes contract deployment permissionless from genesis
	wasmParams := wasmtypes.Params{
		CodeUploadAccess: wasmtypes.AccessConfig{
			Permission: wasmtypes.AccessTypeEverybody,
		},
		InstantiateDefaultPermission: wasmtypes.AccessTypeEverybody,
	}

	// Set the permissionless parameters using the embedded keeper
	if err := k.Keeper.SetParams(ctx, wasmParams); err != nil {
		return err
	}

	k.Keeper.Logger(sdkCtx).Info("✅ Contracts module genesis initialized with permissionless parameters")

	return nil
}

// ExportGenesis returns the module's exported genesis state.
func (k Keeper) ExportGenesis(ctx context.Context) types.GenesisState {
	// Convert context to SDK context for compatibility
	sdkCtx := sdk.UnwrapSDKContext(ctx)

	k.Keeper.Logger(sdkCtx).Info("📤 Exporting contracts module genesis")

	// For our wrapper module, we just export the default genesis state
	// The actual wasm state is handled by the embedded keeper internally
	return *types.DefaultGenesis()
}

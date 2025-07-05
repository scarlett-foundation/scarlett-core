package keeper

import (
	"context"

	"scarlett-core/x/contracts/types"
)

// InitGenesis initializes the module's state from a provided genesis state.
// Since this is a wrapper around wasmd, we don't need to handle genesis state
func (k Keeper) InitGenesis(ctx context.Context, genState types.GenesisState) error {
	// The underlying wasmd keeper handles its own genesis
	// We just need to validate our wrapper is properly initialized
	k.Logger().Info("Initializing contracts module genesis")
	return nil
}

// ExportGenesis returns the module's exported genesis.
// Since this is a wrapper around wasmd, we return default genesis
func (k Keeper) ExportGenesis(ctx context.Context) (*types.GenesisState, error) {
	k.Logger().Info("Exporting contracts module genesis")
	return types.DefaultGenesis(), nil
}

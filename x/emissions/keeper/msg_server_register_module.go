package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"scarlett-core/x/emissions/types"
)

// RegisterModule handles the registration of new modules in the emissions registry
func (m msgServer) RegisterModule(ctx context.Context, req *types.MsgRegisterModule) (*types.MsgRegisterModuleResponse, error) {
	// Convert context to SDK context for logging
	sdkCtx := sdk.UnwrapSDKContext(ctx)

	sdkCtx.Logger().Info("🔧 MODULE REGISTRATION REQUEST",
		"module_name", req.ModuleName,
		"creator", req.Creator,
		"description", req.Description,
	)

	// 1. Validate the module name
	if err := types.ValidateModuleName(req.ModuleName); err != nil {
		sdkCtx.Logger().Error("❌ Invalid module name", "error", err, "module_name", req.ModuleName)
		return nil, err
	}

	// 2. Check if module is already registered
	if m.Keeper.IsModuleRegistered(ctx, req.ModuleName) {
		sdkCtx.Logger().Error("❌ Module already registered", "module_name", req.ModuleName)
		return nil, types.ErrModuleAlreadyRegistered.Wrapf("module '%s' is already registered", req.ModuleName)
	}

	// 3. Verify that the module account exists on-chain
	moduleAccount := m.Keeper.authKeeper.GetModuleAccount(ctx, req.ModuleName)
	if moduleAccount == nil {
		sdkCtx.Logger().Error("❌ Module account does not exist", "module_name", req.ModuleName)
		return nil, types.ErrInvalidDestination.Wrapf("module account '%s' does not exist", req.ModuleName)
	}

	// 4. Create the registered module
	registeredModule := types.NewRegisteredModule(req.ModuleName, req.Creator, req.Description)

	// 5. Validate the registered module
	if err := registeredModule.Validate(); err != nil {
		sdkCtx.Logger().Error("❌ Invalid registered module", "error", err)
		return nil, err
	}

	// 6. Store the registered module in the registry
	if err := m.Keeper.SetRegisteredModule(ctx, registeredModule); err != nil {
		sdkCtx.Logger().Error("❌ Failed to store registered module", "error", err)
		return nil, err
	}

	// 7. Emit registration event
	sdkCtx.EventManager().EmitEvent(
		sdk.NewEvent(
			"module_registered",
			sdk.NewAttribute("module_name", req.ModuleName),
			sdk.NewAttribute("creator", req.Creator),
			sdk.NewAttribute("description", req.Description),
			sdk.NewAttribute("status", string(types.ModuleStatusRegistered)),
			sdk.NewAttribute("module_account", moduleAccount.GetAddress().String()),
		),
	)

	sdkCtx.Logger().Info("✅ MODULE REGISTRATION SUCCESSFUL",
		"module_name", req.ModuleName,
		"creator", req.Creator,
		"status", string(types.ModuleStatusRegistered),
		"module_account", moduleAccount.GetAddress().String(),
	)

	return &types.MsgRegisterModuleResponse{}, nil
}

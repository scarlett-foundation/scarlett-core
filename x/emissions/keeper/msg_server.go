package keeper

import (
	"context"

	"scarlett-core/x/emissions/types"
)

type msgServer struct {
	Keeper
}

// NewMsgServerImpl returns an implementation of the MsgServer interface
// for the provided Keeper.
func NewMsgServerImpl(keeper Keeper) types.MsgServer {
	return &msgServer{Keeper: keeper}
}

var _ types.MsgServer = msgServer{}

// RegisterModule - temporary stub implementation (will be implemented in Step 1.3)
func (m msgServer) RegisterModule(ctx context.Context, req *types.MsgRegisterModule) (*types.MsgRegisterModuleResponse, error) {
	// TODO: Implement in Step 1.3
	return &types.MsgRegisterModuleResponse{}, nil
}

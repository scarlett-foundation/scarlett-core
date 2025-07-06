package keeper

import (
	"context"

	"scarlett-core/x/contracts/types"

	wasmkeeper "github.com/CosmWasm/wasmd/x/wasm/keeper"
	wasmtypes "github.com/CosmWasm/wasmd/x/wasm/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k msgServer) StoreCode(goCtx context.Context, msg *types.MsgStoreCode) (*types.MsgStoreCodeResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Get the wasm keeper with error handling
	wasmKeeper, err := k.getWasmKeeper()
	if err != nil {
		return nil, err
	}

	// Convert our message to wasmd's message format
	wasmMsg := &wasmtypes.MsgStoreCode{
		Sender:       msg.Creator,
		WASMByteCode: msg.WasmByteCode,
		// Set permissionless access by default
		InstantiatePermission: &wasmtypes.AccessConfig{
			Permission: wasmtypes.AccessTypeEverybody,
		},
	}

	// Create a message server for wasmd
	wasmMsgServer := wasmkeeper.NewMsgServerImpl(wasmKeeper)

	// Delegate to wasmd's store code handler
	resp, err := wasmMsgServer.StoreCode(ctx, wasmMsg)
	if err != nil {
		return nil, err
	}

	return &types.MsgStoreCodeResponse{
		CodeId: resp.CodeID,
	}, nil
}

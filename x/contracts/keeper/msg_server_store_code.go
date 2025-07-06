package keeper

import (
	"context"

	"scarlett-core/x/contracts/types"

	errorsmod "cosmossdk.io/errors"
	wasmkeeper "github.com/CosmWasm/wasmd/x/wasm/keeper"
	wasmtypes "github.com/CosmWasm/wasmd/x/wasm/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k msgServer) StoreCode(ctx context.Context, msg *types.MsgStoreCode) (*types.MsgStoreCodeResponse, error) {
	// Validate creator address using SDK validation
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return nil, errorsmod.Wrap(err, "invalid creator address")
	}

	// Create wasmd StoreCode message
	wasmMsg := &wasmtypes.MsgStoreCode{
		Sender:       msg.Creator,
		WASMByteCode: msg.WasmByteCode,
		// Set permissionless access by default
		InstantiatePermission: &wasmtypes.AccessConfig{
			Permission: wasmtypes.AccessTypeEverybody,
		},
	}

	// Get the wasm keeper and create a message server for it
	wasmKeeper := k.getWasmKeeper()
	wasmMsgServer := wasmkeeper.NewMsgServerImpl(wasmKeeper)

	// Delegate to wasmd's StoreCode handler
	wasmResponse, err := wasmMsgServer.StoreCode(ctx, wasmMsg)
	if err != nil {
		return nil, errorsmod.Wrap(err, "failed to store wasm code")
	}

	// Log the successful code storage
	k.Logger(ctx).Info("Successfully stored wasm code",
		"code_id", wasmResponse.CodeID,
		"checksum", wasmResponse.Checksum,
		"creator", msg.Creator)

	// Return response with the generated code ID
	return &types.MsgStoreCodeResponse{
		CodeId: wasmResponse.CodeID,
	}, nil
}

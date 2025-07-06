package keeper

import (
	"context"

	"scarlett-core/x/contracts/types"

	errorsmod "cosmossdk.io/errors"
	wasmkeeper "github.com/CosmWasm/wasmd/x/wasm/keeper"
	wasmtypes "github.com/CosmWasm/wasmd/x/wasm/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k msgServer) ExecuteContract(goCtx context.Context, msg *types.MsgExecuteContract) (*types.MsgExecuteContractResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Validate creator address
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return nil, errorsmod.Wrap(err, "invalid creator address")
	}

	// Check if wasm keeper is initialized
	if k.wasmKeeper == nil {
		return nil, types.ErrWasmNotInitialized
	}

	// Convert our message to wasmd's message format
	wasmMsg := &wasmtypes.MsgExecuteContract{
		Sender:   msg.Creator,
		Contract: msg.Contract,
		Msg:      msg.Msg,
		// Convert funds from string to sdk.Coins
		Funds: nil, // TODO: Parse funds string to sdk.Coins
	}

	// Create a message server for wasmd
	wasmMsgServer := wasmkeeper.NewMsgServerImpl(k.wasmKeeper)

	// Delegate to wasmd's execute contract handler
	resp, err := wasmMsgServer.ExecuteContract(ctx, wasmMsg)
	if err != nil {
		return nil, err
	}

	// Log the successful contract execution
	k.Logger(ctx).Info("Successfully executed contract",
		"contract", msg.Contract,
		"creator", msg.Creator)

	return &types.MsgExecuteContractResponse{
		Data: resp.Data,
	}, nil
}

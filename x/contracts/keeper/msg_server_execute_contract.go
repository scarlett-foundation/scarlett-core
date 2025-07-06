package keeper

import (
	"context"

	"scarlett-core/x/contracts/types"

	errorsmod "cosmossdk.io/errors"
	wasmkeeper "github.com/CosmWasm/wasmd/x/wasm/keeper"
	wasmtypes "github.com/CosmWasm/wasmd/x/wasm/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k msgServer) ExecuteContract(ctx context.Context, msg *types.MsgExecuteContract) (*types.MsgExecuteContractResponse, error) {
	// Validate creator address
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return nil, errorsmod.Wrap(err, "invalid creator address")
	}

	// Create wasmd ExecuteContract message
	wasmMsg := &wasmtypes.MsgExecuteContract{
		Sender:   msg.Creator,
		Contract: msg.Contract,
		Msg:      msg.Msg,
		Funds:    nil, // TODO: Parse funds from string
	}

	// Get the wasm keeper and create a message server for it
	wasmKeeper := k.getWasmKeeper()
	wasmMsgServer := wasmkeeper.NewMsgServerImpl(wasmKeeper)

	// Delegate to wasmd's ExecuteContract handler
	wasmResponse, err := wasmMsgServer.ExecuteContract(ctx, wasmMsg)
	if err != nil {
		return nil, errorsmod.Wrap(err, "failed to execute contract")
	}

	// Log the successful contract execution
	k.Logger(ctx).Info("Successfully executed contract",
		"contract", msg.Contract,
		"creator", msg.Creator)

	// Return response with the execution data
	return &types.MsgExecuteContractResponse{
		Data: wasmResponse.Data,
	}, nil
}

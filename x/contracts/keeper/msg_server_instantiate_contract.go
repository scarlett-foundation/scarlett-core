package keeper

import (
	"context"

	"scarlett-core/x/contracts/types"

	errorsmod "cosmossdk.io/errors"
	wasmkeeper "github.com/CosmWasm/wasmd/x/wasm/keeper"
	wasmtypes "github.com/CosmWasm/wasmd/x/wasm/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k msgServer) InstantiateContract(ctx context.Context, msg *types.MsgInstantiateContract) (*types.MsgInstantiateContractResponse, error) {
	// Validate creator address
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return nil, errorsmod.Wrap(err, "invalid creator address")
	}

	// Create wasmd InstantiateContract message
	wasmMsg := &wasmtypes.MsgInstantiateContract{
		Sender: msg.Creator,
		Admin:  msg.Admin,
		CodeID: msg.CodeId,
		Label:  msg.Label,
		Msg:    msg.Msg,
		Funds:  nil, // TODO: Parse funds from string
	}

	// Get the wasm keeper and create a message server for it
	wasmKeeper := k.getWasmKeeper()
	wasmMsgServer := wasmkeeper.NewMsgServerImpl(wasmKeeper)

	// Delegate to wasmd's InstantiateContract handler
	wasmResponse, err := wasmMsgServer.InstantiateContract(ctx, wasmMsg)
	if err != nil {
		return nil, errorsmod.Wrap(err, "failed to instantiate contract")
	}

	// Log the successful contract instantiation
	k.Logger(ctx).Info("Successfully instantiated contract",
		"contract_address", wasmResponse.Address,
		"code_id", msg.CodeId,
		"creator", msg.Creator)

	// Return response with the contract address
	return &types.MsgInstantiateContractResponse{
		ContractAddress: wasmResponse.Address,
	}, nil
}

package keeper

import (
	"context"

	"scarlett-core/x/contracts/types"

	errorsmod "cosmossdk.io/errors"
	wasmkeeper "github.com/CosmWasm/wasmd/x/wasm/keeper"
	wasmtypes "github.com/CosmWasm/wasmd/x/wasm/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k msgServer) InstantiateContract(goCtx context.Context, msg *types.MsgInstantiateContract) (*types.MsgInstantiateContractResponse, error) {
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
	wasmMsg := &wasmtypes.MsgInstantiateContract{
		Sender: msg.Creator,
		CodeID: msg.CodeId,
		Label:  msg.Label,
		Msg:    msg.Msg,
		// Convert funds from string to sdk.Coins
		Funds: nil, // TODO: Parse funds string to sdk.Coins
		Admin: msg.Admin,
	}

	// Create a message server for wasmd
	wasmMsgServer := wasmkeeper.NewMsgServerImpl(k.wasmKeeper)

	// Delegate to wasmd's instantiate contract handler
	resp, err := wasmMsgServer.InstantiateContract(ctx, wasmMsg)
	if err != nil {
		return nil, err
	}

	// Log the successful contract instantiation
	k.Logger(ctx).Info("Successfully instantiated contract",
		"contract_address", resp.Address,
		"code_id", msg.CodeId,
		"creator", msg.Creator)

	return &types.MsgInstantiateContractResponse{
		ContractAddress: resp.Address,
	}, nil
}

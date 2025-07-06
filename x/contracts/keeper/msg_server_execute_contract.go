package keeper

import (
	"context"

	"scarlett-core/x/contracts/types"

	errorsmod "cosmossdk.io/errors"
)

func (k msgServer) ExecuteContract(ctx context.Context, msg *types.MsgExecuteContract) (*types.MsgExecuteContractResponse, error) {
	if _, err := k.addressCodec.StringToBytes(msg.Creator); err != nil {
		return nil, errorsmod.Wrap(err, "invalid authority address")
	}

	// TODO: Handle the message

	return &types.MsgExecuteContractResponse{}, nil
}

package keeper

import (
	"context"

	"scarlett-core/x/contracts/types"

	errorsmod "cosmossdk.io/errors"
)

func (k msgServer) InstantiateContract(ctx context.Context, msg *types.MsgInstantiateContract) (*types.MsgInstantiateContractResponse, error) {
	if _, err := k.addressCodec.StringToBytes(msg.Creator); err != nil {
		return nil, errorsmod.Wrap(err, "invalid authority address")
	}

	// TODO: Handle the message

	return &types.MsgInstantiateContractResponse{}, nil
}

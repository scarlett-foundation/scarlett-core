package keeper

import (
	"context"

	"scarlett-core/x/contracts/types"

	errorsmod "cosmossdk.io/errors"
)

func (k msgServer) StoreCode(ctx context.Context, msg *types.MsgStoreCode) (*types.MsgStoreCodeResponse, error) {
	if _, err := k.addressCodec.StringToBytes(msg.Creator); err != nil {
		return nil, errorsmod.Wrap(err, "invalid authority address")
	}

	// TODO: Handle the message

	return &types.MsgStoreCodeResponse{}, nil
}

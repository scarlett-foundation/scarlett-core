package keeper

import (
	"context"

	"scarlett-core/x/emissions/types"

	errorsmod "cosmossdk.io/errors"
)

func (k msgServer) UpdateEmissionSplit(ctx context.Context, msg *types.MsgUpdateEmissionSplit) (*types.MsgUpdateEmissionSplitResponse, error) {
	if _, err := k.addressCodec.StringToBytes(msg.Creator); err != nil {
		return nil, errorsmod.Wrap(err, "invalid authority address")
	}

	// TODO: Handle the message

	return &types.MsgUpdateEmissionSplitResponse{}, nil
}

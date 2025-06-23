package keeper

import (
	"context"

	"scarlett-core/x/scarlettcore/types"

	errorsmod "cosmossdk.io/errors"
)

func (k msgServer) BurnTokens(ctx context.Context, msg *types.MsgBurnTokens) (*types.MsgBurnTokensResponse, error) {
	if _, err := k.addressCodec.StringToBytes(msg.Creator); err != nil {
		return nil, errorsmod.Wrap(err, "invalid authority address")
	}

	// TODO: Handle the message

	return &types.MsgBurnTokensResponse{}, nil
}

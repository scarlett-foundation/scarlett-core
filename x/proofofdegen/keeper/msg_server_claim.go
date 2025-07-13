package keeper

import (
	"context"

	"scarlett-core/x/proofofdegen/types"

	errorsmod "cosmossdk.io/errors"
)

func (k msgServer) Claim(ctx context.Context, msg *types.MsgClaim) (*types.MsgClaimResponse, error) {
	if _, err := k.addressCodec.StringToBytes(msg.Creator); err != nil {
		return nil, errorsmod.Wrap(err, "invalid authority address")
	}

	// TODO: Handle the message

	return &types.MsgClaimResponse{}, nil
}

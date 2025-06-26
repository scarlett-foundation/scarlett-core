package keeper

import (
	"context"

	"scarlett-core/x/emissions/types"

	errorsmod "cosmossdk.io/errors"
)

func (k msgServer) RemoveEmissionDestination(ctx context.Context, msg *types.MsgRemoveEmissionDestination) (*types.MsgRemoveEmissionDestinationResponse, error) {
	if _, err := k.addressCodec.StringToBytes(msg.Creator); err != nil {
		return nil, errorsmod.Wrap(err, "invalid authority address")
	}

	// TODO: Handle the message

	return &types.MsgRemoveEmissionDestinationResponse{}, nil
}

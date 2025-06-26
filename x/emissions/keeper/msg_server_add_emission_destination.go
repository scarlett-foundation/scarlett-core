package keeper

import (
	"context"

	"scarlett-core/x/emissions/types"

	errorsmod "cosmossdk.io/errors"
)

func (k msgServer) AddEmissionDestination(ctx context.Context, msg *types.MsgAddEmissionDestination) (*types.MsgAddEmissionDestinationResponse, error) {
	if _, err := k.addressCodec.StringToBytes(msg.Creator); err != nil {
		return nil, errorsmod.Wrap(err, "invalid authority address")
	}

	// TODO: Handle the message

	return &types.MsgAddEmissionDestinationResponse{}, nil
}

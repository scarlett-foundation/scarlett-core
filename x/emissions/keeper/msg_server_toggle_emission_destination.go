package keeper

import (
	"context"

	"scarlett-core/x/emissions/types"

	errorsmod "cosmossdk.io/errors"
)

func (k msgServer) ToggleEmissionDestination(ctx context.Context, msg *types.MsgToggleEmissionDestination) (*types.MsgToggleEmissionDestinationResponse, error) {
	if _, err := k.addressCodec.StringToBytes(msg.Creator); err != nil {
		return nil, errorsmod.Wrap(err, "invalid authority address")
	}

	// TODO: Handle the message

	return &types.MsgToggleEmissionDestinationResponse{}, nil
}

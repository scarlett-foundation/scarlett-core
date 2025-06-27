package keeper

import (
	"bytes"
	"context"

	"scarlett-core/x/emissions/types"

	errorsmod "cosmossdk.io/errors"
)

func (k msgServer) RemoveEmissionDestination(ctx context.Context, msg *types.MsgRemoveEmissionDestination) (*types.MsgRemoveEmissionDestinationResponse, error) {
	// 1. Authority validation - only governance can remove emission destinations
	authority, err := k.addressCodec.StringToBytes(msg.Creator)
	if err != nil {
		return nil, errorsmod.Wrap(err, "invalid authority address")
	}

	if !bytes.Equal(k.GetAuthority(), authority) {
		expectedAuthorityStr, _ := k.addressCodec.BytesToString(k.GetAuthority())
		return nil, errorsmod.Wrapf(types.ErrUnauthorized, "invalid authority; expected %s, got %s", expectedAuthorityStr, msg.Creator)
	}

	// TODO: Handle the message - implement removing emission destination

	return &types.MsgRemoveEmissionDestinationResponse{}, nil
}

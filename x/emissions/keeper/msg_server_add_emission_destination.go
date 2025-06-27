package keeper

import (
	"bytes"
	"context"

	"scarlett-core/x/emissions/types"

	errorsmod "cosmossdk.io/errors"
)

func (k msgServer) AddEmissionDestination(ctx context.Context, msg *types.MsgAddEmissionDestination) (*types.MsgAddEmissionDestinationResponse, error) {
	// 1. Authority validation - only governance can add emission destinations
	authority, err := k.addressCodec.StringToBytes(msg.Creator)
	if err != nil {
		return nil, errorsmod.Wrap(err, "invalid authority address")
	}

	if !bytes.Equal(k.GetAuthority(), authority) {
		expectedAuthorityStr, _ := k.addressCodec.BytesToString(k.GetAuthority())
		return nil, errorsmod.Wrapf(types.ErrUnauthorized, "invalid authority; expected %s, got %s", expectedAuthorityStr, msg.Creator)
	}

	// TODO: Handle the message - implement adding new emission destination

	return &types.MsgAddEmissionDestinationResponse{}, nil
}

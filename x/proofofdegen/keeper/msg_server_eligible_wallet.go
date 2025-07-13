package keeper

import (
	"context"
	"errors"
	"fmt"

	"scarlett-core/x/proofofdegen/types"

	"cosmossdk.io/collections"
	errorsmod "cosmossdk.io/errors"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

func (k msgServer) CreateEligibleWallet(ctx context.Context, msg *types.MsgCreateEligibleWallet) (*types.MsgCreateEligibleWalletResponse, error) {
	if _, err := k.addressCodec.StringToBytes(msg.Creator); err != nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidAddress, fmt.Sprintf("invalid address: %s", err))
	}

	// Check if the value already exists
	ok, err := k.EligibleWallet.Has(ctx, msg.Index)
	if err != nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrLogic, err.Error())
	} else if ok {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "index already set")
	}

	var eligibleWallet = types.EligibleWallet{
		Creator:   msg.Creator,
		Index:     msg.Index,
		Address:   msg.Address,
		Claimed:   msg.Claimed,
		ClaimTime: msg.ClaimTime,
		Weight:    msg.Weight,
	}

	if err := k.EligibleWallet.Set(ctx, eligibleWallet.Index, eligibleWallet); err != nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrLogic, err.Error())
	}

	return &types.MsgCreateEligibleWalletResponse{}, nil
}

func (k msgServer) UpdateEligibleWallet(ctx context.Context, msg *types.MsgUpdateEligibleWallet) (*types.MsgUpdateEligibleWalletResponse, error) {
	if _, err := k.addressCodec.StringToBytes(msg.Creator); err != nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidAddress, fmt.Sprintf("invalid signer address: %s", err))
	}

	// Check if the value exists
	val, err := k.EligibleWallet.Get(ctx, msg.Index)
	if err != nil {
		if errors.Is(err, collections.ErrNotFound) {
			return nil, errorsmod.Wrap(sdkerrors.ErrKeyNotFound, "index not set")
		}

		return nil, errorsmod.Wrap(sdkerrors.ErrLogic, err.Error())
	}

	// Checks if the msg creator is the same as the current owner
	if msg.Creator != val.Creator {
		return nil, errorsmod.Wrap(sdkerrors.ErrUnauthorized, "incorrect owner")
	}

	var eligibleWallet = types.EligibleWallet{
		Creator:   msg.Creator,
		Index:     msg.Index,
		Address:   msg.Address,
		Claimed:   msg.Claimed,
		ClaimTime: msg.ClaimTime,
		Weight:    msg.Weight,
	}

	if err := k.EligibleWallet.Set(ctx, eligibleWallet.Index, eligibleWallet); err != nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrLogic, "failed to update eligibleWallet")
	}

	return &types.MsgUpdateEligibleWalletResponse{}, nil
}

func (k msgServer) DeleteEligibleWallet(ctx context.Context, msg *types.MsgDeleteEligibleWallet) (*types.MsgDeleteEligibleWalletResponse, error) {
	if _, err := k.addressCodec.StringToBytes(msg.Creator); err != nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidAddress, fmt.Sprintf("invalid signer address: %s", err))
	}

	// Check if the value exists
	val, err := k.EligibleWallet.Get(ctx, msg.Index)
	if err != nil {
		if errors.Is(err, collections.ErrNotFound) {
			return nil, errorsmod.Wrap(sdkerrors.ErrKeyNotFound, "index not set")
		}

		return nil, errorsmod.Wrap(sdkerrors.ErrLogic, err.Error())
	}

	// Checks if the msg creator is the same as the current owner
	if msg.Creator != val.Creator {
		return nil, errorsmod.Wrap(sdkerrors.ErrUnauthorized, "incorrect owner")
	}

	if err := k.EligibleWallet.Remove(ctx, msg.Index); err != nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrLogic, "failed to remove eligibleWallet")
	}

	return &types.MsgDeleteEligibleWalletResponse{}, nil
}

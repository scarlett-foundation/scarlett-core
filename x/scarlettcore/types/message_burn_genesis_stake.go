package types

import (
	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func NewMsgBurnGenesisStake(creator string, amount string) *MsgBurnGenesisStake {
	return &MsgBurnGenesisStake{
		Creator: creator,
		Amount:  amount,
	}
}

func (msg *MsgBurnGenesisStake) ValidateBasic() error {
	// Validate creator address
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return errorsmod.Wrapf(ErrInvalidAddress, "invalid creator address (%s)", err)
	}

	// Validate amount
	if msg.Amount == "" {
		return errorsmod.Wrap(ErrInvalidAmount, "amount cannot be empty")
	}

	// Parse and validate the amount
	coin, err := sdk.ParseCoinNormalized(msg.Amount)
	if err != nil {
		return errorsmod.Wrapf(ErrInvalidAmount, "invalid amount format: %s", err)
	}

	// Ensure amount is positive
	if !coin.Amount.IsPositive() {
		return errorsmod.Wrap(ErrInvalidAmount, "amount must be positive")
	}

	// Ensure denomination is sclt
	if coin.Denom != "sclt" {
		return errorsmod.Wrapf(ErrInvalidAmount, "invalid denomination, expected 'sclt', got '%s'", coin.Denom)
	}

	return nil
}

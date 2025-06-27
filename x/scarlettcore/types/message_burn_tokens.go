package types

func NewMsgBurnTokens(creator string, amount string, denom string) *MsgBurnTokens {
	return &MsgBurnTokens{
		Creator: creator,
		Amount:  amount,
		Denom:   denom,
	}
}

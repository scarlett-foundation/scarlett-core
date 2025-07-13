package types

func NewMsgClaim(creator string, address string) *MsgClaim {
	return &MsgClaim{
		Creator: creator,
		Address: address,
	}
}

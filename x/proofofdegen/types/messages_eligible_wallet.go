package types

func NewMsgCreateEligibleWallet(
	creator string,
	index string,
	address string,
	claimed bool,
	claimTime uint64,
	weight uint64,

) *MsgCreateEligibleWallet {
	return &MsgCreateEligibleWallet{
		Creator:   creator,
		Index:     index,
		Address:   address,
		Claimed:   claimed,
		ClaimTime: claimTime,
		Weight:    weight,
	}
}

func NewMsgUpdateEligibleWallet(
	creator string,
	index string,
	address string,
	claimed bool,
	claimTime uint64,
	weight uint64,

) *MsgUpdateEligibleWallet {
	return &MsgUpdateEligibleWallet{
		Creator:   creator,
		Index:     index,
		Address:   address,
		Claimed:   claimed,
		ClaimTime: claimTime,
		Weight:    weight,
	}
}

func NewMsgDeleteEligibleWallet(
	creator string,
	index string,

) *MsgDeleteEligibleWallet {
	return &MsgDeleteEligibleWallet{
		Creator: creator,
		Index:   index,
	}
}

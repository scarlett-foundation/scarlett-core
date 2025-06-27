package types

func NewMsgBurnGenesisStake(creator string) *MsgBurnGenesisStake {
	return &MsgBurnGenesisStake{
		Creator: creator,
	}
}

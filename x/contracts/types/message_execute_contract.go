package types

func NewMsgExecuteContract(creator string, contract string, msg []byte, funds string) *MsgExecuteContract {
	return &MsgExecuteContract{
		Creator:  creator,
		Contract: contract,
		Msg:      msg,
		Funds:    funds,
	}
}

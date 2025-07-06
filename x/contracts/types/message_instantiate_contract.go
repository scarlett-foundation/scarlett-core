package types

func NewMsgInstantiateContract(creator string, codeId uint64, msg []byte, label string, funds string, admin string) *MsgInstantiateContract {
	return &MsgInstantiateContract{
		Creator: creator,
		CodeId:  codeId,
		Msg:     msg,
		Label:   label,
		Funds:   funds,
		Admin:   admin,
	}
}

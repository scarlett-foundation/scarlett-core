package types

func NewMsgDeployContract(creator string, code []byte, instantiateMsg string, label string) *MsgDeployContract {
	return &MsgDeployContract{
		Creator:        creator,
		Code:           code,
		InstantiateMsg: instantiateMsg,
		Label:          label,
	}
}

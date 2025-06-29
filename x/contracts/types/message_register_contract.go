package types

func NewMsgRegisterContract(creator string, contractAddress string, name string, description string) *MsgRegisterContract {
	return &MsgRegisterContract{
		Creator:         creator,
		ContractAddress: contractAddress,
		Name:            name,
		Description:     description,
	}
}

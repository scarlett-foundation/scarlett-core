package types

func NewMsgRegisterModule(creator string, moduleName string, description string) *MsgRegisterModule {
	return &MsgRegisterModule{
		Creator:     creator,
		ModuleName:  moduleName,
		Description: description,
	}
}

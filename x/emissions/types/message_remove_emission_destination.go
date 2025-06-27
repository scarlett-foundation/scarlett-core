package types

func NewMsgRemoveEmissionDestination(creator string, module string, reason string) *MsgRemoveEmissionDestination {
	return &MsgRemoveEmissionDestination{
		Creator: creator,
		Module:  module,
		Reason:  reason,
	}
}

package types

func NewMsgToggleEmissionDestination(creator string, module string, enabled bool, reason string) *MsgToggleEmissionDestination {
	return &MsgToggleEmissionDestination{
		Creator: creator,
		Module:  module,
		Enabled: enabled,
		Reason:  reason,
	}
}

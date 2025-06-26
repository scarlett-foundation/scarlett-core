package types

func NewMsgAddEmissionDestination(creator string, module string, weight string, description string, reason string) *MsgAddEmissionDestination {
	return &MsgAddEmissionDestination{
		Creator:     creator,
		Module:      module,
		Weight:      weight,
		Description: description,
		Reason:      reason,
	}
}

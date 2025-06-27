package types

func NewMsgUpdateEmissionSplit(creator string, destinations []string, weights []string, reason string) *MsgUpdateEmissionSplit {
	return &MsgUpdateEmissionSplit{
		Creator:      creator,
		Destinations: destinations,
		Weights:      weights,
		Reason:       reason,
	}
}

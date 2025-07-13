package types

func NewMsgCreateCampaign(
	creator string,
	index string,
	name string,
	active bool,
	totalAllocation uint64,

) *MsgCreateCampaign {
	return &MsgCreateCampaign{
		Creator:         creator,
		Index:           index,
		Name:            name,
		Active:          active,
		TotalAllocation: totalAllocation,
	}
}

func NewMsgUpdateCampaign(
	creator string,
	index string,
	name string,
	active bool,
	totalAllocation uint64,

) *MsgUpdateCampaign {
	return &MsgUpdateCampaign{
		Creator:         creator,
		Index:           index,
		Name:            name,
		Active:          active,
		TotalAllocation: totalAllocation,
	}
}

func NewMsgDeleteCampaign(
	creator string,
	index string,

) *MsgDeleteCampaign {
	return &MsgDeleteCampaign{
		Creator: creator,
		Index:   index,
	}
}

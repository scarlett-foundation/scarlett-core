package types

import "fmt"

// DefaultGenesis returns the default genesis state
func DefaultGenesis() *GenesisState {
	return &GenesisState{
		Params:      DefaultParams(),
		CampaignMap: []Campaign{}, EligibleWalletMap: []EligibleWallet{}}
}

// Validate performs basic genesis state validation returning an error upon any
// failure.
func (gs GenesisState) Validate() error {
	campaignIndexMap := make(map[string]struct{})

	for _, elem := range gs.CampaignMap {
		index := fmt.Sprint(elem.Index)
		if _, ok := campaignIndexMap[index]; ok {
			return fmt.Errorf("duplicated index for campaign")
		}
		campaignIndexMap[index] = struct{}{}
	}
	eligibleWalletIndexMap := make(map[string]struct{})

	for _, elem := range gs.EligibleWalletMap {
		index := fmt.Sprint(elem.Index)
		if _, ok := eligibleWalletIndexMap[index]; ok {
			return fmt.Errorf("duplicated index for eligibleWallet")
		}
		eligibleWalletIndexMap[index] = struct{}{}
	}

	return gs.Params.Validate()
}

package types

import "fmt"

// DefaultGenesis returns the default genesis state
func DefaultGenesis() *GenesisState {
	return &GenesisState{
		Params: DefaultParams(),
		CampaignMap: []Campaign{
			{
				Index:           "genesis", // Unique index for genesis campaign
				Name:            "Genesis",
				Active:          true,
				TotalAllocation: 1000000, // 1M SCLT allocated to Genesis campaign
			},
		},
		EligibleWalletMap: []EligibleWallet{
			// Test eligible wallets with different weights for testing patience mechanics
			{
				Index:     "alice",                                           // Unique index using nickname
				Address:   "scarlett1tpsw9anauv2peeurtudrjs8rq8ngjpvc45a0zl", // Alice
				Claimed:   false,
				ClaimTime: 0,
				Weight:    150, // Higher weight for early contributor
			},
			{
				Index:     "bob",                                             // Unique index using nickname
				Address:   "scarlett1estg4c5xtq37xhsxtwd5gqtppnhpa7s607gt0d", // Bob
				Claimed:   false,
				ClaimTime: 0,
				Weight:    100, // Standard weight
			},
			{
				Index:     "validator1",                                      // Unique index using nickname
				Address:   "scarlett1v55ru2vwhrdz6jy8264dh43p43knas24tmtv45", // Validator1
				Claimed:   false,
				ClaimTime: 0,
				Weight:    200, // Higher weight for validator
			},
			{
				Index:     "validator2",                                      // Unique index using nickname
				Address:   "scarlett16e8hn309ntektqcrjgznq5u0w2azt2w93e7ptc", // Validator2
				Claimed:   false,
				ClaimTime: 0,
				Weight:    75, // Lower weight (late joiner)
			},
			{
				Index:     "validator3",                                      // Unique index using nickname
				Address:   "scarlett1583vwdqgqkr3jht8pzquz988d2x52sfqrdmwu8", // Validator3
				Claimed:   false,
				ClaimTime: 0,
				Weight:    125, // Medium weight
			},
		},
		// this line is used by starport scaffolding # genesis/types/default
	}
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

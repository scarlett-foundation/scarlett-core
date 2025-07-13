package proofofdegen

import (
	"math/rand"

	"github.com/cosmos/cosmos-sdk/types/module"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	"github.com/cosmos/cosmos-sdk/x/simulation"

	"scarlett-core/testutil/sample"
	proofofdegensimulation "scarlett-core/x/proofofdegen/simulation"
	"scarlett-core/x/proofofdegen/types"
)

// GenerateGenesisState creates a randomized GenState of the module.
func (AppModule) GenerateGenesisState(simState *module.SimulationState) {
	accs := make([]string, len(simState.Accounts))
	for i, acc := range simState.Accounts {
		accs[i] = acc.Address.String()
	}
	proofofdegenGenesis := types.GenesisState{
		Params: types.DefaultParams(),
		CampaignMap: []types.Campaign{{Creator: sample.AccAddress(),
			Index: "0",
		}, {Creator: sample.AccAddress(),
			Index: "1",
		}}, EligibleWalletMap: []types.EligibleWallet{{Creator: sample.AccAddress(),
			Index: "0",
		}, {Creator: sample.AccAddress(),
			Index: "1",
		}}}
	simState.GenState[types.ModuleName] = simState.Cdc.MustMarshalJSON(&proofofdegenGenesis)
}

// RegisterStoreDecoder registers a decoder.
func (am AppModule) RegisterStoreDecoder(_ simtypes.StoreDecoderRegistry) {}

// WeightedOperations returns the all the gov module operations with their respective weights.
func (am AppModule) WeightedOperations(simState module.SimulationState) []simtypes.WeightedOperation {
	operations := make([]simtypes.WeightedOperation, 0)
	const (
		opWeightMsgCreateCampaign          = "op_weight_msg_proofofdegen"
		defaultWeightMsgCreateCampaign int = 100
	)

	var weightMsgCreateCampaign int
	simState.AppParams.GetOrGenerate(opWeightMsgCreateCampaign, &weightMsgCreateCampaign, nil,
		func(_ *rand.Rand) {
			weightMsgCreateCampaign = defaultWeightMsgCreateCampaign
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgCreateCampaign,
		proofofdegensimulation.SimulateMsgCreateCampaign(am.authKeeper, am.bankKeeper, am.keeper, simState.TxConfig),
	))
	const (
		opWeightMsgUpdateCampaign          = "op_weight_msg_proofofdegen"
		defaultWeightMsgUpdateCampaign int = 100
	)

	var weightMsgUpdateCampaign int
	simState.AppParams.GetOrGenerate(opWeightMsgUpdateCampaign, &weightMsgUpdateCampaign, nil,
		func(_ *rand.Rand) {
			weightMsgUpdateCampaign = defaultWeightMsgUpdateCampaign
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgUpdateCampaign,
		proofofdegensimulation.SimulateMsgUpdateCampaign(am.authKeeper, am.bankKeeper, am.keeper, simState.TxConfig),
	))
	const (
		opWeightMsgDeleteCampaign          = "op_weight_msg_proofofdegen"
		defaultWeightMsgDeleteCampaign int = 100
	)

	var weightMsgDeleteCampaign int
	simState.AppParams.GetOrGenerate(opWeightMsgDeleteCampaign, &weightMsgDeleteCampaign, nil,
		func(_ *rand.Rand) {
			weightMsgDeleteCampaign = defaultWeightMsgDeleteCampaign
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgDeleteCampaign,
		proofofdegensimulation.SimulateMsgDeleteCampaign(am.authKeeper, am.bankKeeper, am.keeper, simState.TxConfig),
	))
	const (
		opWeightMsgCreateEligibleWallet          = "op_weight_msg_proofofdegen"
		defaultWeightMsgCreateEligibleWallet int = 100
	)

	var weightMsgCreateEligibleWallet int
	simState.AppParams.GetOrGenerate(opWeightMsgCreateEligibleWallet, &weightMsgCreateEligibleWallet, nil,
		func(_ *rand.Rand) {
			weightMsgCreateEligibleWallet = defaultWeightMsgCreateEligibleWallet
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgCreateEligibleWallet,
		proofofdegensimulation.SimulateMsgCreateEligibleWallet(am.authKeeper, am.bankKeeper, am.keeper, simState.TxConfig),
	))
	const (
		opWeightMsgUpdateEligibleWallet          = "op_weight_msg_proofofdegen"
		defaultWeightMsgUpdateEligibleWallet int = 100
	)

	var weightMsgUpdateEligibleWallet int
	simState.AppParams.GetOrGenerate(opWeightMsgUpdateEligibleWallet, &weightMsgUpdateEligibleWallet, nil,
		func(_ *rand.Rand) {
			weightMsgUpdateEligibleWallet = defaultWeightMsgUpdateEligibleWallet
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgUpdateEligibleWallet,
		proofofdegensimulation.SimulateMsgUpdateEligibleWallet(am.authKeeper, am.bankKeeper, am.keeper, simState.TxConfig),
	))
	const (
		opWeightMsgDeleteEligibleWallet          = "op_weight_msg_proofofdegen"
		defaultWeightMsgDeleteEligibleWallet int = 100
	)

	var weightMsgDeleteEligibleWallet int
	simState.AppParams.GetOrGenerate(opWeightMsgDeleteEligibleWallet, &weightMsgDeleteEligibleWallet, nil,
		func(_ *rand.Rand) {
			weightMsgDeleteEligibleWallet = defaultWeightMsgDeleteEligibleWallet
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgDeleteEligibleWallet,
		proofofdegensimulation.SimulateMsgDeleteEligibleWallet(am.authKeeper, am.bankKeeper, am.keeper, simState.TxConfig),
	))

	return operations
}

// ProposalMsgs returns msgs used for governance proposals for simulations.
func (am AppModule) ProposalMsgs(simState module.SimulationState) []simtypes.WeightedProposalMsg {
	return []simtypes.WeightedProposalMsg{}
}

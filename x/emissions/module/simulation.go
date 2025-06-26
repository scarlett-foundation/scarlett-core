package emissions

import (
	"math/rand"

	"github.com/cosmos/cosmos-sdk/types/module"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	"github.com/cosmos/cosmos-sdk/x/simulation"

	emissionssimulation "scarlett-core/x/emissions/simulation"
	"scarlett-core/x/emissions/types"
)

// GenerateGenesisState creates a randomized GenState of the module.
func (AppModule) GenerateGenesisState(simState *module.SimulationState) {
	accs := make([]string, len(simState.Accounts))
	for i, acc := range simState.Accounts {
		accs[i] = acc.Address.String()
	}
	emissionsGenesis := types.GenesisState{
		Params: types.DefaultParams(),
	}
	simState.GenState[types.ModuleName] = simState.Cdc.MustMarshalJSON(&emissionsGenesis)
}

// RegisterStoreDecoder registers a decoder.
func (am AppModule) RegisterStoreDecoder(_ simtypes.StoreDecoderRegistry) {}

// WeightedOperations returns the all the gov module operations with their respective weights.
func (am AppModule) WeightedOperations(simState module.SimulationState) []simtypes.WeightedOperation {
	operations := make([]simtypes.WeightedOperation, 0)
	const (
		opWeightMsgUpdateEmissionSplit          = "op_weight_msg_emissions"
		defaultWeightMsgUpdateEmissionSplit int = 100
	)

	var weightMsgUpdateEmissionSplit int
	simState.AppParams.GetOrGenerate(opWeightMsgUpdateEmissionSplit, &weightMsgUpdateEmissionSplit, nil,
		func(_ *rand.Rand) {
			weightMsgUpdateEmissionSplit = defaultWeightMsgUpdateEmissionSplit
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgUpdateEmissionSplit,
		emissionssimulation.SimulateMsgUpdateEmissionSplit(am.authKeeper, am.bankKeeper, am.keeper, simState.TxConfig),
	))
	const (
		opWeightMsgAddEmissionDestination          = "op_weight_msg_emissions"
		defaultWeightMsgAddEmissionDestination int = 100
	)

	var weightMsgAddEmissionDestination int
	simState.AppParams.GetOrGenerate(opWeightMsgAddEmissionDestination, &weightMsgAddEmissionDestination, nil,
		func(_ *rand.Rand) {
			weightMsgAddEmissionDestination = defaultWeightMsgAddEmissionDestination
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgAddEmissionDestination,
		emissionssimulation.SimulateMsgAddEmissionDestination(am.authKeeper, am.bankKeeper, am.keeper, simState.TxConfig),
	))

	return operations
}

// ProposalMsgs returns msgs used for governance proposals for simulations.
func (am AppModule) ProposalMsgs(simState module.SimulationState) []simtypes.WeightedProposalMsg {
	return []simtypes.WeightedProposalMsg{}
}

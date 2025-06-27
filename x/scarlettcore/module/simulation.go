package scarlettcore

import (
	"math/rand"

	"github.com/cosmos/cosmos-sdk/types/module"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	"github.com/cosmos/cosmos-sdk/x/simulation"

	scarlettcoresimulation "scarlett-core/x/scarlettcore/simulation"
	"scarlett-core/x/scarlettcore/types"
)

// GenerateGenesisState creates a randomized GenState of the module.
func (AppModule) GenerateGenesisState(simState *module.SimulationState) {
	accs := make([]string, len(simState.Accounts))
	for i, acc := range simState.Accounts {
		accs[i] = acc.Address.String()
	}
	scarlettcoreGenesis := types.GenesisState{
		Params: types.DefaultParams(),
	}
	simState.GenState[types.ModuleName] = simState.Cdc.MustMarshalJSON(&scarlettcoreGenesis)
}

// RegisterStoreDecoder registers a decoder.
func (am AppModule) RegisterStoreDecoder(_ simtypes.StoreDecoderRegistry) {}

// WeightedOperations returns the all the gov module operations with their respective weights.
func (am AppModule) WeightedOperations(simState module.SimulationState) []simtypes.WeightedOperation {
	operations := make([]simtypes.WeightedOperation, 0)
	const (
		opWeightMsgBurnTokens          = "op_weight_msg_scarlettcore"
		defaultWeightMsgBurnTokens int = 100
	)

	var weightMsgBurnTokens int
	simState.AppParams.GetOrGenerate(opWeightMsgBurnTokens, &weightMsgBurnTokens, nil,
		func(_ *rand.Rand) {
			weightMsgBurnTokens = defaultWeightMsgBurnTokens
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgBurnTokens,
		scarlettcoresimulation.SimulateMsgBurnTokens(am.authKeeper, am.bankKeeper, am.keeper, simState.TxConfig),
	))
	const (
		opWeightMsgBurnGenesisStake          = "op_weight_msg_scarlettcore"
		defaultWeightMsgBurnGenesisStake int = 100
	)

	var weightMsgBurnGenesisStake int
	simState.AppParams.GetOrGenerate(opWeightMsgBurnGenesisStake, &weightMsgBurnGenesisStake, nil,
		func(_ *rand.Rand) {
			weightMsgBurnGenesisStake = defaultWeightMsgBurnGenesisStake
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgBurnGenesisStake,
		scarlettcoresimulation.SimulateMsgBurnGenesisStake(am.authKeeper, am.bankKeeper, am.keeper, simState.TxConfig),
	))

	return operations
}

// ProposalMsgs returns msgs used for governance proposals for simulations.
func (am AppModule) ProposalMsgs(simState module.SimulationState) []simtypes.WeightedProposalMsg {
	return []simtypes.WeightedProposalMsg{}
}

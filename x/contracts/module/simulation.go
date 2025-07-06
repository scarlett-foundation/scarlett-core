package contracts

import (
	"math/rand"

	"github.com/cosmos/cosmos-sdk/types/module"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	"github.com/cosmos/cosmos-sdk/x/simulation"

	contractssimulation "scarlett-core/x/contracts/simulation"
	"scarlett-core/x/contracts/types"
)

// GenerateGenesisState creates a randomized GenState of the module.
func (AppModule) GenerateGenesisState(simState *module.SimulationState) {
	accs := make([]string, len(simState.Accounts))
	for i, acc := range simState.Accounts {
		accs[i] = acc.Address.String()
	}
	contractsGenesis := types.GenesisState{
		Params: types.DefaultParams(),
	}
	simState.GenState[types.ModuleName] = simState.Cdc.MustMarshalJSON(&contractsGenesis)
}

// RegisterStoreDecoder registers a decoder.
func (am AppModule) RegisterStoreDecoder(_ simtypes.StoreDecoderRegistry) {}

// WeightedOperations returns the all the gov module operations with their respective weights.
func (am AppModule) WeightedOperations(simState module.SimulationState) []simtypes.WeightedOperation {
	operations := make([]simtypes.WeightedOperation, 0)
	const (
		opWeightMsgStoreCode          = "op_weight_msg_contracts"
		defaultWeightMsgStoreCode int = 100
	)

	var weightMsgStoreCode int
	simState.AppParams.GetOrGenerate(opWeightMsgStoreCode, &weightMsgStoreCode, nil,
		func(_ *rand.Rand) {
			weightMsgStoreCode = defaultWeightMsgStoreCode
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgStoreCode,
		contractssimulation.SimulateMsgStoreCode(am.authKeeper, am.bankKeeper, am.keeper, simState.TxConfig),
	))
	const (
		opWeightMsgInstantiateContract          = "op_weight_msg_contracts"
		defaultWeightMsgInstantiateContract int = 100
	)

	var weightMsgInstantiateContract int
	simState.AppParams.GetOrGenerate(opWeightMsgInstantiateContract, &weightMsgInstantiateContract, nil,
		func(_ *rand.Rand) {
			weightMsgInstantiateContract = defaultWeightMsgInstantiateContract
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgInstantiateContract,
		contractssimulation.SimulateMsgInstantiateContract(am.authKeeper, am.bankKeeper, am.keeper, simState.TxConfig),
	))
	const (
		opWeightMsgExecuteContract          = "op_weight_msg_contracts"
		defaultWeightMsgExecuteContract int = 100
	)

	var weightMsgExecuteContract int
	simState.AppParams.GetOrGenerate(opWeightMsgExecuteContract, &weightMsgExecuteContract, nil,
		func(_ *rand.Rand) {
			weightMsgExecuteContract = defaultWeightMsgExecuteContract
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgExecuteContract,
		contractssimulation.SimulateMsgExecuteContract(am.authKeeper, am.bankKeeper, am.keeper, simState.TxConfig),
	))

	return operations
}

// ProposalMsgs returns msgs used for governance proposals for simulations.
func (am AppModule) ProposalMsgs(simState module.SimulationState) []simtypes.WeightedProposalMsg {
	return []simtypes.WeightedProposalMsg{}
}

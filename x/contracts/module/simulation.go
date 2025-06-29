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
		opWeightMsgDeployContract          = "op_weight_msg_contracts"
		defaultWeightMsgDeployContract int = 100
	)

	var weightMsgDeployContract int
	simState.AppParams.GetOrGenerate(opWeightMsgDeployContract, &weightMsgDeployContract, nil,
		func(_ *rand.Rand) {
			weightMsgDeployContract = defaultWeightMsgDeployContract
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgDeployContract,
		contractssimulation.SimulateMsgDeployContract(am.authKeeper, am.bankKeeper, am.keeper, simState.TxConfig),
	))
	const (
		opWeightMsgRegisterContract          = "op_weight_msg_contracts"
		defaultWeightMsgRegisterContract int = 100
	)

	var weightMsgRegisterContract int
	simState.AppParams.GetOrGenerate(opWeightMsgRegisterContract, &weightMsgRegisterContract, nil,
		func(_ *rand.Rand) {
			weightMsgRegisterContract = defaultWeightMsgRegisterContract
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgRegisterContract,
		contractssimulation.SimulateMsgRegisterContract(am.authKeeper, am.bankKeeper, am.keeper, simState.TxConfig),
	))

	return operations
}

// ProposalMsgs returns msgs used for governance proposals for simulations.
func (am AppModule) ProposalMsgs(simState module.SimulationState) []simtypes.WeightedProposalMsg {
	return []simtypes.WeightedProposalMsg{}
}

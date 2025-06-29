package contracts

import (
	autocliv1 "cosmossdk.io/api/cosmos/autocli/v1"

	"scarlett-core/x/contracts/types"
)

// AutoCLIOptions implements the autocli.HasAutoCLIConfig interface.
func (am AppModule) AutoCLIOptions() *autocliv1.ModuleOptions {
	return &autocliv1.ModuleOptions{
		Query: &autocliv1.ServiceCommandDescriptor{
			Service: types.Query_serviceDesc.ServiceName,
			RpcCommandOptions: []*autocliv1.RpcCommandOptions{
				{
					RpcMethod: "Params",
					Use:       "params",
					Short:     "Shows the parameters of the module",
				},
				{
					RpcMethod:      "ListContracts",
					Use:            "list-contracts ",
					Short:          "Query list-contracts",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{},
				},

				{
					RpcMethod:      "ContractInfo",
					Use:            "contract-info [contract-address]",
					Short:          "Query contract-info",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "contract_address"}},
				},

				// this line is used by ignite scaffolding # autocli/query
			},
		},
		Tx: &autocliv1.ServiceCommandDescriptor{
			Service:              types.Msg_serviceDesc.ServiceName,
			EnhanceCustomCommand: true, // only required if you want to use the custom command
			RpcCommandOptions: []*autocliv1.RpcCommandOptions{
				{
					RpcMethod: "UpdateParams",
					Skip:      true, // skipped because authority gated
				},
				{
					RpcMethod:      "DeployContract",
					Use:            "deploy-contract [code] [instantiate-msg] [label]",
					Short:          "Send a deploy-contract tx",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "code"}, {ProtoField: "instantiate_msg"}, {ProtoField: "label"}},
				},
				{
					RpcMethod:      "RegisterContract",
					Use:            "register-contract [contract-address] [name] [description]",
					Short:          "Send a register-contract tx",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "contract_address"}, {ProtoField: "name"}, {ProtoField: "description"}},
				},
				// this line is used by ignite scaffolding # autocli/tx
			},
		},
	}
}

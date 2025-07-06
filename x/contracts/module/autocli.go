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
				// TODO: Add wasm query commands once we implement the proto methods
				{
					RpcMethod:      "ContractInfo",
					Use:            "contract-info [contract]",
					Short:          "Query contract-info",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "contract"}},
				},

				{
					RpcMethod:      "SmartContractState",
					Use:            "smart-contract-state [contract] [query-data]",
					Short:          "Query smart-contract-state",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "contract"}, {ProtoField: "query_data", Varargs: true}},
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
				// TODO: Add wasm transaction commands once we implement the proto methods
				{
					RpcMethod:      "StoreCode",
					Use:            "store-code [wasm-byte-code]",
					Short:          "Send a store-code tx",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "wasm_byte_code", Varargs: true}},
				},
				{
					RpcMethod:      "InstantiateContract",
					Use:            "instantiate-contract [code-id] [msg] [label] [funds] [admin]",
					Short:          "Send a instantiate-contract tx",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "code_id"}, {ProtoField: "msg"}, {ProtoField: "label"}, {ProtoField: "funds"}, {ProtoField: "admin"}},
				},
				{
					RpcMethod:      "ExecuteContract",
					Use:            "execute-contract [contract] [msg] [funds]",
					Short:          "Send a execute-contract tx",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "contract"}, {ProtoField: "msg"}, {ProtoField: "funds"}},
				},
				// this line is used by ignite scaffolding # autocli/tx
			},
		},
	}
}

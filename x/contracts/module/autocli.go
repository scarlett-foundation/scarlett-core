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
				// Wasm query commands - expose the standard wasm queries
				{
					RpcMethod: "ContractInfo",
					Use:       "contract-info [contract-address]",
					Short:     "Prints out metadata of a contract given its address",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{
						{ProtoField: "address"},
					},
				},
				{
					RpcMethod: "ContractHistory",
					Use:       "contract-history [contract-address]",
					Short:     "Prints out the code history for a contract given its address",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{
						{ProtoField: "address"},
					},
				},
				{
					RpcMethod: "ContractsByCode",
					Use:       "list-contract-by-code [code-id]",
					Short:     "List all contracts by code ID",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{
						{ProtoField: "code_id"},
					},
				},
				{
					RpcMethod: "AllContractState",
					Use:       "contract-state-all [contract-address]",
					Short:     "Prints out all internal state of a contract given its address",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{
						{ProtoField: "address"},
					},
				},
				{
					RpcMethod: "RawContractState",
					Use:       "contract-state-raw [contract-address] [key]",
					Short:     "Prints out internal state for key of a contract given its address",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{
						{ProtoField: "address"},
						{ProtoField: "query_data"},
					},
				},
				{
					RpcMethod: "SmartContractState",
					Use:       "contract-state-smart [contract-address] [query]",
					Short:     "Calls contract with given address with query data and prints the returned result",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{
						{ProtoField: "address"},
						{ProtoField: "query_data"},
					},
				},
				{
					RpcMethod: "Code",
					Use:       "code [code-id]",
					Short:     "Downloads wasm bytecode for given code ID",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{
						{ProtoField: "code_id"},
					},
				},
				{
					RpcMethod: "Codes",
					Use:       "list-code",
					Short:     "List all wasm bytecode on the chain",
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
				// Wasm transaction commands - expose the standard wasm transactions
				{
					RpcMethod: "StoreCode",
					Use:       "store [wasm-file]",
					Short:     "Upload a wasm binary",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{
						{ProtoField: "wasm_byte_code"},
					},
				},
				{
					RpcMethod: "InstantiateContract",
					Use:       "instantiate [code-id] [json-encoded-init-args]",
					Short:     "Instantiate a wasm contract",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{
						{ProtoField: "code_id"},
						{ProtoField: "msg"},
					},
				},
				{
					RpcMethod: "ExecuteContract",
					Use:       "execute [contract-address] [json-encoded-send-args]",
					Short:     "Execute a command on a wasm contract",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{
						{ProtoField: "contract"},
						{ProtoField: "msg"},
					},
				},
				{
					RpcMethod: "MigrateContract",
					Use:       "migrate [contract-address] [new-code-id] [json-encoded-migration-args]",
					Short:     "Migrate a wasm contract to a new code version",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{
						{ProtoField: "contract"},
						{ProtoField: "code_id"},
						{ProtoField: "msg"},
					},
				},
				{
					RpcMethod: "UpdateAdmin",
					Use:       "set-contract-admin [contract-address] [new-admin]",
					Short:     "Set new admin for a contract",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{
						{ProtoField: "contract"},
						{ProtoField: "new_admin"},
					},
				},
				{
					RpcMethod: "ClearAdmin",
					Use:       "clear-contract-admin [contract-address]",
					Short:     "Clears admin for a contract to prevent further migrations",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{
						{ProtoField: "contract"},
					},
				},
				// this line is used by ignite scaffolding # autocli/tx
			},
		},
	}
}

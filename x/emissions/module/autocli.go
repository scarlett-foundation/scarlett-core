package emissions

import (
	autocliv1 "cosmossdk.io/api/cosmos/autocli/v1"

	"scarlett-core/x/emissions/types"
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
					RpcMethod:      "UpdateEmissionSplit",
					Use:            "update-emission-split [destinations] [weights] [reason]",
					Short:          "Send a update-emission-split tx",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "destinations"}, {ProtoField: "weights"}, {ProtoField: "reason"}},
				},
				{
					RpcMethod:      "AddEmissionDestination",
					Use:            "add-emission-destination [module] [weight] [description] [reason]",
					Short:          "Send a add-emission-destination tx",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "module"}, {ProtoField: "weight"}, {ProtoField: "description"}, {ProtoField: "reason"}},
				},
				{
					RpcMethod:      "RemoveEmissionDestination",
					Use:            "remove-emission-destination [module] [reason]",
					Short:          "Send a remove-emission-destination tx",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "module"}, {ProtoField: "reason"}},
				},
				{
					RpcMethod:      "ToggleEmissionDestination",
					Use:            "toggle-emission-destination [module] [enabled] [reason]",
					Short:          "Send a toggle-emission-destination tx",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "module"}, {ProtoField: "enabled"}, {ProtoField: "reason"}},
				},
				// this line is used by ignite scaffolding # autocli/tx
			},
		},
	}
}

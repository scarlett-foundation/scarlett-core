package proofofdegen

import (
	autocliv1 "cosmossdk.io/api/cosmos/autocli/v1"

	"scarlett-core/x/proofofdegen/types"
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
					RpcMethod: "ListCampaign",
					Use:       "list-campaign",
					Short:     "List all campaign",
				},
				{
					RpcMethod:      "GetCampaign",
					Use:            "get-campaign [id]",
					Short:          "Gets a campaign",
					Alias:          []string{"show-campaign"},
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "index"}},
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
					RpcMethod:      "CreateCampaign",
					Use:            "create-campaign [index] [name] [active] [total-allocation]",
					Short:          "Create a new campaign",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "index"}, {ProtoField: "name"}, {ProtoField: "active"}, {ProtoField: "total_allocation"}},
				},
				{
					RpcMethod:      "UpdateCampaign",
					Use:            "update-campaign [index] [name] [active] [total-allocation]",
					Short:          "Update campaign",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "index"}, {ProtoField: "name"}, {ProtoField: "active"}, {ProtoField: "total_allocation"}},
				},
				{
					RpcMethod:      "DeleteCampaign",
					Use:            "delete-campaign [index]",
					Short:          "Delete campaign",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "index"}},
				},
				// this line is used by ignite scaffolding # autocli/tx
			},
		},
	}
}

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
				{
					RpcMethod: "ListEligibleWallet",
					Use:       "list-eligible-wallet",
					Short:     "List all eligible-wallet",
				},
				{
					RpcMethod:      "GetEligibleWallet",
					Use:            "get-eligible-wallet [id]",
					Short:          "Gets a eligible-wallet",
					Alias:          []string{"show-eligible-wallet"},
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "index"}},
				},
				{
					RpcMethod:      "EligibleAmount",
					Use:            "eligible-amount [address]",
					Short:          "Query eligible-amount",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "address"}},
				},

				{
					RpcMethod:      "CampaignInfo",
					Use:            "campaign-info ",
					Short:          "Query campaign-info",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{},
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
				{
					RpcMethod:      "CreateEligibleWallet",
					Use:            "create-eligible-wallet [index] [address] [claimed] [claim-time] [weight]",
					Short:          "Create a new eligible-wallet",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "index"}, {ProtoField: "address"}, {ProtoField: "claimed"}, {ProtoField: "claim_time"}, {ProtoField: "weight"}},
				},
				{
					RpcMethod:      "UpdateEligibleWallet",
					Use:            "update-eligible-wallet [index] [address] [claimed] [claim-time] [weight]",
					Short:          "Update eligible-wallet",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "index"}, {ProtoField: "address"}, {ProtoField: "claimed"}, {ProtoField: "claim_time"}, {ProtoField: "weight"}},
				},
				{
					RpcMethod:      "DeleteEligibleWallet",
					Use:            "delete-eligible-wallet [index]",
					Short:          "Delete eligible-wallet",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "index"}},
				},
				{
					RpcMethod:      "Claim",
					Use:            "claim [address]",
					Short:          "Send a claim tx",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "address"}},
				},
				// this line is used by ignite scaffolding # autocli/tx
			},
		},
	}
}

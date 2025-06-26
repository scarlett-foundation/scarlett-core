package keeper

import (
	"fmt"

	"cosmossdk.io/collections"
	"cosmossdk.io/core/address"
	corestore "cosmossdk.io/core/store"
	"github.com/cosmos/cosmos-sdk/codec"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	govkeeper "github.com/cosmos/cosmos-sdk/x/gov/keeper"
	mintkeeper "github.com/cosmos/cosmos-sdk/x/mint/keeper"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"

	"scarlett-core/x/emissions/types"
)

type Keeper struct {
	storeService corestore.KVStoreService
	cdc          codec.Codec
	addressCodec address.Codec
	// Address capable of executing a MsgUpdateParams message.
	// Typically, this should be the x/gov module account.
	authority []byte

	// Keeper dependencies for emission distribution
	bankKeeper    bankkeeper.Keeper
	govKeeper     govkeeper.Keeper
	stakingKeeper stakingkeeper.Keeper
	mintKeeper    mintkeeper.Keeper

	// Collections for state management
	Schema             collections.Schema
	Params             collections.Item[types.Params]
	EmissionParams     collections.Item[types.EmissionParams]
	EmissionHistory    collections.Map[int64, types.EmissionParams]    // Block height -> params
	DestinationMetrics collections.Map[string, types.DestinationStats] // Module -> stats
}

func NewKeeper(
	storeService corestore.KVStoreService,
	cdc codec.Codec,
	addressCodec address.Codec,
	authority []byte,
	bankKeeper bankkeeper.Keeper,
	govKeeper govkeeper.Keeper,
	stakingKeeper stakingkeeper.Keeper,
	mintKeeper mintkeeper.Keeper,
) Keeper {
	if _, err := addressCodec.BytesToString(authority); err != nil {
		panic(fmt.Sprintf("invalid authority address %s: %s", authority, err))
	}

	sb := collections.NewSchemaBuilder(storeService)

	k := Keeper{
		storeService:  storeService,
		cdc:           cdc,
		addressCodec:  addressCodec,
		authority:     authority,
		bankKeeper:    bankKeeper,
		govKeeper:     govKeeper,
		stakingKeeper: stakingKeeper,
		mintKeeper:    mintKeeper,

		// Initialize collections
		Params:             collections.NewItem(sb, types.ParamsKey, "params", codec.CollValue[types.Params](cdc)),
		EmissionParams:     collections.NewItem(sb, types.EmissionParamsKey, "emission_params", codec.CollValue[types.EmissionParams](cdc)),
		EmissionHistory:    collections.NewMap(sb, types.EmissionHistoryPrefix, "emission_history", collections.Int64Key, codec.CollValue[types.EmissionParams](cdc)),
		DestinationMetrics: collections.NewMap(sb, types.DestinationMetricsPrefix, "destination_metrics", collections.StringKey, codec.CollValue[types.DestinationStats](cdc)),
	}

	schema, err := sb.Build()
	if err != nil {
		panic(err)
	}
	k.Schema = schema

	return k
}

// GetAuthority returns the module's authority.
func (k Keeper) GetAuthority() []byte {
	return k.authority
}

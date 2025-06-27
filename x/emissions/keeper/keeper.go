package keeper

import (
	"context"
	"encoding/json"
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
	Schema collections.Schema
	Params collections.Item[types.Params]
	// Temporarily use string storage for custom types until proto implementation
	EmissionParamsRaw     collections.Item[string]
	EmissionHistoryRaw    collections.Map[int64, string]
	DestinationMetricsRaw collections.Map[string, string]
}

func NewKeeper(
	storeService corestore.KVStoreService,
	cdc codec.Codec,
	addressCodec address.Codec,
	authority []byte,
	bankKeeper bankkeeper.Keeper,
	govKeeper govkeeper.Keeper,
	stakingKeeper stakingkeeper.Keeper,
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

		Params:                collections.NewItem(sb, types.ParamsKey, "params", codec.CollValue[types.Params](cdc)),
		EmissionParamsRaw:     collections.NewItem(sb, types.EmissionParamsKey, "emission_params", collections.StringValue),
		EmissionHistoryRaw:    collections.NewMap(sb, types.EmissionHistoryPrefix, "emission_history", collections.Int64Key, collections.StringValue),
		DestinationMetricsRaw: collections.NewMap(sb, types.DestinationMetricsPrefix, "destination_metrics", collections.StringKey, collections.StringValue),
	}

	schema, err := sb.Build()
	if err != nil {
		panic(err)
	}
	k.Schema = schema

	return k
}

// SetMintKeeper sets the mint keeper for the emissions keeper.
func (k *Keeper) SetMintKeeper(mk mintkeeper.Keeper) {
	k.mintKeeper = mk
}

// GetAuthority returns the module's authority.
func (k Keeper) GetAuthority() []byte {
	return k.authority
}

// Helper methods for JSON serialization (temporary until proto implementation)

// GetEmissionParams retrieves emission parameters from storage
func (k Keeper) GetEmissionParams(ctx context.Context) (types.EmissionParams, error) {
	jsonStr, err := k.EmissionParamsRaw.Get(ctx)
	if err != nil {
		return types.DefaultEmissionParams(), err
	}

	var params types.EmissionParams
	if err := json.Unmarshal([]byte(jsonStr), &params); err != nil {
		return types.DefaultEmissionParams(), err
	}

	return params, nil
}

// SetEmissionParams stores emission parameters to storage
func (k Keeper) SetEmissionParams(ctx context.Context, params types.EmissionParams) error {
	jsonBytes, err := json.Marshal(params)
	if err != nil {
		return err
	}

	return k.EmissionParamsRaw.Set(ctx, string(jsonBytes))
}

// SetEmissionHistory stores emission parameters history
func (k Keeper) SetEmissionHistory(ctx context.Context, height int64, params types.EmissionParams) error {
	jsonBytes, err := json.Marshal(params)
	if err != nil {
		return err
	}

	return k.EmissionHistoryRaw.Set(ctx, height, string(jsonBytes))
}

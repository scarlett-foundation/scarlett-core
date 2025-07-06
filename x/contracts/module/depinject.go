package contracts

import (
	"cosmossdk.io/core/appmodule"
	"cosmossdk.io/core/store"
	"cosmossdk.io/depinject"
	"cosmossdk.io/depinject/appconfig"
	"cosmossdk.io/log"

	"github.com/cosmos/cosmos-sdk/codec"
	authkeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"

	"scarlett-core/x/contracts/keeper"
	"scarlett-core/x/contracts/types"
)

var _ depinject.OnePerModuleType = AppModule{}

// IsOnePerModuleType implements the depinject.OnePerModuleType interface.
func (am AppModule) IsOnePerModuleType() {}

func init() {
	appconfig.RegisterModule(&types.Module{},
		appconfig.Provide(ProvideModule),
	)
}

type ModuleInputs struct {
	depinject.In

	Cdc          codec.Codec
	StoreService store.KVStoreService
	Logger       log.Logger

	// Use the exact same pattern as the emissions module
	// Some keepers are value types, some are pointer types
	AccountKeeper authkeeper.AccountKeeper // Value type like emissions
	BankKeeper    bankkeeper.Keeper        // Value type like emissions
	StakingKeeper *stakingkeeper.Keeper    // Pointer type like emissions
	// DistributionKeeper and ChannelKeeper are set to nil in NewKeeper for now
	// since they don't match the exact interfaces wasmd expects
}

type ModuleOutputs struct {
	depinject.Out

	ContractsKeeper keeper.Keeper
	Module          appmodule.AppModule
}

func ProvideModule(in ModuleInputs) ModuleOutputs {
	// Create the contracts keeper which will handle the manual wasm.Keeper instantiation
	keeper := keeper.NewKeeper(
		in.Cdc,
		in.StoreService,
		in.Logger,
		in.AccountKeeper,
		in.BankKeeper,
		*in.StakingKeeper,
	)

	m := NewAppModule(in.Cdc, keeper)

	return ModuleOutputs{ContractsKeeper: keeper, Module: m}
}

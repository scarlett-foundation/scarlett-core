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
	contractstypes "scarlett-core/x/contracts/types"
)

var _ depinject.OnePerModuleType = AppModule{}

// IsOnePerModuleType implements the depinject.OnePerModuleType interface.
func (am AppModule) IsOnePerModuleType() {}

func init() {
	appconfig.RegisterModule(
		&contractstypes.Module{},
		appconfig.Provide(ProvideModule),
	)
}

type ModuleInputs struct {
	depinject.In

	Cdc           codec.Codec
	StoreService  store.KVStoreService
	Logger        log.Logger
	AccountKeeper authkeeper.AccountKeeper
	BankKeeper    bankkeeper.Keeper
	StakingKeeper *stakingkeeper.Keeper

	Config *contractstypes.Module
}

type ModuleOutputs struct {
	depinject.Out

	ContractsKeeper keeper.Keeper
	Module          appmodule.AppModule
}

func ProvideModule(in ModuleInputs) ModuleOutputs {
	// Use a default home directory approach that works with depinject
	// The actual home directory will be set by the runtime when the app starts
	// For now, we use a placeholder that will be replaced during keeper initialization
	defaultHomeDir := ""

	k := keeper.NewKeeper(
		in.Cdc,
		in.StoreService,
		in.Logger,
		in.AccountKeeper,
		in.BankKeeper,
		*in.StakingKeeper,
		defaultHomeDir,
	)

	m := NewAppModule(in.Cdc, k)

	return ModuleOutputs{
		ContractsKeeper: k,
		Module:          m,
	}
}

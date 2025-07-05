package contracts

import (
	"cosmossdk.io/core/appmodule"
	"cosmossdk.io/core/store"
	"cosmossdk.io/depinject"
	"cosmossdk.io/depinject/appconfig"
	"cosmossdk.io/log"
	wasmkeeper "github.com/CosmWasm/wasmd/x/wasm/keeper"
	"github.com/cosmos/cosmos-sdk/codec"

	"scarlett-core/x/contracts/keeper"
	"scarlett-core/x/contracts/types"
)

var _ depinject.OnePerModuleType = AppModule{}

// IsOnePerModuleType implements the depinject.OnePerModuleType interface.
func (AppModule) IsOnePerModuleType() {}

func init() {
	appconfig.Register(
		&types.Module{},
		appconfig.Provide(ProvideModule),
	)
}

type ModuleInputs struct {
	depinject.In

	Config       *types.Module
	StoreService store.KVStoreService
	Cdc          codec.Codec
	Logger       log.Logger

	// We need the wasmd keeper to wrap it
	WasmKeeper *wasmkeeper.Keeper

	AuthKeeper    types.AuthKeeper
	BankKeeper    types.BankKeeper
	AccountKeeper types.AccountKeeper
}

type ModuleOutputs struct {
	depinject.Out

	ContractsKeeper keeper.Keeper
	Module          appmodule.AppModule
}

func ProvideModule(in ModuleInputs) ModuleOutputs {
	k := keeper.NewKeeper(
		in.Cdc,
		in.StoreService,
		in.Logger,
		in.WasmKeeper,
		in.AuthKeeper,
		in.BankKeeper,
		in.AccountKeeper,
	)
	m := NewAppModule(in.Cdc, k, in.AuthKeeper, in.BankKeeper)

	return ModuleOutputs{ContractsKeeper: k, Module: m}
}

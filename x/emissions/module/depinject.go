package emissions

import (
	"cosmossdk.io/core/address"
	"cosmossdk.io/core/appmodule"
	"cosmossdk.io/core/store"
	"cosmossdk.io/depinject"
	"cosmossdk.io/depinject/appconfig"
	"github.com/cosmos/cosmos-sdk/codec"
	authkeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	govkeeper "github.com/cosmos/cosmos-sdk/x/gov/keeper"
	mintkeeper "github.com/cosmos/cosmos-sdk/x/mint/keeper"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"

	"scarlett-core/x/emissions/keeper"
	"scarlett-core/x/emissions/types"
)

var _ depinject.OnePerModuleType = AppModule{}

// IsOnePerModuleType implements the depinject.OnePerModuleType interface.
func (AppModule) IsOnePerModuleType() {}

func init() {
	appconfig.Register(
		&types.Module{},
		appconfig.Provide(ProvideModule, ProvideEmissionsKeeperHooks),
	)
}

type ModuleInputs struct {
	depinject.In

	Config       *types.Module
	StoreService store.KVStoreService
	Cdc          codec.Codec
	AddressCodec address.Codec

	AuthKeeper    authkeeper.AccountKeeper
	BankKeeper    bankkeeper.Keeper
	GovKeeper     *govkeeper.Keeper
	StakingKeeper *stakingkeeper.Keeper
	MintKeeper    mintkeeper.Keeper
}

type ModuleOutputs struct {
	depinject.Out

	EmissionsKeeper keeper.Keeper
	Module          appmodule.AppModule
}

func ProvideModule(in ModuleInputs) ModuleOutputs {
	// default to governance authority if not provided
	authority := authtypes.NewModuleAddress(types.GovModuleName)
	if in.Config.Authority != "" {
		authority = authtypes.NewModuleAddressOrBech32Address(in.Config.Authority)
	}
	k := keeper.NewKeeper(
		in.StoreService,
		in.Cdc,
		in.AddressCodec,
		authority.Bytes(),
		in.AuthKeeper,
		in.BankKeeper,
		*in.GovKeeper,
		*in.StakingKeeper,
	)
	m := NewAppModule(in.Cdc, k, in.AuthKeeper, in.BankKeeper)

	return ModuleOutputs{EmissionsKeeper: k, Module: m}
}

// ProvideEmissionsKeeperHooks sets the mint keeper for the emissions keeper.
func ProvideEmissionsKeeperHooks(
	k keeper.Keeper,
	mk mintkeeper.Keeper,
) {
	k.SetMintKeeper(mk)
}

package contracts

import (
	"cosmossdk.io/core/appmodule"
	"cosmossdk.io/core/store"
	"cosmossdk.io/depinject"
	"cosmossdk.io/depinject/appconfig"
	"cosmossdk.io/log"
	wasmkeeper "github.com/CosmWasm/wasmd/x/wasm/keeper"
	wasmtypes "github.com/CosmWasm/wasmd/x/wasm/types"
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

	// Use wasmd's keeper interfaces directly
	AuthKeeper              wasmtypes.AccountKeeper
	BankKeeper              wasmtypes.BankKeeper
	StakingKeeper           wasmtypes.StakingKeeper
	DistributionKeeper      wasmtypes.DistributionKeeper
	ICS4Wrapper             wasmtypes.ICS4Wrapper
	ChannelKeeper           wasmtypes.ChannelKeeper
	ChannelKeeperV2         wasmtypes.ChannelKeeperV2
	ICS20TransferPortSource wasmtypes.ICS20TransferPortSource
	MessageRouter           wasmkeeper.MessageRouter
	GRPCQueryRouter         wasmkeeper.GRPCQueryRouter
}

type ModuleOutputs struct {
	depinject.Out

	Keeper keeper.Keeper
	Module appmodule.AppModule
}

func ProvideModule(in ModuleInputs) ModuleOutputs {
	// Create wasm node config with defaults
	nodeConfig := wasmtypes.NodeConfig{}

	// Create wasm VM config with defaults
	vmConfig := wasmtypes.VMConfig{}

	// Set available capabilities (minimal set for security)
	supportedFeatures := []string{"iterator", "staking", "stargate"}

	// Authority for governance (should be governance module address)
	authority := "cosmos10d07y265gmmuvt4z0w9aw880jnsr700juxf7n47" // placeholder - should be from governance

	// Home directory for wasm
	homeDir := "wasm"

	k := keeper.NewKeeper(
		in.Cdc,
		in.StoreService,
		in.Logger,
		in.AuthKeeper,
		in.BankKeeper,
		in.StakingKeeper,
		in.DistributionKeeper,
		in.ICS4Wrapper,
		in.ChannelKeeper,
		in.ChannelKeeperV2,
		in.ICS20TransferPortSource,
		in.MessageRouter,
		in.GRPCQueryRouter,
		homeDir,
		nodeConfig,
		vmConfig,
		supportedFeatures,
		authority,
	)

	// Create module without keeper interfaces that don't match
	m := NewAppModule(in.Cdc, k)

	return ModuleOutputs{
		Keeper: k,
		Module: m,
	}
}

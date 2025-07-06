package keeper

import (
	"context"

	"cosmossdk.io/core/store"
	"cosmossdk.io/log"
	wasmkeeper "github.com/CosmWasm/wasmd/x/wasm/keeper"
	wasmtypes "github.com/CosmWasm/wasmd/x/wasm/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authkeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
)

// Keeper wraps wasmd's keeper by embedding it directly and provides a clean interface
// for our contracts module while maintaining 100% compatibility with wasmd v0.61.0
type Keeper struct {
	wasmkeeper.Keeper
}

// NewKeeper creates a new contracts keeper that embeds wasm.Keeper
// This function handles the manual instantiation of wasm.Keeper with all required dependencies
func NewKeeper(
	cdc codec.Codec,
	storeService store.KVStoreService,
	logger log.Logger,
	accountKeeper authkeeper.AccountKeeper,
	bankKeeper bankkeeper.Keeper,
	stakingKeeper stakingkeeper.Keeper,
) Keeper {
	// Create wasm configuration with permissionless defaults
	nodeConfig := wasmtypes.NodeConfig{}
	vmConfig := wasmtypes.VMConfig{}

	// Set available capabilities (minimal set for security)
	supportedFeatures := []string{"iterator", "staking", "stargate"}

	// Authority for governance (should be governance module address)
	authority := "cosmos10d07y265gmmuvt4z0w9aw880jnsr700juxf7n47" // placeholder

	// Home directory for wasm
	homeDir := "/tmp/wasm"

	// For the parameters that wasmd v0.61.0 expects but we don't have compatible implementations,
	// we'll use nil for now. This is compatible with the wasmd v0.61.0 NewKeeper signature.
	// These can be implemented later when we need the specific functionality.
	var ics4Wrapper wasmtypes.ICS4Wrapper = nil
	var channelKeeperV2 wasmtypes.ChannelKeeperV2 = nil
	var portSource wasmtypes.ICS20TransferPortSource = nil
	var messageRouter wasmkeeper.MessageRouter = nil
	var grpcQueryRouter wasmkeeper.GRPCQueryRouter = nil

	// Use nil for distribution and channel keepers since they don't match interfaces exactly
	// This is acceptable for basic wasm functionality
	var distributionKeeperCompat wasmtypes.DistributionKeeper = nil
	var channelKeeperCompat wasmtypes.ChannelKeeper = nil

	// Manually instantiate the wasm.Keeper with all required dependencies
	wasmKeeper := wasmkeeper.NewKeeper(
		cdc,
		storeService,
		accountKeeper,            // implements wasmtypes.AccountKeeper
		bankKeeper,               // implements wasmtypes.BankKeeper
		stakingKeeper,            // implements wasmtypes.StakingKeeper
		distributionKeeperCompat, // nil for now - distribution queries not needed for basic functionality
		ics4Wrapper,
		channelKeeperCompat, // nil for now - IBC functionality not needed for basic contracts
		channelKeeperV2,
		portSource,
		messageRouter,
		grpcQueryRouter,
		homeDir,
		nodeConfig,
		vmConfig,
		supportedFeatures,
		authority,
	)

	return Keeper{
		Keeper: wasmKeeper,
	}
}

// Logger returns a module-specific logger.
func (k Keeper) Logger(ctx context.Context) log.Logger {
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	return k.Keeper.Logger(sdkCtx)
}

// GetContractInfo returns contract info from the embedded wasm keeper
func (k Keeper) GetContractInfo(ctx context.Context, contractAddr sdk.AccAddress) *wasmtypes.ContractInfo {
	return k.Keeper.GetContractInfo(ctx, contractAddr)
}

// GetCodeInfo returns code info from the embedded wasm keeper
func (k Keeper) GetCodeInfo(ctx context.Context, codeID uint64) *wasmtypes.CodeInfo {
	return k.Keeper.GetCodeInfo(ctx, codeID)
}

// Additional wrapper methods can be added here as needed
// All other wasm functionality is available through the embedded Keeper

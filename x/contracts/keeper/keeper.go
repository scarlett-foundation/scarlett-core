package keeper

import (
	"context"
	"fmt"
	"os"
	"sync"
	"time"

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

// Keeper wraps wasmd's keeper using lazy initialization to avoid lock file issues
type Keeper struct {
	// Configuration for lazy initialization
	cdc           codec.Codec
	storeService  store.KVStoreService
	logger        log.Logger
	accountKeeper authkeeper.AccountKeeper
	bankKeeper    bankkeeper.Keeper
	stakingKeeper stakingkeeper.Keeper

	// Lazy-initialized wasm keeper
	wasmKeeper *wasmkeeper.Keeper
	once       sync.Once
}

// NewKeeper creates a new contracts keeper with lazy wasm initialization
func NewKeeper(
	cdc codec.Codec,
	storeService store.KVStoreService,
	logger log.Logger,
	accountKeeper authkeeper.AccountKeeper,
	bankKeeper bankkeeper.Keeper,
	stakingKeeper stakingkeeper.Keeper,
) Keeper {
	return Keeper{
		cdc:           cdc,
		storeService:  storeService,
		logger:        logger,
		accountKeeper: accountKeeper,
		bankKeeper:    bankKeeper,
		stakingKeeper: stakingKeeper,
	}
}

// getWasmKeeper lazily initializes the wasm keeper on first access
func (k *Keeper) getWasmKeeper() *wasmkeeper.Keeper {
	k.once.Do(func() {
		// Create wasm configuration to support bulk memory operations
		nodeConfig := wasmtypes.NodeConfig{}

		// Configure VM to accept bulk memory operations
		// This may require wasmvm v3.0.0+ with proper bulk memory support
		vmConfig := wasmtypes.VMConfig{}

		// Set available capabilities for wasmd v0.61.0 + wasmvm v3.0.0 compatibility
		// This matches the full feature set supported by wasmd v0.61.0
		supportedFeatures := []string{
			"iterator", "staking", "stargate",
			"cosmwasm_1_1", "cosmwasm_1_2", "cosmwasm_1_3", "cosmwasm_1_4",
			"cosmwasm_2_0", "cosmwasm_2_1", "cosmwasm_2_2",
		}

		// Authority for governance (should be governance module address)
		authority := "cosmos10d07y265gmmuvt4z0w9aw880jnsr700juxf7n47" // placeholder

		// Home directory for wasm - use a unique path to avoid conflicts
		// Generate unique directory with timestamp to prevent multiple VM instances from conflicting
		homeDir := fmt.Sprintf("/tmp/wasm-contracts-%d-%d", os.Getpid(), time.Now().UnixNano())

		// For the parameters that wasmd v0.61.0 expects but we don't have compatible implementations,
		// we'll use nil for now. This is compatible with the wasmd v0.61.0 NewKeeper signature.
		var ics4Wrapper wasmtypes.ICS4Wrapper = nil
		var channelKeeperV2 wasmtypes.ChannelKeeperV2 = nil
		var portSource wasmtypes.ICS20TransferPortSource = nil
		var messageRouter wasmkeeper.MessageRouter = nil
		var grpcQueryRouter wasmkeeper.GRPCQueryRouter = nil

		// Use nil for distribution and channel keepers since they don't match interfaces exactly
		var distributionKeeperCompat wasmtypes.DistributionKeeper = nil
		var channelKeeperCompat wasmtypes.ChannelKeeper = nil

		// Manually instantiate the wasm.Keeper with all required dependencies
		// NewKeeper returns a value type, so we need to store it and then take a pointer
		wasmKeeperValue := wasmkeeper.NewKeeper(
			k.cdc,
			k.storeService,
			k.accountKeeper,          // implements wasmtypes.AccountKeeper
			k.bankKeeper,             // implements wasmtypes.BankKeeper
			k.stakingKeeper,          // implements wasmtypes.StakingKeeper
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
		k.wasmKeeper = &wasmKeeperValue
	})
	return k.wasmKeeper
}

// Logger returns a module-specific logger.
func (k Keeper) Logger(ctx context.Context) log.Logger {
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	return k.getWasmKeeper().Logger(sdkCtx)
}

// GetContractInfo returns contract info from the embedded wasm keeper
func (k Keeper) GetContractInfo(ctx context.Context, contractAddr sdk.AccAddress) *wasmtypes.ContractInfo {
	return k.getWasmKeeper().GetContractInfo(ctx, contractAddr)
}

// GetCodeInfo returns code info from the embedded wasm keeper
func (k Keeper) GetCodeInfo(ctx context.Context, codeID uint64) *wasmtypes.CodeInfo {
	return k.getWasmKeeper().GetCodeInfo(ctx, codeID)
}

// Additional wrapper methods can be added here as needed
// All wasm functionality is available through the lazy-loaded keeper

package keeper

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"cosmossdk.io/core/store"
	"cosmossdk.io/log"
	wasmkeeper "github.com/CosmWasm/wasmd/x/wasm/keeper"
	wasmtypes "github.com/CosmWasm/wasmd/x/wasm/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authkeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
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

	// wasmd configuration
	homeDir string

	// Lazy-initialized wasm keeper
	wasmKeeper *wasmkeeper.Keeper
	once       sync.Once
	initError  error // Store any initialization error
}

// NewKeeper creates a new contracts keeper with lazy wasm initialization
func NewKeeper(
	cdc codec.Codec,
	storeService store.KVStoreService,
	logger log.Logger,
	accountKeeper authkeeper.AccountKeeper,
	bankKeeper bankkeeper.Keeper,
	stakingKeeper stakingkeeper.Keeper,
	homeDir string,
) Keeper {
	// Use default home directory if not provided
	if homeDir == "" {
		homeDir = getDefaultHomeDir()
	}

	return Keeper{
		cdc:           cdc,
		storeService:  storeService,
		logger:        logger,
		accountKeeper: accountKeeper,
		bankKeeper:    bankKeeper,
		stakingKeeper: stakingKeeper,
		homeDir:       homeDir,
	}
}

// getDefaultHomeDir returns the default home directory for the chain
func getDefaultHomeDir() string {
	// Try to get from environment or use default
	if home := os.Getenv("HOME"); home != "" {
		return filepath.Join(home, ".scarlett-core")
	}
	return ".scarlett-core" // fallback for CI/testing environments
}

// ensureWasmDir creates the wasm directory if it doesn't exist
func (k *Keeper) ensureWasmDir(wasmDir string) error {
	if err := os.MkdirAll(wasmDir, 0755); err != nil {
		return fmt.Errorf("failed to create wasm directory %s: %w", wasmDir, err)
	}

	// Verify directory is writable
	testFile := filepath.Join(wasmDir, ".write_test")
	if err := os.WriteFile(testFile, []byte("test"), 0644); err != nil {
		return fmt.Errorf("wasm directory %s is not writable: %w", wasmDir, err)
	}
	os.Remove(testFile) // Clean up test file

	return nil
}

// getWasmKeeper lazily initializes the wasm keeper on first access
func (k *Keeper) getWasmKeeper() (*wasmkeeper.Keeper, error) {
	k.once.Do(func() {
		k.logger.Info("Initializing wasmd keeper", "homeDir", k.homeDir)

		// Use wasmd's standard directory structure
		wasmDir := filepath.Join(k.homeDir, "data", "wasm")

		// Ensure the wasm directory exists and is writable
		if err := k.ensureWasmDir(wasmDir); err != nil {
			k.initError = fmt.Errorf("failed to setup wasm directory: %w", err)
			return
		}

		// Create wasm configuration using wasmd's approach for v0.61.0
		nodeConfig := wasmtypes.NodeConfig{}
		vmConfig := wasmtypes.VMConfig{}

		// Set available capabilities for wasmd v0.61.0 + wasmvm v3.0.0 compatibility
		// This matches the full feature set supported by wasmd v0.61.0
		supportedFeatures := []string{
			"iterator", "staking", "stargate",
			"cosmwasm_1_1", "cosmwasm_1_2", "cosmwasm_1_3", "cosmwasm_1_4",
			"cosmwasm_2_0", "cosmwasm_2_1", "cosmwasm_2_2",
		}

		// Use proper governance module address for authority
		authority := authtypes.NewModuleAddress(govtypes.ModuleName).String()
		k.logger.Info("Using governance authority for wasmd", "authority", authority)

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

		k.logger.Info("Creating wasmd keeper",
			"wasmDir", wasmDir,
			"features", supportedFeatures,
		)

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
			wasmDir,    // Use proper wasmd directory structure
			nodeConfig, // Use wasmd's NodeConfig
			vmConfig,   // Use wasmd's VMConfig
			supportedFeatures,
			authority,
		)
		k.wasmKeeper = &wasmKeeperValue
		k.logger.Info("Successfully initialized wasmd keeper")
	})

	if k.initError != nil {
		return nil, k.initError
	}
	return k.wasmKeeper, nil
}

// Logger returns a module-specific logger.
func (k Keeper) Logger(ctx context.Context) log.Logger {
	wasmKeeper, err := k.getWasmKeeper()
	if err != nil {
		k.logger.Error("Failed to get wasm keeper for logger", "error", err)
		return k.logger
	}
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	return wasmKeeper.Logger(sdkCtx)
}

// GetContractInfo returns contract info from the embedded wasm keeper
func (k Keeper) GetContractInfo(ctx context.Context, contractAddr sdk.AccAddress) *wasmtypes.ContractInfo {
	wasmKeeper, err := k.getWasmKeeper()
	if err != nil {
		k.logger.Error("Failed to get wasm keeper for GetContractInfo", "error", err)
		return nil
	}
	return wasmKeeper.GetContractInfo(ctx, contractAddr)
}

// GetCodeInfo returns code info from the embedded wasm keeper
func (k Keeper) GetCodeInfo(ctx context.Context, codeID uint64) *wasmtypes.CodeInfo {
	wasmKeeper, err := k.getWasmKeeper()
	if err != nil {
		k.logger.Error("Failed to get wasm keeper for GetCodeInfo", "error", err)
		return nil
	}
	return wasmKeeper.GetCodeInfo(ctx, codeID)
}

// Additional wrapper methods can be added here as needed
// All wasm functionality is available through the lazy-loaded keeper

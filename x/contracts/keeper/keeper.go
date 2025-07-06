package keeper

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

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

// Keeper wraps wasmd's keeper with proper initialization
type Keeper struct {
	// Configuration
	cdc           codec.Codec
	storeService  store.KVStoreService
	logger        log.Logger
	accountKeeper authkeeper.AccountKeeper
	bankKeeper    bankkeeper.Keeper
	stakingKeeper stakingkeeper.Keeper

	// wasmd configuration
	homeDir string

	// Initialized wasm keeper
	wasmKeeper *wasmkeeper.Keeper
}

// NewKeeper creates a new contracts keeper with immediate wasm initialization
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

	k := Keeper{
		cdc:           cdc,
		storeService:  storeService,
		logger:        logger,
		accountKeeper: accountKeeper,
		bankKeeper:    bankKeeper,
		stakingKeeper: stakingKeeper,
		homeDir:       homeDir,
	}

	// Initialize wasmd keeper immediately during app startup
	wasmKeeper, err := k.initializeWasmKeeper()
	if err != nil {
		// Log error but don't panic - let the app start and handle errors gracefully
		logger.Error("Failed to initialize wasmd keeper", "error", err)
		// Return keeper with nil wasmKeeper - methods will handle this gracefully
		return k
	}

	k.wasmKeeper = wasmKeeper
	return k
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

// cleanupStaleLockFiles removes stale exclusive.lock files that might prevent VM startup
func (k *Keeper) cleanupStaleLockFiles(wasmDir string) {
	lockFile := filepath.Join(wasmDir, "wasm", "exclusive.lock")

	// Check if lock file exists
	if _, err := os.Stat(lockFile); os.IsNotExist(err) {
		return // No lock file to clean up
	}

	// Remove the lock file - it's likely stale from a previous run
	if err := os.Remove(lockFile); err != nil {
		k.logger.Warn("Failed to remove stale lock file", "lockFile", lockFile, "error", err)
	} else {
		k.logger.Info("Removed stale lock file", "lockFile", lockFile)
	}
}

// initializeWasmKeeper initializes the wasm keeper during app startup
func (k *Keeper) initializeWasmKeeper() (*wasmkeeper.Keeper, error) {
	k.logger.Info("Initializing wasmd keeper", "homeDir", k.homeDir)

	// Use wasmd's standard directory structure
	// wasmd will create a 'wasm' subdirectory, so we use 'data' as the base
	wasmDir := filepath.Join(k.homeDir, "data")

	// Ensure the wasm directory exists and is writable
	if err := k.ensureWasmDir(wasmDir); err != nil {
		return nil, fmt.Errorf("failed to setup wasm directory: %w", err)
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

	// Clean up any stale lock files immediately before VM initialization
	k.cleanupStaleLockFiles(wasmDir)

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
	k.logger.Info("Successfully initialized wasmd keeper")
	return &wasmKeeperValue, nil
}

// Logger returns a module-specific logger.
func (k Keeper) Logger(ctx context.Context) log.Logger {
	if k.wasmKeeper == nil {
		k.logger.Error("Wasm keeper not initialized")
		return k.logger
	}
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	return k.wasmKeeper.Logger(sdkCtx)
}

// GetContractInfo returns contract info from the embedded wasm keeper
func (k Keeper) GetContractInfo(ctx context.Context, contractAddr sdk.AccAddress) *wasmtypes.ContractInfo {
	if k.wasmKeeper == nil {
		k.logger.Error("Wasm keeper not initialized")
		return nil
	}
	return k.wasmKeeper.GetContractInfo(ctx, contractAddr)
}

// GetCodeInfo returns code info from the embedded wasm keeper
func (k Keeper) GetCodeInfo(ctx context.Context, codeID uint64) *wasmtypes.CodeInfo {
	if k.wasmKeeper == nil {
		k.logger.Error("Wasm keeper not initialized")
		return nil
	}
	return k.wasmKeeper.GetCodeInfo(ctx, codeID)
}

// Additional wrapper methods can be added here as needed
// All wasm functionality is available through the lazy-loaded keeper

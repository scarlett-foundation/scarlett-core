package keeper

import (
	"fmt"
	"path/filepath"

	"cosmossdk.io/collections"
	"cosmossdk.io/core/address"
	corestore "cosmossdk.io/core/store"
	"github.com/cosmos/cosmos-sdk/codec"

	wasmvm "github.com/CosmWasm/wasmvm/v3"

	"scarlett-core/x/contracts/types"

	"cosmossdk.io/log"
)

type Keeper struct {
	storeService corestore.KVStoreService
	cdc          codec.Codec
	addressCodec address.Codec
	// Address capable of executing a MsgUpdateParams message.
	// Typically, this should be the x/gov module account.
	authority []byte

	Schema collections.Schema
	Params collections.Item[types.Params]

	// WasmVM instance for contract execution (lazy initialization)
	wasmVM     *wasmvm.VM
	wasmDir    string
	wasmVMInit bool

	// Gas configuration for contract operations
	gasConfig types.GasConfig

	// Contract storage
	ContractCode collections.Map[[]byte, []byte]             // codeID -> wasm bytecode
	ContractInfo collections.Map[string, types.ContractInfo] // address -> contract info
	ContractSeq  collections.Sequence                        // for generating contract IDs

	bankKeeper      types.BankKeeper
	emissionsKeeper types.EmissionsKeeper
}

func NewKeeper(
	storeService corestore.KVStoreService,
	cdc codec.Codec,
	addressCodec address.Codec,
	authority []byte,
	wasmDir string,

	bankKeeper types.BankKeeper,
	emissionsKeeper types.EmissionsKeeper,
) Keeper {
	if _, err := addressCodec.BytesToString(authority); err != nil {
		panic(fmt.Sprintf("invalid authority address %s: %s", authority, err))
	}

	sb := collections.NewSchemaBuilder(storeService)

	k := Keeper{
		storeService: storeService,
		cdc:          cdc,
		addressCodec: addressCodec,
		authority:    authority,
		wasmDir:      wasmDir,
		wasmVMInit:   false,                // Lazy initialization
		gasConfig:    types.GetGasConfig(), // Load embedded gas configuration

		bankKeeper:      bankKeeper,
		emissionsKeeper: emissionsKeeper,

		Params:       collections.NewItem(sb, types.ParamsKey, "params", codec.CollValue[types.Params](cdc)),
		ContractCode: collections.NewMap(sb, types.ContractCodeKey, "contract_code", collections.BytesKey, collections.BytesValue),
		ContractInfo: collections.NewMap(sb, types.ContractInfoKey, "contract_info", collections.StringKey, codec.CollValue[types.ContractInfo](cdc)),
		ContractSeq:  collections.NewSequence(sb, types.ContractSeqKey, "contract_seq"),
	}

	schema, err := sb.Build()
	if err != nil {
		panic(err)
	}
	k.Schema = schema

	return k
}

// GetAuthority returns the module's authority.
func (k Keeper) GetAuthority() []byte {
	return k.authority
}

// GetGasConfig returns the current gas configuration
func (k Keeper) GetGasConfig() types.GasConfig {
	return k.gasConfig
}

// initWasmVM initializes WasmVM lazily when first needed
func (k *Keeper) initWasmVM(logger log.Logger) error {
	if k.wasmVMInit {
		return nil
	}

	// Initialize WasmVM with production-ready configuration
	// NewVM(dataDir string, supportedCapabilities []string, memoryLimit uint32, printDebug bool, cacheSize uint32) (*VM, error)
	wasmCacheDir := filepath.Join(k.wasmDir, "wasm", "cache")
	supportedCapabilities := []string{"iterator", "staking", "stargate", "cosmwasm_1_1", "cosmwasm_1_2"} // Conservative capabilities to avoid v3.0.0 bugs
	memoryLimit := uint32(256)                                                                           // 256 MiB memory limit - sufficient for large contracts
	printDebug := true                                                                                   // Enable debug to see WasmVM internals
	cacheSize := uint32(200)                                                                             // 200 MiB cache - increased for better performance

	// Log WasmVM initialization parameters using proper SDK logging
	logger.Info("🔧 Initializing WasmVM with configuration",
		"cache_directory", wasmCacheDir,
		"memory_limit_mib", memoryLimit,
		"cache_size_mib", cacheSize,
		"capabilities", supportedCapabilities,
		"debug_mode", printDebug)

	vm, err := wasmvm.NewVM(wasmCacheDir, supportedCapabilities, memoryLimit, printDebug, cacheSize)
	if err != nil {
		logger.Error("❌ Failed to initialize WasmVM",
			"error", err,
			"attempted_memory_limit_mib", memoryLimit,
			"attempted_cache_size_mib", cacheSize,
			"cache_directory", wasmCacheDir)
		return fmt.Errorf("failed to initialize WasmVM: %w", err)
	}

	logger.Info("✅ WasmVM initialized successfully",
		"memory_limit_mib", memoryLimit,
		"cache_size_mib", cacheSize)

	k.wasmVM = vm
	k.wasmVMInit = true
	return nil
}

// GetWasmVM returns the WasmVM instance, initializing it if needed
func (k *Keeper) GetWasmVM(logger log.Logger) (*wasmvm.VM, error) {
	if err := k.initWasmVM(logger); err != nil {
		return nil, err
	}
	return k.wasmVM, nil
}

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
		wasmVMInit:   false, // Lazy initialization

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

// initWasmVM initializes WasmVM lazily when first needed
func (k *Keeper) initWasmVM() error {
	if k.wasmVMInit {
		return nil
	}

	// Initialize WasmVM with secure configuration
	// NewVM(dataDir string, supportedCapabilities []string, memoryLimit uint32, printDebug bool, cacheSize uint32) (*VM, error)
	wasmCacheDir := filepath.Join(k.wasmDir, "wasm", "cache")
	supportedCapabilities := []string{"iterator", "staking", "stargate", "cosmwasm_1_1", "cosmwasm_1_2", "cosmwasm_2_0", "bulk-memory"} // Standard capabilities with bulk-memory (supported in wasmvm v3)
	memoryLimit := uint32(32)                                                                                                           // 32 MiB memory limit per contract
	printDebug := false                                                                                                                 // Don't print debug in production
	cacheSize := uint32(256)                                                                                                            // 256 MiB cache

	vm, err := wasmvm.NewVM(wasmCacheDir, supportedCapabilities, memoryLimit, printDebug, cacheSize)
	if err != nil {
		return fmt.Errorf("failed to initialize WasmVM: %w", err)
	}

	k.wasmVM = vm
	k.wasmVMInit = true
	return nil
}

// GetWasmVM returns the WasmVM instance, initializing it if needed
func (k *Keeper) GetWasmVM() (*wasmvm.VM, error) {
	if err := k.initWasmVM(); err != nil {
		return nil, err
	}
	return k.wasmVM, nil
}

package keeper

import (
	"context"

	"cosmossdk.io/core/address"
	"cosmossdk.io/core/store"
	"cosmossdk.io/log"
	wasmkeeper "github.com/CosmWasm/wasmd/x/wasm/keeper"
	wasmtypes "github.com/CosmWasm/wasmd/x/wasm/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"scarlett-core/x/contracts/types"
)

// Keeper wraps wasmd's keeper to provide contracts module functionality
type Keeper struct {
	cdc      codec.BinaryCodec
	storeKey store.KVStoreService
	logger   log.Logger

	// Embed wasmd's keeper to delegate all wasm operations
	*wasmkeeper.Keeper

	// Expected keepers for additional functionality
	authKeeper    types.AuthKeeper
	bankKeeper    types.BankKeeper
	accountKeeper types.AccountKeeper
}

// NewKeeper creates a new contracts keeper that wraps wasmd's functionality
func NewKeeper(
	cdc codec.BinaryCodec,
	storeKey store.KVStoreService,
	logger log.Logger,
	wasmKeeper *wasmkeeper.Keeper,
	authKeeper types.AuthKeeper,
	bankKeeper types.BankKeeper,
	accountKeeper types.AccountKeeper,
) Keeper {
	return Keeper{
		cdc:           cdc,
		storeKey:      storeKey,
		logger:        logger,
		Keeper:        wasmKeeper, // Embed wasmd's keeper
		authKeeper:    authKeeper,
		bankKeeper:    bankKeeper,
		accountKeeper: accountKeeper,
	}
}

// Logger returns the module logger
func (k Keeper) Logger() log.Logger {
	return k.logger.With("module", "x/"+types.ModuleName)
}

// GetContractInfo delegates to wasmd's keeper
func (k Keeper) GetContractInfo(ctx context.Context, contractAddress sdk.AccAddress) *wasmtypes.ContractInfo {
	return k.Keeper.GetContractInfo(ctx, contractAddress)
}

// GetCodeInfo delegates to wasmd's keeper
func (k Keeper) GetCodeInfo(ctx context.Context, codeID uint64) *wasmtypes.CodeInfo {
	return k.Keeper.GetCodeInfo(ctx, codeID)
}

// GetKeeper returns the embedded wasmd keeper for compatibility
func (k Keeper) GetKeeper() *wasmkeeper.Keeper {
	return k.Keeper
}

// AddressCodec returns the address codec from auth keeper
func (k Keeper) AddressCodec() address.Codec {
	return k.authKeeper.AddressCodec()
}

// ValidateContractAddress validates that an address is a valid contract address
func (k Keeper) ValidateContractAddress(ctx context.Context, addr sdk.AccAddress) error {
	contractInfo := k.GetContractInfo(ctx, addr)
	if contractInfo == nil {
		return types.ErrInvalidSigner.Wrapf("address %s is not a contract", addr.String())
	}
	return nil
}

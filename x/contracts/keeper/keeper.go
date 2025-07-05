package keeper

import (
	"context"

	"cosmossdk.io/core/store"
	"cosmossdk.io/log"
	wasmkeeper "github.com/CosmWasm/wasmd/x/wasm/keeper"
	wasmtypes "github.com/CosmWasm/wasmd/x/wasm/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Keeper wraps wasmd's keeper by embedding it directly and providing contracts module functionality
type Keeper struct {
	cdc      codec.BinaryCodec
	storeKey store.KVStoreService
	logger   log.Logger

	// Embed wasmd's keeper directly (not a pointer)
	wasmkeeper.Keeper
}

// NewKeeper creates a new contracts keeper that manually instantiates wasm.Keeper
func NewKeeper(
	cdc codec.Codec,
	storeKey store.KVStoreService,
	logger log.Logger,
	authKeeper wasmtypes.AccountKeeper,
	bankKeeper wasmtypes.BankKeeper,
	stakingKeeper wasmtypes.StakingKeeper,
	distributionKeeper wasmtypes.DistributionKeeper,
	ics4Wrapper wasmtypes.ICS4Wrapper,
	channelKeeper wasmtypes.ChannelKeeper,
	channelKeeperV2 wasmtypes.ChannelKeeperV2,
	ics20TransferPortSource wasmtypes.ICS20TransferPortSource,
	messageRouter wasmkeeper.MessageRouter,
	grpcQueryRouter wasmkeeper.GRPCQueryRouter,
	homeDir string,
	nodeConfig wasmtypes.NodeConfig,
	vmConfig wasmtypes.VMConfig,
	supportedFeatures []string,
	authority string,
	opts ...wasmkeeper.Option,
) Keeper {
	// Manually instantiate wasmd's keeper
	wasmKeeper := wasmkeeper.NewKeeper(
		cdc,
		storeKey,
		authKeeper,
		bankKeeper,
		stakingKeeper,
		distributionKeeper,
		ics4Wrapper,
		channelKeeper,
		channelKeeperV2,
		ics20TransferPortSource,
		messageRouter,
		grpcQueryRouter,
		homeDir,
		nodeConfig,
		vmConfig,
		supportedFeatures,
		authority,
		opts...,
	)

	return Keeper{
		cdc:      cdc,
		storeKey: storeKey,
		logger:   logger,
		Keeper:   wasmKeeper,
	}
}

// Logger returns a module-specific logger.
func (k Keeper) Logger() log.Logger {
	return k.logger.With("module", "x/contracts")
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

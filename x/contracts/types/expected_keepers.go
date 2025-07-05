package types

import (
	"context"

	"cosmossdk.io/core/address"
	wasmtypes "github.com/CosmWasm/wasmd/x/wasm/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth/types"
)

// AuthKeeper defines the expected interface for the Auth module.
type AuthKeeper interface {
	AddressCodec() address.Codec
	GetAccount(context.Context, sdk.AccAddress) sdk.AccountI // only used for simulation
	// Methods imported from account should be defined here
}

// BankKeeper defines the expected interface for the Bank module.
type BankKeeper interface {
	SpendableCoins(context.Context, sdk.AccAddress) sdk.Coins
	// Methods imported from bank should be defined here
}

// ParamSubspace defines the expected Subspace interface for parameters.
type ParamSubspace interface {
	Get(context.Context, []byte, interface{})
	Set(context.Context, []byte, interface{})
}

// AccountKeeper defines the expected interface for the Account module.
type AccountKeeper interface {
	GetAccount(context.Context, sdk.AccAddress) types.AccountI // only used for simulation
}

// WasmKeeper defines the expected interface for the Wasm module.
// This interface delegates to wasmd's actual keeper methods.
type WasmKeeper interface {
	// Contract information queries
	GetContractInfo(ctx context.Context, contractAddr sdk.AccAddress) *wasmtypes.ContractInfo
	GetCodeInfo(ctx context.Context, codeID uint64) *wasmtypes.CodeInfo

	// Basic wasm operations that we can delegate to
	HasContractInfo(ctx context.Context, contractAddr sdk.AccAddress) bool
}

// Note: We reference wasmd's keeper interfaces directly since they're already defined
// This ensures 100% compatibility with wasmd v0.61.0

// StakingKeeper is an alias for wasmd's staking keeper interface
type StakingKeeper = wasmtypes.StakingKeeper

// DistributionKeeper is an alias for wasmd's distribution keeper interface
type DistributionKeeper = wasmtypes.DistributionKeeper

// ChannelKeeper is an alias for wasmd's IBC channel keeper interface
type ChannelKeeper = wasmtypes.ChannelKeeper

// Note: PortKeeper, CapabilityKeeper, and ScopedKeeper are not available in wasmd v0.61.0
// They were removed or are not exposed in the public API

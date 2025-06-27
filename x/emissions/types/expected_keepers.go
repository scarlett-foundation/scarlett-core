package types

import (
	"context"

	"cosmossdk.io/core/address"
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types/v1beta1"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"
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
	SendCoinsFromModuleToModule(ctx context.Context, senderModule, recipientModule string, amt sdk.Coins) error
	GetBalance(ctx context.Context, addr sdk.AccAddress, denom string) sdk.Coin
	GetAllBalances(ctx context.Context, addr sdk.AccAddress) sdk.Coins
	// Methods imported from bank should be defined here
}

// GovKeeper defines the expected interface for the Gov module.
type GovKeeper interface {
	GetGovernanceAccount(ctx context.Context) sdk.AccAddress
	GetParams(ctx context.Context) (govtypes.Params, error)
}

// StakingKeeper defines the expected interface for the Staking module.
type StakingKeeper interface {
	BondedRatio(ctx context.Context) (math.LegacyDec, error)
	TotalBondedTokens(ctx context.Context) (math.Int, error)
}

// MintKeeper defines the expected interface for the Mint module.
type MintKeeper interface {
	GetMinter(ctx context.Context) (minttypes.Minter, error)
	GetParams(ctx context.Context) (minttypes.Params, error)
	MintCoins(ctx context.Context, newCoins sdk.Coins) error
}

// ParamSubspace defines the expected Subspace interface for parameters.
type ParamSubspace interface {
	Get(context.Context, []byte, interface{})
	Set(context.Context, []byte, interface{})
}

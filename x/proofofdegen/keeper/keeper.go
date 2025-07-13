package keeper

import (
	"fmt"

	"cosmossdk.io/collections"
	"cosmossdk.io/core/address"
	corestore "cosmossdk.io/core/store"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"scarlett-core/x/proofofdegen/types"
)

type Keeper struct {
	storeService corestore.KVStoreService
	cdc          codec.Codec
	addressCodec address.Codec
	// Address capable of executing a MsgUpdateParams message.
	// Typically, this should be the x/gov module account.
	authority []byte

	authKeeper types.AuthKeeper
	bankKeeper types.BankKeeper

	Schema         collections.Schema
	Params         collections.Item[types.Params]
	Campaign       collections.Map[string, types.Campaign]
	EligibleWallet collections.Map[string, types.EligibleWallet]
}

func NewKeeper(
	storeService corestore.KVStoreService,
	cdc codec.Codec,
	addressCodec address.Codec,
	authority []byte,
	authKeeper types.AuthKeeper,
	bankKeeper types.BankKeeper,

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
		authKeeper:   authKeeper,
		bankKeeper:   bankKeeper,

		Params:         collections.NewItem(sb, types.ParamsKey, "params", codec.CollValue[types.Params](cdc)),
		Campaign:       collections.NewMap(sb, types.CampaignKey, "campaign", collections.StringKey, codec.CollValue[types.Campaign](cdc)),
		EligibleWallet: collections.NewMap(sb, types.EligibleWalletKey, "eligibleWallet", collections.StringKey, codec.CollValue[types.EligibleWallet](cdc)),
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

// DistributeEmissions handles emissions received from the main emissions module
// In the accumulation model, emissions stay in the module account and accumulate
// until users claim their weighted shares
func (k Keeper) DistributeEmissions(ctx sdk.Context) error {
	// Get current module balance (accumulated emissions)
	moduleAccount := k.authKeeper.GetModuleAccount(ctx, types.ModuleName)
	if moduleAccount == nil {
		return fmt.Errorf("module account not found: %s", types.ModuleName)
	}

	moduleAddr := moduleAccount.GetAddress()
	moduleBalance := k.bankKeeper.GetBalance(ctx, moduleAddr, "sclt")
	currentBalance := moduleBalance.Amount.Uint64()

	// Emit event for transparency about emissions accumulation
	if currentBalance > 0 {
		ctx.EventManager().EmitEvent(
			sdk.NewEvent(
				"emissions_accumulated",
				sdk.NewAttribute("module", types.ModuleName),
				sdk.NewAttribute("total_balance", fmt.Sprintf("%d", currentBalance)),
				sdk.NewAttribute("block_height", fmt.Sprintf("%d", ctx.BlockHeight())),
			),
		)
	}

	// In accumulation model, we don't distribute tokens immediately
	// They stay in the module account until users claim their weighted shares
	// This creates the patience reward mechanics where later claimers get bigger shares

	return nil
}

// GetUnclaimedWallets returns all eligible wallets that haven't claimed yet
func (k Keeper) GetUnclaimedWallets(ctx sdk.Context) []types.EligibleWallet {
	var unclaimedWallets []types.EligibleWallet

	// Iterate through all eligible wallets
	err := k.EligibleWallet.Walk(ctx, nil, func(key string, wallet types.EligibleWallet) (bool, error) {
		if !wallet.Claimed {
			unclaimedWallets = append(unclaimedWallets, wallet)
		}
		return false, nil // Continue iteration
	})

	if err != nil {
		return nil
	}

	return unclaimedWallets
}

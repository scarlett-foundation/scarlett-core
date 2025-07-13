package keeper

import (
	"fmt"

	"cosmossdk.io/collections"
	"cosmossdk.io/core/address"
	corestore "cosmossdk.io/core/store"
	"cosmossdk.io/math"
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
		bankKeeper:   bankKeeper,

		Params:   collections.NewItem(sb, types.ParamsKey, "params", codec.CollValue[types.Params](cdc)),
		Campaign: collections.NewMap(sb, types.CampaignKey, "campaign", collections.StringKey, codec.CollValue[types.Campaign](cdc)), EligibleWallet: collections.NewMap(sb, types.EligibleWalletKey, "eligibleWallet", collections.StringKey, codec.CollValue[types.EligibleWallet](cdc))}

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

// DistributeEmissions distributes emissions to unclaimed wallets based on their weights
func (k Keeper) DistributeEmissions(ctx sdk.Context) error {
	// Get emissions received from main emissions module this block
	moduleAddr := k.bankKeeper.GetModuleAddress(types.ModuleName)
	moduleBalance := k.bankKeeper.GetBalance(ctx, moduleAddr, "sclt")
	receivedEmissions := moduleBalance.Amount.Uint64()

	if receivedEmissions == 0 {
		return nil // No emissions to distribute
	}

	// Get all unclaimed wallets from Genesis campaign
	unclaimedWallets := k.GetUnclaimedWallets(ctx)
	if len(unclaimedWallets) == 0 {
		return nil // No unclaimed wallets
	}

	// Calculate total weight of unclaimed wallets
	totalWeight := uint64(0)
	for _, wallet := range unclaimedWallets {
		totalWeight += wallet.Weight
	}

	if totalWeight == 0 {
		return nil // No weight to distribute
	}

	// Distribute emissions based on weights
	for _, wallet := range unclaimedWallets {
		// Calculate this wallet's share: (receivedEmissions * walletWeight) / totalWeight
		walletShare := (receivedEmissions * wallet.Weight) / totalWeight

		if walletShare > 0 {
			// Convert wallet address string to sdk.AccAddress
			walletAddr, err := sdk.AccAddressFromBech32(wallet.Address)
			if err != nil {
				continue // Skip invalid addresses
			}

			// Send tokens to the wallet
			coins := sdk.NewCoins(sdk.NewCoin("sclt", math.NewInt(int64(walletShare))))
			err = k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, walletAddr, coins)
			if err != nil {
				return err
			}
		}
	}

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

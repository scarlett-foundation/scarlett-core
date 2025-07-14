package keeper

import (
	"context"

	errorsmod "cosmossdk.io/errors"
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"scarlett-core/x/proofofdegen/types"
)

func (k msgServer) Claim(ctx context.Context, msg *types.MsgClaim) (*types.MsgClaimResponse, error) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)

	// 1. Validate creator address
	claimerAddr, err := k.addressCodec.StringToBytes(msg.Creator)
	if err != nil {
		return nil, errorsmod.Wrap(err, "invalid claimer address")
	}

	// 2. Validate that creator and address are the same (self-claiming only)
	if msg.Creator != msg.Address {
		return nil, errorsmod.Wrap(types.ErrUnauthorized, "can only claim for your own address")
	}

	// 3. Get and validate Merkle proof for address
	merkleProof, err := k.GetMerkleProofForAddress(ctx, msg.Address)
	if err != nil {
		return nil, errorsmod.Wrap(types.ErrNotEligible, "address is not eligible for airdrop")
	}

	// 4. Get campaign to check merkle root and claim status
	campaign, err := k.Campaign.Get(ctx, "genesis") // For now, hardcode to "genesis" campaign
	if err != nil {
		return nil, errorsmod.Wrap(err, "campaign not found")
	}

	// 5. Validate Merkle proof against campaign root
	weight, isValid := k.ValidateMerkleProof(msg.Address, merkleProof, campaign.MerkleRoot)
	if !isValid {
		return nil, errorsmod.Wrap(types.ErrInvalidProof, "invalid merkle proof for address")
	}

	// 6. Check if already claimed (we'll need to track this differently now)
	// For now, check if address exists in EligibleWallet and is claimed
	// TODO: Replace this with a more efficient claim tracking system
	eligibleWallet, err := k.EligibleWallet.Get(ctx, msg.Address)
	if err == nil && eligibleWallet.Claimed {
		return nil, errorsmod.Wrap(types.ErrAlreadyClaimed, "address has already claimed airdrop")
	}

	// 7. Calculate current claimable amount based on remaining claimers
	claimableAmount := k.CalculateShare(sdkCtx, msg.Address)
	if claimableAmount == 0 {
		return nil, errorsmod.Wrap(types.ErrNoTokensTolaim, "no tokens available to claim")
	}

	// 8. Create or update eligible wallet record with Merkle proof weight
	if err != nil {
		// Create new record if it doesn't exist
		eligibleWallet = types.EligibleWallet{
			Index:     msg.Address, // Use address as index
			Address:   msg.Address,
			Claimed:   true,
			ClaimTime: uint64(sdkCtx.BlockHeight()),
			Weight:    weight, // Use weight from Merkle proof
		}
	} else {
		// Update existing record
		eligibleWallet.Claimed = true
		eligibleWallet.ClaimTime = uint64(sdkCtx.BlockHeight())
		eligibleWallet.Weight = weight // Update with weight from Merkle proof
	}

	// 9. Update eligible wallet record
	if err := k.EligibleWallet.Set(ctx, msg.Address, eligibleWallet); err != nil {
		return nil, errorsmod.Wrap(err, "failed to update eligible wallet record")
	}

	// 8. Transfer tokens to claimer
	claimerAccAddr := sdk.AccAddress(claimerAddr)
	coins := sdk.NewCoins(sdk.NewCoin("sclt", math.NewInt(int64(claimableAmount))))

	err = k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, claimerAccAddr, coins)
	if err != nil {
		return nil, errorsmod.Wrap(err, "failed to transfer tokens to claimer")
	}

	// 9. Emit claim event
	sdkCtx.EventManager().EmitEvent(
		sdk.NewEvent(
			"airdrop_claimed",
			sdk.NewAttribute("claimer", msg.Address),
			sdk.NewAttribute("amount", coins.String()),
			sdk.NewAttribute("block_height", string(rune(sdkCtx.BlockHeight()))),
		),
	)

	return &types.MsgClaimResponse{
		ClaimedAmount: claimableAmount,
	}, nil
}

// CalculateShare calculates the claimable amount for a wallet based on current distribution
func (k Keeper) CalculateShare(ctx sdk.Context, walletAddress string) uint64 {
	// Get module balance (total emissions available for distribution)
	moduleAccount := k.authKeeper.GetModuleAccount(ctx, types.ModuleName)
	if moduleAccount == nil {
		return 0
	}

	moduleAddr := moduleAccount.GetAddress()
	moduleBalance := k.bankKeeper.GetBalance(ctx, moduleAddr, "sclt")
	availableTokens := moduleBalance.Amount.Uint64()

	if availableTokens == 0 {
		return 0
	}

	// Get all unclaimed wallets
	unclaimedWallets := k.GetUnclaimedWallets(ctx)
	if len(unclaimedWallets) == 0 {
		return 0
	}

	// Find this wallet in the unclaimed list
	var thisWallet *types.EligibleWallet
	totalWeight := uint64(0)

	for _, wallet := range unclaimedWallets {
		totalWeight += wallet.Weight
		if wallet.Address == walletAddress {
			thisWallet = &wallet
		}
	}

	if thisWallet == nil {
		return 0 // Wallet not found or already claimed
	}

	// Calculate this wallet's share: (availableTokens * walletWeight) / totalWeight
	if totalWeight == 0 {
		return 0
	}

	walletShare := (availableTokens * thisWallet.Weight) / totalWeight
	return walletShare
}

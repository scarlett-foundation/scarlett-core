package keeper

import (
	"context"

	"scarlett-core/x/proofofdegen/types"
)

// CalculateShareFromAggregates calculates claimable amount using aggregate tracking
// This replaces the O(n) iteration through individual wallets
func (k Keeper) CalculateShareFromAggregates(ctx context.Context, walletAddress string, walletWeight uint64) uint64 {
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

	// Get campaign to access aggregate data
	campaign, err := k.Campaign.Get(ctx, "genesis")
	if err != nil {
		return 0
	}

	// Calculate remaining unclaimed weight using aggregates
	// total_weight is stored in campaign (from Merkle tree generation)
	// claimed_weight is tracked as claims happen
	totalWeight := campaign.TotalAllocation // TODO: This should be total_weight, not total_allocation
	unclaimedWeight := totalWeight - campaign.ClaimedWeight

	if unclaimedWeight == 0 {
		return 0
	}

	// Calculate this wallet's share: (availableTokens * walletWeight) / unclaimedWeight
	walletShare := (availableTokens * walletWeight) / unclaimedWeight
	return walletShare
}

// GetAggregateStats returns campaign statistics using aggregate data
func (k Keeper) GetAggregateStats(ctx context.Context, campaignIndex string) (types.Campaign, uint64, uint64, error) {
	campaign, err := k.Campaign.Get(ctx, campaignIndex)
	if err != nil {
		return types.Campaign{}, 0, 0, err
	}

	// Get module balance
	moduleAccount := k.authKeeper.GetModuleAccount(ctx, types.ModuleName)
	var moduleBalance uint64
	if moduleAccount != nil {
		moduleAddr := moduleAccount.GetAddress()
		balance := k.bankKeeper.GetBalance(ctx, moduleAddr, "sclt")
		moduleBalance = balance.Amount.Uint64()
	}

	// Calculate unclaimed metrics using aggregates
	totalWeight := campaign.TotalAllocation // TODO: Replace with total_weight field
	unclaimedWeight := totalWeight - campaign.ClaimedWeight

	return campaign, moduleBalance, unclaimedWeight, nil
}

package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"scarlett-core/x/proofofdegen/types"
)

func (q queryServer) CampaignInfo(ctx context.Context, req *types.QueryCampaignInfoRequest) (*types.QueryCampaignInfoResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	sdkCtx := sdk.UnwrapSDKContext(ctx)

	// 1. Get current module balance (available for distribution)
	var currentModuleBalance uint64
	moduleAccount := q.k.authKeeper.GetModuleAccount(ctx, types.ModuleName)
	if moduleAccount != nil {
		moduleAddr := moduleAccount.GetAddress()
		balance := q.k.bankKeeper.GetBalance(ctx, moduleAddr, "sclt")
		currentModuleBalance = balance.Amount.Uint64()
	}

	// 2. Get all unclaimed wallets statistics
	unclaimedWallets := q.k.GetUnclaimedWallets(sdkCtx)
	totalUnclaimedWallets := uint64(len(unclaimedWallets))

	// Calculate total unclaimed weight
	totalUnclaimedWeight := uint64(0)
	for _, wallet := range unclaimedWallets {
		totalUnclaimedWeight += wallet.Weight
	}

	// 3. Get total claimed wallets count
	var totalClaimedWallets uint64
	err := q.k.EligibleWallet.Walk(ctx, nil, func(key string, wallet types.EligibleWallet) (bool, error) {
		if wallet.Claimed {
			totalClaimedWallets++
		}
		return false, nil
	})
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to count claimed wallets")
	}

	// 4. Calculate total emissions received (this would be tracked in a real implementation)
	// For now, we can estimate based on current balance + already distributed tokens
	// In a production system, this should be tracked in campaign storage
	totalEmissionsReceived := currentModuleBalance // Simplified for MVP

	// 5. Determine campaign status
	campaignStatus := "active"
	if totalUnclaimedWallets == 0 {
		campaignStatus = "completed"
	} else if currentModuleBalance == 0 {
		campaignStatus = "no_funds"
	}

	return &types.QueryCampaignInfoResponse{
		TotalEmissionsReceived: totalEmissionsReceived,
		TotalUnclaimedWallets:  totalUnclaimedWallets,
		TotalClaimedWallets:    totalClaimedWallets,
		TotalUnclaimedWeight:   totalUnclaimedWeight,
		CurrentModuleBalance:   currentModuleBalance,
		CampaignStatus:         campaignStatus,
	}, nil
}

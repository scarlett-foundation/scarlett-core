package keeper

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"scarlett-core/x/proofofdegen/types"
)

func (q queryServer) CampaignInfo(ctx context.Context, req *types.QueryCampaignInfoRequest) (*types.QueryCampaignInfoResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	// Use aggregate-based statistics for O(1) performance
	campaign, currentModuleBalance, totalUnclaimedWeight, err := q.k.GetAggregateStats(ctx, "genesis")
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to get campaign stats")
	}

	// Calculate metrics from campaign aggregates
	totalClaimedWallets := campaign.ClaimedCount
	totalWeight := campaign.TotalAllocation                    // TODO: Replace with total_weight when available
	totalUnclaimedWallets := totalWeight - totalClaimedWallets // Approximation

	// Calculate total emissions received
	totalEmissionsReceived := currentModuleBalance // Simplified for MVP

	// Determine campaign status
	campaignStatus := "active"
	if !campaign.Active {
		campaignStatus = "inactive"
	} else if totalUnclaimedWallets == 0 {
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

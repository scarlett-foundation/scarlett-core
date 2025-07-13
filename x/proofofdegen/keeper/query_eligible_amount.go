package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"scarlett-core/x/proofofdegen/types"
)

func (q queryServer) EligibleAmount(ctx context.Context, req *types.QueryEligibleAmountRequest) (*types.QueryEligibleAmountResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	if req.Address == "" {
		return nil, status.Error(codes.InvalidArgument, "address cannot be empty")
	}

	sdkCtx := sdk.UnwrapSDKContext(ctx)

	// 1. Check if address is eligible
	eligibleWallet, err := q.k.EligibleWallet.Get(ctx, req.Address)
	if err != nil {
		return &types.QueryEligibleAmountResponse{
			IsEligible:      false,
			ClaimableAmount: 0,
			Status:          "not_eligible",
		}, nil
	}

	// 2. Check if already claimed
	if eligibleWallet.Claimed {
		return &types.QueryEligibleAmountResponse{
			IsEligible:      true,
			ClaimableAmount: 0,
			Status:          "already_claimed",
			ClaimedAt:       eligibleWallet.ClaimTime,
		}, nil
	}

	// 3. Calculate current claimable amount using dynamic share calculation
	claimableAmount := q.k.CalculateShare(sdkCtx, req.Address)

	// 4. Get additional statistics for transparency
	unclaimedWallets := q.k.GetUnclaimedWallets(sdkCtx)
	totalUnclaimedWallets := uint64(len(unclaimedWallets))

	// Calculate user's weight percentage
	totalWeight := uint64(0)
	for _, wallet := range unclaimedWallets {
		totalWeight += wallet.Weight
	}

	var weightPercentage float64
	if totalWeight > 0 {
		weightPercentage = float64(eligibleWallet.Weight) / float64(totalWeight) * 100
	}

	return &types.QueryEligibleAmountResponse{
		IsEligible:            true,
		ClaimableAmount:       claimableAmount,
		Status:                "claimable",
		WalletWeight:          eligibleWallet.Weight,
		TotalUnclaimedWallets: totalUnclaimedWallets,
		WeightPercentage:      weightPercentage,
	}, nil
}

package keeper

import (
	"context"
	"errors"

	"scarlett-core/x/proofofdegen/types"

	"cosmossdk.io/collections"
	"github.com/cosmos/cosmos-sdk/types/query"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (q queryServer) ListEligibleWallet(ctx context.Context, req *types.QueryAllEligibleWalletRequest) (*types.QueryAllEligibleWalletResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	eligibleWallets, pageRes, err := query.CollectionPaginate(
		ctx,
		q.k.EligibleWallet,
		req.Pagination,
		func(_ string, value types.EligibleWallet) (types.EligibleWallet, error) {
			return value, nil
		},
	)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryAllEligibleWalletResponse{EligibleWallet: eligibleWallets, Pagination: pageRes}, nil
}

func (q queryServer) GetEligibleWallet(ctx context.Context, req *types.QueryGetEligibleWalletRequest) (*types.QueryGetEligibleWalletResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	val, err := q.k.EligibleWallet.Get(ctx, req.Index)
	if err != nil {
		if errors.Is(err, collections.ErrNotFound) {
			return nil, status.Error(codes.NotFound, "not found")
		}

		return nil, status.Error(codes.Internal, "internal error")
	}

	return &types.QueryGetEligibleWalletResponse{EligibleWallet: val}, nil
}

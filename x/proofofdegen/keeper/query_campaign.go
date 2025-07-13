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

func (q queryServer) ListCampaign(ctx context.Context, req *types.QueryAllCampaignRequest) (*types.QueryAllCampaignResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	campaigns, pageRes, err := query.CollectionPaginate(
		ctx,
		q.k.Campaign,
		req.Pagination,
		func(_ string, value types.Campaign) (types.Campaign, error) {
			return value, nil
		},
	)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryAllCampaignResponse{Campaign: campaigns, Pagination: pageRes}, nil
}

func (q queryServer) GetCampaign(ctx context.Context, req *types.QueryGetCampaignRequest) (*types.QueryGetCampaignResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	val, err := q.k.Campaign.Get(ctx, req.Index)
	if err != nil {
		if errors.Is(err, collections.ErrNotFound) {
			return nil, status.Error(codes.NotFound, "not found")
		}

		return nil, status.Error(codes.Internal, "internal error")
	}

	return &types.QueryGetCampaignResponse{Campaign: val}, nil
}

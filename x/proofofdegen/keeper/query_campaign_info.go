package keeper

import (
	"context"

	"scarlett-core/x/proofofdegen/types"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (q queryServer) CampaignInfo(ctx context.Context, req *types.QueryCampaignInfoRequest) (*types.QueryCampaignInfoResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	// TODO: Process the query

	return &types.QueryCampaignInfoResponse{}, nil
}

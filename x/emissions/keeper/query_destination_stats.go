package keeper

import (
	"context"

    "scarlett-core/x/emissions/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (q queryServer) DestinationStats(ctx context.Context,  req *types.QueryDestinationStatsRequest) (*types.QueryDestinationStatsResponse, error) {
	if req == nil {
        return nil, status.Error(codes.InvalidArgument, "invalid request")
    }

    // TODO: Process the query

	return &types.QueryDestinationStatsResponse{}, nil
}

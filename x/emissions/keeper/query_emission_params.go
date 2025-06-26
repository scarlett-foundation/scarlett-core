package keeper

import (
	"context"

	"scarlett-core/x/emissions/types"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (q queryServer) EmissionParams(ctx context.Context, req *types.QueryEmissionParamsRequest) (*types.QueryEmissionParamsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	// Get current emission parameters
	_, err := q.k.EmissionParams.Get(ctx)
	if err != nil {
		// Return default parameters if not found - will implement response fields in proto
	}

	return &types.QueryEmissionParamsResponse{}, nil
}

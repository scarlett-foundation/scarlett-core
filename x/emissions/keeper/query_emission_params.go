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

	// Get current emission parameters using helper method
	_, err := q.k.GetEmissionParams(ctx)
	if err != nil {
		// Return default parameters if not found - will add to proto response later
	}

	return &types.QueryEmissionParamsResponse{}, nil
}

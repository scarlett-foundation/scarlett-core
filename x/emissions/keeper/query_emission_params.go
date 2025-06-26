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

	// TODO: Process the query

	return &types.QueryEmissionParamsResponse{}, nil
}

package keeper

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"scarlett-core/x/contracts/types"
)

func (q queryServer) Params(ctx context.Context, req *types.QueryParamsRequest) (*types.QueryParamsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	// Since we're wrapping wasmd's keeper, we return default params
	// In a real implementation, you might want to store module-specific params
	params := types.DefaultParams()

	return &types.QueryParamsResponse{Params: params}, nil
}

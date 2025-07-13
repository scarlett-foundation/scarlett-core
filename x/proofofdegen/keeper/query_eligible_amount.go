package keeper

import (
	"context"

	"scarlett-core/x/proofofdegen/types"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (q queryServer) EligibleAmount(ctx context.Context, req *types.QueryEligibleAmountRequest) (*types.QueryEligibleAmountResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	// TODO: Process the query

	return &types.QueryEligibleAmountResponse{}, nil
}

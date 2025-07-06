package keeper

import (
	"context"

	"scarlett-core/x/contracts/types"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (q queryServer) SmartContractState(ctx context.Context, req *types.QuerySmartContractStateRequest) (*types.QuerySmartContractStateResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	// TODO: Process the query

	return &types.QuerySmartContractStateResponse{}, nil
}

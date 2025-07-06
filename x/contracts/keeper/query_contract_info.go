package keeper

import (
	"context"

	"scarlett-core/x/contracts/types"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (q queryServer) ContractInfo(ctx context.Context, req *types.QueryContractInfoRequest) (*types.QueryContractInfoResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	// TODO: Process the query

	return &types.QueryContractInfoResponse{}, nil
}

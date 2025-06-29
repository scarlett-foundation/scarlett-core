package keeper

import (
	"context"

	"scarlett-core/x/contracts/types"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (q queryServer) ListContracts(ctx context.Context, req *types.QueryListContractsRequest) (*types.QueryListContractsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	// TODO: In Phase 2, this will query actual contract storage
	// For V1, we'll return mock data to show the structure works
	contracts := []*types.ContractSummary{
		// This will be replaced with actual contract storage queries in Phase 2
		// Example structure:
		// {
		//     Address: "scarlett1abcd1234...",
		//     Label: "Demo Contract",
		//     Creator: "scarlett1creator...",
		//     IsRegistered: true,
		// },
	}

	// TODO: Phase 2 implementation will:
	// 1. Iterate through stored contracts
	// 2. Check registration status for each
	// 3. Populate ContractSummary for each deployed contract
	// 4. Return paginated results

	return &types.QueryListContractsResponse{
		Contracts: contracts,
	}, nil
}

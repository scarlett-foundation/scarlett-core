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

	var contracts []*types.ContractSummary

	// Iterate through all deployed contracts
	err := q.k.ContractInfo.Walk(ctx, nil, func(contractAddr string, contractInfo types.ContractInfo) (bool, error) {
		// Check if contract is registered for emissions
		isRegistered := q.k.emissionsKeeper.IsModuleRegistered(ctx, contractAddr)

		// Create contract summary
		summary := &types.ContractSummary{
			Address:      contractAddr,
			Label:        contractInfo.Label,
			Creator:      contractInfo.Creator,
			IsRegistered: isRegistered,
		}

		contracts = append(contracts, summary)
		return false, nil // Continue iteration
	})

	if err != nil {
		return nil, status.Error(codes.Internal, "failed to iterate contracts: "+err.Error())
	}

	return &types.QueryListContractsResponse{
		Contracts: contracts,
	}, nil
}

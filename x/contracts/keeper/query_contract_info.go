package keeper

import (
	"context"
	"strings"

	"scarlett-core/x/contracts/types"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (q queryServer) ContractInfo(ctx context.Context, req *types.QueryContractInfoRequest) (*types.QueryContractInfoResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	// Validate contract address format
	if strings.TrimSpace(req.ContractAddress) == "" {
		return nil, status.Error(codes.InvalidArgument, "contract address cannot be empty")
	}

	// Validate contract address format
	if _, err := q.k.addressCodec.StringToBytes(req.ContractAddress); err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid contract address format")
	}

	// TODO: In Phase 2, this will query actual contract storage and registration
	// For V1, we'll return a placeholder response showing the structure
	contractInfo := &types.ContractInfo{
		Address:                 req.ContractAddress,
		Label:                   "",    // Will be populated from contract deployment storage
		Creator:                 "",    // Will be populated from contract deployment storage
		IsRegistered:            false, // Will be checked against registration storage
		RegistrationName:        "",    // Will be populated if contract is registered
		RegistrationDescription: "",    // Will be populated if contract is registered
		RegistrationCreator:     "",    // Will be populated if contract is registered
	}

	// TODO: Phase 2 implementation will:
	// 1. Query contract deployment storage by address
	// 2. Check if contract exists (return not found if missing)
	// 3. Query registration storage to check if registered
	// 4. Populate all fields with actual data
	// 5. Return comprehensive contract information

	// For V1, return empty info to indicate "contract not found" behavior
	// This will be populated with actual data in Phase 2
	return &types.QueryContractInfoResponse{
		ContractInfo: contractInfo,
	}, nil
}

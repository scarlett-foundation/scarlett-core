package keeper

import (
	"context"
	"strings"

	"scarlett-core/x/contracts/types"
	emissionstypes "scarlett-core/x/emissions/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
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

	// Get contract deployment information
	contractInfo, err := q.k.ContractInfo.Get(ctx, req.ContractAddress)
	if err != nil {
		return nil, status.Error(codes.NotFound, "contract not found: "+req.ContractAddress)
	}

	// Check if contract is registered for emissions
	isRegistered := q.k.emissionsKeeper.IsModuleRegistered(ctx, req.ContractAddress)

	// Initialize response with deployment information
	response := &types.ContractInfo{
		Address:      req.ContractAddress,
		Label:        contractInfo.Label,
		Creator:      contractInfo.Creator,
		IsRegistered: isRegistered,
	}

	// If registered, fetch registration details
	if isRegistered {
		registeredModule, err := q.k.emissionsKeeper.GetRegisteredModule(ctx, req.ContractAddress)
		if err == nil {
			response.RegistrationName = extractNameFromLabel(contractInfo.Label)
			response.RegistrationDescription = registeredModule.Description
			response.RegistrationCreator = registeredModule.Creator
		}
		// If we can't get registration details, continue with what we have
		// The IsRegistered flag will still be true
	}

	return &types.QueryContractInfoResponse{
		ContractInfo: response,
	}, nil
}

// GetContractBalance returns the token balance of a contract address
func (q queryServer) GetContractBalance(ctx context.Context, contractAddress string) (sdk.Coins, error) {
	contractAddr, err := q.k.addressCodec.StringToBytes(contractAddress)
	if err != nil {
		return nil, err
	}

	// Get contract balance using bank keeper
	return q.k.bankKeeper.SpendableCoins(ctx, contractAddr), nil
}

// GetContractDeploymentInfo returns detailed deployment information for a contract
func (q queryServer) GetContractDeploymentInfo(ctx context.Context, contractAddress string) (*types.ContractInfo, error) {
	contractInfo, err := q.k.ContractInfo.Get(ctx, contractAddress)
	if err != nil {
		return nil, err
	}

	return &contractInfo, nil
}

// GetContractRegistrationStatus returns detailed registration status and history
func (q queryServer) GetContractRegistrationStatus(ctx context.Context, contractAddress string) (bool, *emissionstypes.RegisteredModule, error) {
	isRegistered := q.k.emissionsKeeper.IsModuleRegistered(ctx, contractAddress)

	if !isRegistered {
		return false, nil, nil
	}

	registeredModule, err := q.k.emissionsKeeper.GetRegisteredModule(ctx, contractAddress)
	if err != nil {
		return true, nil, err // Registered but can't get details
	}

	return true, &registeredModule, nil
}

// IsContractEligibleForEmissions checks if a contract meets all requirements for emission funding
func (q queryServer) IsContractEligibleForEmissions(ctx context.Context, contractAddress string) (bool, string, error) {
	// 1. Check if contract exists
	_, err := q.k.ContractInfo.Get(ctx, contractAddress)
	if err != nil {
		return false, "contract not deployed", nil
	}

	// 2. Check if contract is registered
	isRegistered := q.k.emissionsKeeper.IsModuleRegistered(ctx, contractAddress)
	if !isRegistered {
		return false, "contract not registered for emissions", nil
	}

	// 3. Check registration status
	registeredModule, err := q.k.emissionsKeeper.GetRegisteredModule(ctx, contractAddress)
	if err != nil {
		return false, "failed to get registration details", err
	}

	if registeredModule.Status != emissionstypes.ModuleStatusRegistered {
		return false, "contract registration not active", nil
	}

	return true, "contract eligible for emission funding", nil
}

// extractNameFromLabel extracts the contract name from the label which may have been modified during registration
func extractNameFromLabel(label string) string {
	// Labels are stored as "name (Registered for Emissions)" during registration
	// Extract the original name
	if strings.Contains(label, " (Registered for Emissions)") {
		return strings.TrimSuffix(label, " (Registered for Emissions)")
	}
	return label
}

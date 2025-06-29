package keeper

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"cosmossdk.io/collections"
	"github.com/cosmos/cosmos-sdk/types/query"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"scarlett-core/x/emissions/types"
)

// ListRegisteredModules returns all registered modules in the emissions registry
func (q queryServer) ListRegisteredModules(ctx context.Context, req *types.QueryListRegisteredModulesRequest) (*types.QueryListRegisteredModulesResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	var modules []*types.RegisteredModuleInfo

	results, pageRes, err := query.CollectionPaginate(
		ctx,
		q.k.ModuleRegistryRaw,
		req.Pagination,
		func(key string, value string) (*types.RegisteredModuleInfo, error) {
			// Parse the stored module from JSON
			registeredModule, err := q.parseRegisteredModuleFromJSON(value)
			if err != nil {
				return nil, fmt.Errorf("failed to parse registered module '%s': %w", key, err)
			}

			// Convert to proto message format with enhanced contract detection
			moduleInfo := q.registeredModuleToInfoEnhanced(ctx, registeredModule)
			return &moduleInfo, nil
		},
	)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	// Extract modules from pagination result
	modules = append(modules, results...)

	return &types.QueryListRegisteredModulesResponse{
		Modules:    modules,
		Pagination: pageRes,
	}, nil
}

// GetRegisteredModule returns a specific registered module by name
func (q queryServer) GetRegisteredModule(ctx context.Context, req *types.QueryGetRegisteredModuleRequest) (*types.QueryGetRegisteredModuleResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	if req.ModuleName == "" {
		return nil, status.Error(codes.InvalidArgument, "module name cannot be empty")
	}

	// Get the registered module from storage
	registeredModule, err := q.k.GetRegisteredModule(ctx, req.ModuleName)
	if err != nil {
		if err == collections.ErrNotFound {
			return nil, status.Error(codes.NotFound, fmt.Sprintf("module '%s' not found in registry", req.ModuleName))
		}
		return nil, status.Error(codes.Internal, err.Error())
	}

	// Convert to proto message format with enhanced contract information
	moduleInfo := q.registeredModuleToInfoEnhanced(ctx, registeredModule)

	return &types.QueryGetRegisteredModuleResponse{
		Module: &moduleInfo,
	}, nil
}

// Helper function to parse RegisteredModule from JSON string
func (q queryServer) parseRegisteredModuleFromJSON(jsonStr string) (types.RegisteredModule, error) {
	var module types.RegisteredModule
	if err := json.Unmarshal([]byte(jsonStr), &module); err != nil {
		return types.RegisteredModule{}, err
	}
	return module, nil
}

// Helper function to convert RegisteredModule to RegisteredModuleInfo proto message
func (q queryServer) registeredModuleToInfo(module types.RegisteredModule) types.RegisteredModuleInfo {
	return types.RegisteredModuleInfo{
		ModuleName:  module.ModuleName,
		Creator:     module.Creator,
		Description: module.Description,
		Status:      string(module.Status),
		CreatedAt:   module.CreatedAt.Format("2006-01-02T15:04:05Z"),
		UpdatedAt:   module.UpdatedAt.Format("2006-01-02T15:04:05Z"),
	}
}

// Helper function to determine if a module name is actually a contract address
func (q queryServer) isContractAddress(moduleName string) bool {
	// Contract addresses start with "scarlett1" and are longer than typical module names
	return strings.HasPrefix(moduleName, "scarlett1") && len(moduleName) >= 45
}

// Enhanced helper function to convert RegisteredModule to RegisteredModuleInfo with contract detection
func (q queryServer) registeredModuleToInfoEnhanced(ctx context.Context, module types.RegisteredModule) types.RegisteredModuleInfo {
	info := types.RegisteredModuleInfo{
		ModuleName:  module.ModuleName,
		Creator:     module.Creator,
		Description: module.Description,
		Status:      string(module.Status),
		CreatedAt:   module.CreatedAt.Format("2006-01-02T15:04:05Z"),
		UpdatedAt:   module.UpdatedAt.Format("2006-01-02T15:04:05Z"),
	}

	// Add contract-specific information if this is a contract
	if q.isContractAddress(module.ModuleName) {
		info.IsContract = true
		info.ContractAddress = module.ModuleName

		// Enhance description to indicate this is a contract
		if !strings.Contains(info.Description, "(Contract)") {
			info.Description = fmt.Sprintf("%s (Contract)", info.Description)
		}
	} else {
		info.IsContract = false
		info.ContractAddress = ""

		// Enhance description to indicate this is a module
		if !strings.Contains(info.Description, "(Module)") {
			info.Description = fmt.Sprintf("%s (Module)", info.Description)
		}
	}

	return info
}

// GetContractFundingStatus returns detailed funding information for a contract
func (q queryServer) GetContractFundingStatus(ctx context.Context, contractAddress string) (*types.ContractFundingInfo, error) {
	// Check if contract is registered
	registeredModule, err := q.k.GetRegisteredModule(ctx, contractAddress)
	if err != nil {
		return &types.ContractFundingInfo{
			ContractAddress: contractAddress,
			IsRegistered:    false,
			IsEligible:      false,
			Status:          "not_registered",
		}, nil
	}

	// Get current emission configuration
	emissionParams, err := q.k.GetEmissionParams(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get emission parameters: %w", err)
	}

	// Calculate funding information
	isEligible := registeredModule.Status == types.ModuleStatusRegistered
	fundingWeight := float64(0)

	// Find the contract in current emission destinations
	for _, destination := range emissionParams.Destinations {
		if destination.ModuleName == contractAddress {
			fundingWeight, _ = destination.Weight.Float64()
			break
		}
	}

	return &types.ContractFundingInfo{
		ContractAddress:  contractAddress,
		IsRegistered:     true,
		IsEligible:       isEligible,
		Status:           string(registeredModule.Status),
		CurrentWeight:    fundingWeight,
		RegistrationName: registeredModule.ModuleName,
		Description:      registeredModule.Description,
		Creator:          registeredModule.Creator,
		RegisteredAt:     registeredModule.CreatedAt.Format("2006-01-02T15:04:05Z"),
		LastUpdated:      registeredModule.UpdatedAt.Format("2006-01-02T15:04:05Z"),
	}, nil
}

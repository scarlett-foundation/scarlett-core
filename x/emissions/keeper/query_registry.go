package keeper

import (
	"context"
	"encoding/json"
	"fmt"

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

			// Convert to proto message format
			moduleInfo := q.registeredModuleToInfo(registeredModule)
			return &moduleInfo, nil
		},
	)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	// Extract modules from pagination result
	for _, result := range results {
		modules = append(modules, result)
	}

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

	// Convert to proto message format
	moduleInfo := q.registeredModuleToInfo(registeredModule)

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

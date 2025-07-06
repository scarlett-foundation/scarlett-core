package keeper

import (
	"context"

	"scarlett-core/x/contracts/types"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (q queryServer) SmartContractState(goCtx context.Context, req *types.QuerySmartContractStateRequest) (*types.QuerySmartContractStateResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	// Validate contract address
	if req.Contract == "" {
		return nil, status.Error(codes.InvalidArgument, "contract address cannot be empty")
	}

	// Validate query data
	if len(req.QueryData) == 0 {
		return nil, status.Error(codes.InvalidArgument, "query data cannot be empty")
	}

	// Parse contract address
	contractAddr, err := sdk.AccAddressFromBech32(req.Contract)
	if err != nil {
		return nil, errorsmod.Wrap(err, "invalid contract address")
	}

	// Get the wasm keeper with error handling
	wasmKeeper, err := q.k.getWasmKeeper()
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	// Execute smart contract query via wasmd
	queryResult, err := wasmKeeper.QuerySmart(ctx, contractAddr, req.QueryData)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	// Create response with query result data
	response := &types.QuerySmartContractStateResponse{
		Data: queryResult,
	}

	// Log the successful query
	q.k.Logger(ctx).Info("Successfully executed smart contract query",
		"contract_address", req.Contract,
		"query_data_length", len(req.QueryData),
		"result_length", len(queryResult))

	return response, nil
}

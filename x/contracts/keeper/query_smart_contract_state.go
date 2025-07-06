package keeper

import (
	"context"

	"scarlett-core/x/contracts/types"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (q queryServer) SmartContractState(ctx context.Context, req *types.QuerySmartContractStateRequest) (*types.QuerySmartContractStateResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

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

	// Get the wasm keeper
	wasmKeeper := q.k.getWasmKeeper()

	// Execute smart contract query
	queryResult, err := wasmKeeper.QuerySmart(ctx, contractAddr, req.QueryData)
	if err != nil {
		return nil, errorsmod.Wrap(err, "failed to execute smart contract query")
	}

	// Log the successful query
	q.k.Logger(ctx).Info("Successfully executed smart contract query",
		"contract_address", req.Contract,
		"query_data_length", len(req.QueryData),
		"result_length", len(queryResult))

	// For now, return empty response since proto doesn't have fields defined
	// TODO: Update proto to include query result fields and populate response
	// The queryResult contains the actual JSON response from the contract
	return &types.QuerySmartContractStateResponse{}, nil
}

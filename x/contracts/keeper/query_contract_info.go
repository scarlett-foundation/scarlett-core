package keeper

import (
	"context"

	"scarlett-core/x/contracts/types"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (q queryServer) ContractInfo(ctx context.Context, req *types.QueryContractInfoRequest) (*types.QueryContractInfoResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	// Validate contract address
	if req.Contract == "" {
		return nil, status.Error(codes.InvalidArgument, "contract address cannot be empty")
	}

	// Parse contract address
	contractAddr, err := sdk.AccAddressFromBech32(req.Contract)
	if err != nil {
		return nil, errorsmod.Wrap(err, "invalid contract address")
	}

	// Get the wasm keeper
	wasmKeeper := q.k.getWasmKeeper()

	// Query contract info from wasmd keeper
	contractInfo := wasmKeeper.GetContractInfo(ctx, contractAddr)
	if contractInfo == nil {
		return nil, status.Error(codes.NotFound, "contract not found")
	}

	// Log the successful query
	q.k.Logger(ctx).Info("Successfully queried contract info",
		"contract_address", req.Contract,
		"code_id", contractInfo.CodeID,
		"creator", contractInfo.Creator,
		"admin", contractInfo.Admin)

	// For now, return empty response since proto doesn't have fields defined
	// TODO: Update proto to include contract info fields and populate response
	return &types.QueryContractInfoResponse{}, nil
}

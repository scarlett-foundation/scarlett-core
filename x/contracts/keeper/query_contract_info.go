package keeper

import (
	"context"

	"scarlett-core/x/contracts/types"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (q queryServer) ContractInfo(goCtx context.Context, req *types.QueryContractInfoRequest) (*types.QueryContractInfoResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	// Validate contract address
	if req.Contract == "" {
		return nil, status.Error(codes.InvalidArgument, "contract address cannot be empty")
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

	// Query contract info from wasmd
	contractInfo := wasmKeeper.GetContractInfo(ctx, contractAddr)
	if contractInfo == nil {
		return nil, status.Error(codes.NotFound, "contract not found")
	}

	// Map wasmd ContractInfo to our proto response
	response := &types.QueryContractInfoResponse{
		ContractAddress: req.Contract,
		Creator:         contractInfo.Creator,
		Admin:           contractInfo.Admin,
		CodeId:          contractInfo.CodeID,
		Label:           contractInfo.Label,
		IbcPortId:       contractInfo.IBCPortID,
	}

	// Map created position if available
	if contractInfo.Created != nil {
		response.Created = &types.AbsoluteTxPosition{
			BlockHeight: contractInfo.Created.BlockHeight,
			TxIndex:     contractInfo.Created.TxIndex,
		}
	}

	// Log the successful query
	q.k.Logger(ctx).Info("Successfully queried contract info",
		"contract_address", req.Contract,
		"code_id", contractInfo.CodeID,
		"creator", contractInfo.Creator,
		"admin", contractInfo.Admin)

	return response, nil
}

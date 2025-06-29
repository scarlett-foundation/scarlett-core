package keeper

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"

	"scarlett-core/x/contracts/types"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k msgServer) DeployContract(ctx context.Context, msg *types.MsgDeployContract) (*types.MsgDeployContractResponse, error) {
	if _, err := k.addressCodec.StringToBytes(msg.Creator); err != nil {
		return nil, errorsmod.Wrap(err, "invalid creator address")
	}

	// Validate contract code is not empty
	if len(msg.Code) == 0 {
		return nil, errorsmod.Wrap(types.ErrInvalidInput, "contract code cannot be empty")
	}

	// Validate label is not empty
	if msg.Label == "" {
		return nil, errorsmod.Wrap(types.ErrInvalidInput, "contract label cannot be empty")
	}

	// Store contract code and get checksum
	checksum, err := k.storeContractCode(ctx, msg.Code)
	if err != nil {
		return nil, errorsmod.Wrap(err, "failed to store contract code")
	}

	// Generate contract address
	contractAddr := k.generateContractAddress(ctx, msg.Creator, checksum)

	// Instantiate contract with the provided message
	err = k.instantiateContract(ctx, checksum, contractAddr, msg.Creator, msg.InstantiateMsg, msg.Label)
	if err != nil {
		return nil, errorsmod.Wrap(err, "failed to instantiate contract")
	}

	// Emit contract deployment event
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	sdkCtx.EventManager().EmitEvent(
		sdk.NewEvent(
			"contract_deployed",
			sdk.NewAttribute("creator", msg.Creator),
			sdk.NewAttribute("contract_address", contractAddr),
			sdk.NewAttribute("code_checksum", hex.EncodeToString(checksum)),
			sdk.NewAttribute("label", msg.Label),
		),
	)

	return &types.MsgDeployContractResponse{
		ContractAddress: contractAddr,
	}, nil
}

// storeContractCode validates and stores the contract code, returning the checksum
func (k msgServer) storeContractCode(ctx context.Context, code []byte) ([]byte, error) {
	// Create checksum of the code
	hash := sha256.Sum256(code)
	checksum := hash[:]

	// TODO: In a full implementation, we would:
	// 1. Initialize a wasmvm instance
	// 2. Validate the wasm code
	// 3. Store the code in the VM
	// For now, we'll just store the checksum and validate basic wasm format

	// Basic wasm magic number validation
	if len(code) < 8 || string(code[:4]) != "\x00asm" {
		return nil, errorsmod.Wrap(types.ErrInvalidInput, "invalid wasm format")
	}

	// Store code mapping in keeper state (checksum -> code)
	// This would be implemented with proper state management

	return checksum, nil
}

// generateContractAddress creates a deterministic contract address
func (k msgServer) generateContractAddress(ctx context.Context, creator string, checksum []byte) string {
	// Get current block height for uniqueness
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	height := sdkCtx.BlockHeight()

	// Create deterministic address from creator + checksum + height
	data := fmt.Sprintf("%s%s%d", creator, hex.EncodeToString(checksum), height)
	hash := sha256.Sum256([]byte(data))

	// Use first 20 bytes as address (similar to Ethereum)
	addr := hash[:20]

	// Convert to bech32 format (this is simplified)
	return fmt.Sprintf("scarlett1%s", hex.EncodeToString(addr)[:38])
}

// instantiateContract initializes the contract with the given parameters
func (k msgServer) instantiateContract(ctx context.Context, checksum []byte, contractAddr, creator, instantiateMsg, label string) error {
	// TODO: In a full implementation, we would:
	// 1. Call wasmvm.Instantiate with proper environment
	// 2. Set up contract storage
	// 3. Execute the instantiate function
	// 4. Store contract metadata

	// For now, validate that instantiate message is valid JSON-like
	if instantiateMsg != "" && instantiateMsg != "{}" {
		// Basic validation that it looks like JSON
		if !(instantiateMsg[0] == '{' && instantiateMsg[len(instantiateMsg)-1] == '}') {
			return errorsmod.Wrap(types.ErrInvalidInput, "instantiate message must be valid JSON object")
		}
	}

	// Store contract metadata (this would be implemented with proper state management)
	// - Contract address -> checksum mapping
	// - Contract info (creator, label, etc.)

	return nil
}

package keeper

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"

	"scarlett-core/x/contracts/types"

	wasmvm "github.com/CosmWasm/wasmvm/v3"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k msgServer) DeployContract(ctx context.Context, msg *types.MsgDeployContract) (*types.MsgDeployContractResponse, error) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	logger := sdkCtx.Logger()
	logger.Info("🚀 DeployContract handler initiated")

	if _, err := k.addressCodec.StringToBytes(msg.Creator); err != nil {
		return nil, errorsmod.Wrap(err, "invalid creator address")
	}

	// Validate contract code is not empty
	if len(msg.Code) == 0 {
		return nil, errorsmod.Wrap(types.ErrInvalidInput, "contract code cannot be empty")
	}

	// Get gas configuration
	gasConfig := k.GetGasConfig()
	logger.Info("✅ Gas configuration loaded", "config", fmt.Sprintf("%+v", gasConfig))

	// Validate contract code size
	if err := gasConfig.ValidateContractSize(msg.Code); err != nil {
		logger.Error("❌ Contract size validation failed", "error", err)
		return nil, errorsmod.Wrap(types.ErrInvalidInput, err.Error())
	}

	// Validate contract code is valid WASM
	if !isValidWasmBytecode(msg.Code) {
		logger.Error("❌ Invalid WASM bytecode format")
		return nil, errorsmod.Wrap(types.ErrInvalidInput, "invalid WASM bytecode format")
	}

	// Validate label is not empty
	if len(msg.Label) == 0 {
		return nil, errorsmod.Wrap(types.ErrInvalidInput, "contract label cannot be empty")
	}

	// Store contract code in WasmVM and get checksum
	wasmCode := wasmvm.WasmCode(msg.Code)
	deploymentGasLimit := gasConfig.DeploymentGasLimit // Use configurable gas limit (50M gas)
	logger.Info("⛽️ Storing contract code in WasmVM", "gas_limit", deploymentGasLimit)
	vm, err := k.GetWasmVM()
	if err != nil {
		logger.Error("❌ Failed to get WasmVM", "error", err)
		return nil, errorsmod.Wrapf(types.ErrInvalidInput, "failed to get WasmVM: %s", err)
	}
	checksum, gasUsed, err := vm.StoreCode(wasmCode, deploymentGasLimit)
	if err != nil {
		logger.Error("❌ Failed to store contract code", "error", err, "gas_used", gasUsed)
		return nil, errorsmod.Wrapf(types.ErrInvalidInput, "failed to store contract code: %s", err)
	}
	logger.Info("✅ Contract code stored successfully", "gas_used", gasUsed, "checksum", hex.EncodeToString(checksum))
	_ = gasUsed // We'll use this for gas accounting in future versions

	// Generate next contract ID
	contractSeq, err := k.ContractSeq.Next(ctx)
	if err != nil {
		return nil, errorsmod.Wrap(err, "failed to generate contract sequence")
	}

	// Generate deterministic contract address
	contractAddress := generateContractAddress(msg.Creator, contractSeq)

	// Store contract code mapping (checksum -> bytecode)
	checksumBytes := []byte(checksum)
	if err := k.ContractCode.Set(ctx, checksumBytes, msg.Code); err != nil {
		return nil, errorsmod.Wrap(err, "failed to store contract code mapping")
	}

	// Create contract info
	contractInfo := types.ContractInfo{
		Address:                 contractAddress,
		Label:                   msg.Label,
		Creator:                 msg.Creator,
		IsRegistered:            false, // Not registered for emissions yet
		RegistrationName:        "",
		RegistrationDescription: "",
		RegistrationCreator:     "",
	}

	// Store contract info
	if err := k.ContractInfo.Set(ctx, contractAddress, contractInfo); err != nil {
		return nil, errorsmod.Wrap(err, "failed to store contract info")
	}

	// Emit deployment event
	sdkCtx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			"contract_deployed",
			sdk.NewAttribute("contract_address", contractAddress),
			sdk.NewAttribute("creator", msg.Creator),
			sdk.NewAttribute("label", msg.Label),
			sdk.NewAttribute("code_id", fmt.Sprintf("%d", contractSeq)),
			sdk.NewAttribute("checksum", hex.EncodeToString(checksumBytes)),
		),
	})

	return &types.MsgDeployContractResponse{
		ContractAddress: contractAddress,
	}, nil
}

// isValidWasmBytecode performs basic validation to check if bytes represent valid WASM
func isValidWasmBytecode(code []byte) bool {
	// WASM files start with magic number: 0x00 0x61 0x73 0x6D (null + "asm")
	if len(code) < 4 {
		return false
	}
	return code[0] == 0x00 && code[1] == 0x61 && code[2] == 0x73 && code[3] == 0x6D
}

// generateContractAddress creates a deterministic contract address
func generateContractAddress(creator string, contractSeq uint64) string {
	// Create a deterministic address based on creator and sequence
	// Format: sha256(creator + contractSeq) -> bech32 address
	input := fmt.Sprintf("%s:%d", creator, contractSeq)
	hash := sha256.Sum256([]byte(input))

	// Take first 20 bytes for address (standard cosmos address length)
	addrBytes := hash[:20]

	// Convert to bech32 format with "scarlett" prefix
	// Note: In production, you'd use the proper bech32 encoding
	// For now, we'll create a deterministic hex-based address
	return fmt.Sprintf("scarlett1%s", hex.EncodeToString(addrBytes)[:38])
}

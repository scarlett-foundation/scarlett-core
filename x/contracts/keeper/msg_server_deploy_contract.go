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
	logger.Info("🚀 DeployContract handler initiated",
		"creator", msg.Creator,
		"contract_size_bytes", len(msg.Code),
		"label", msg.Label)

	if _, err := k.addressCodec.StringToBytes(msg.Creator); err != nil {
		logger.Error("❌ Invalid creator address", "creator", msg.Creator, "error", err)
		return nil, errorsmod.Wrap(err, "invalid creator address")
	}

	// Validate contract code is not empty
	if len(msg.Code) == 0 {
		logger.Error("❌ Empty contract code provided")
		return nil, errorsmod.Wrap(types.ErrInvalidInput, "contract code cannot be empty")
	}

	// Get gas configuration
	gasConfig := k.GetGasConfig()
	logger.Info("✅ Gas configuration loaded",
		"deployment_gas_limit", gasConfig.DeploymentGasLimit,
		"max_contract_size", gasConfig.MaxContractSize,
		"max_gas_limit", gasConfig.MaxGasLimit)

	// Validate contract code size
	if err := gasConfig.ValidateContractSize(msg.Code); err != nil {
		logger.Error("❌ Contract size validation failed",
			"contract_size", len(msg.Code),
			"max_allowed", gasConfig.MaxContractSize,
			"error", err)
		return nil, errorsmod.Wrap(types.ErrInvalidInput, err.Error())
	}

	// Validate contract code is valid WASM
	if !isValidWasmBytecode(msg.Code) {
		logger.Error("❌ Invalid WASM bytecode format",
			"contract_size", len(msg.Code),
			"first_4_bytes", fmt.Sprintf("%x", msg.Code[:4]))
		return nil, errorsmod.Wrap(types.ErrInvalidInput, "invalid WASM bytecode format")
	}

	// Validate label is not empty
	if len(msg.Label) == 0 {
		logger.Error("❌ Empty contract label provided")
		return nil, errorsmod.Wrap(types.ErrInvalidInput, "contract label cannot be empty")
	}

	logger.Info("✅ All validation checks passed",
		"contract_size", len(msg.Code),
		"label", msg.Label)

	// Store contract code in WasmVM and get checksum
	wasmCode := wasmvm.WasmCode(msg.Code)

	// Calculate gas needed based on contract size
	// Use 1 gas per byte as base cost, plus 50M base cost for compilation
	gasNeeded := uint64(len(msg.Code)) + 50_000_000

	// Add 20% buffer for safety
	deploymentGasLimit := gasNeeded + (gasNeeded / 5)

	// Ensure we don't exceed max gas limit
	if deploymentGasLimit > gasConfig.MaxGasLimit {
		deploymentGasLimit = gasConfig.MaxGasLimit
	}

	logger.Info("⛽️ Preparing to store contract code in WasmVM",
		"contract_size", len(msg.Code),
		"gas_needed", gasNeeded,
		"gas_limit", deploymentGasLimit)

	vm, err := k.GetWasmVM(logger)
	if err != nil {
		logger.Error("❌ Failed to get WasmVM instance",
			"error", err,
			"gas_limit", deploymentGasLimit)
		return nil, errorsmod.Wrapf(types.ErrInvalidInput, "failed to get WasmVM: %s", err)
	}

	logger.Info("✅ WasmVM instance obtained, calling StoreCode",
		"vm_initialized", true,
		"gas_limit", deploymentGasLimit)

	checksum, gasUsed, err := vm.StoreCode(wasmCode, deploymentGasLimit)
	if err != nil {
		logger.Error("❌ Failed to store contract code in WasmVM",
			"error", err,
			"gas_used", gasUsed,
			"gas_limit", deploymentGasLimit,
			"contract_size", len(msg.Code),
			"gas_efficiency", fmt.Sprintf("%.2f bytes/gas", float64(len(msg.Code))/float64(gasUsed)))
		return nil, errorsmod.Wrapf(types.ErrInvalidInput, "failed to store contract code: %s", err)
	}

	logger.Info("✅ Contract code stored successfully in WasmVM",
		"gas_used", gasUsed,
		"gas_limit", deploymentGasLimit,
		"gas_remaining", deploymentGasLimit-gasUsed,
		"gas_efficiency", fmt.Sprintf("%.2f bytes/gas", float64(len(msg.Code))/float64(gasUsed)),
		"checksum", hex.EncodeToString(checksum))

	// Generate next contract ID
	contractSeq, err := k.ContractSeq.Next(ctx)
	if err != nil {
		logger.Error("❌ Failed to generate contract sequence", "error", err)
		return nil, errorsmod.Wrap(err, "failed to generate contract sequence")
	}

	logger.Info("✅ Contract sequence generated", "contract_seq", contractSeq)

	// Generate deterministic contract address
	contractAddress := generateContractAddress(msg.Creator, contractSeq)
	logger.Info("✅ Contract address generated",
		"contract_address", contractAddress,
		"creator", msg.Creator,
		"sequence", contractSeq)

	// Store contract code mapping (checksum -> bytecode)
	checksumBytes := []byte(checksum)
	if err := k.ContractCode.Set(ctx, checksumBytes, msg.Code); err != nil {
		logger.Error("❌ Failed to store contract code mapping",
			"error", err,
			"checksum", hex.EncodeToString(checksumBytes))
		return nil, errorsmod.Wrap(err, "failed to store contract code mapping")
	}

	logger.Info("✅ Contract code mapping stored",
		"checksum", hex.EncodeToString(checksumBytes),
		"code_size", len(msg.Code))

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
		logger.Error("❌ Failed to store contract info",
			"error", err,
			"contract_address", contractAddress)
		return nil, errorsmod.Wrap(err, "failed to store contract info")
	}

	logger.Info("✅ Contract info stored successfully",
		"contract_address", contractAddress,
		"label", msg.Label,
		"creator", msg.Creator)

	// Emit deployment event
	sdkCtx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			"contract_deployed",
			sdk.NewAttribute("contract_address", contractAddress),
			sdk.NewAttribute("creator", msg.Creator),
			sdk.NewAttribute("label", msg.Label),
			sdk.NewAttribute("code_id", fmt.Sprintf("%d", contractSeq)),
			sdk.NewAttribute("checksum", hex.EncodeToString(checksumBytes)),
			sdk.NewAttribute("gas_used", fmt.Sprintf("%d", gasUsed)),
			sdk.NewAttribute("contract_size", fmt.Sprintf("%d", len(msg.Code))),
		),
	})

	logger.Info("🎉 CONTRACT DEPLOYMENT COMPLETED SUCCESSFULLY",
		"contract_address", contractAddress,
		"creator", msg.Creator,
		"label", msg.Label,
		"gas_used", gasUsed,
		"contract_size", len(msg.Code),
		"checksum", hex.EncodeToString(checksumBytes))

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

package keeper

import (
	"context"
	"fmt"
	"time"

	"scarlett-core/x/contracts/types"

	wasmvmtypes "github.com/CosmWasm/wasmvm/v3/types"

	corestore "cosmossdk.io/core/store"
	errorsmod "cosmossdk.io/errors"
	"cosmossdk.io/log"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Note: Security configuration is now managed via types.GasConfig
// This provides centralized configuration with embedded defaults for consensus safety

// ContractExecutionParams holds parameters for secure contract execution
type ContractExecutionParams struct {
	GasLimit  uint64
	GasPrice  sdk.DecCoins
	Timeout   time.Duration
	Caller    string
	Funds     sdk.Coins
	BlockInfo wasmvmtypes.BlockInfo
	TxInfo    wasmvmtypes.TransactionInfo
}

// ValidateContractExecution performs security validation before contract execution
func (k Keeper) ValidateContractExecution(ctx context.Context, contractAddr string, params ContractExecutionParams) error {
	// Validate contract address format
	if _, err := k.addressCodec.StringToBytes(contractAddr); err != nil {
		return errorsmod.Wrapf(types.ErrInvalidInput, "invalid contract address: %s", err)
	}

	// Validate contract exists
	_, err := k.ContractInfo.Get(ctx, contractAddr)
	if err != nil {
		return errorsmod.Wrapf(types.ErrInvalidInput, "contract not found: %s", contractAddr)
	}

	// Validate gas limits
	if err := k.validateGasLimits(params.GasLimit); err != nil {
		return err
	}

	// Validate execution timeout
	if err := k.validateTimeout(params.Timeout); err != nil {
		return err
	}

	// Validate caller permissions
	if err := k.validateCallerPermissions(ctx, params.Caller, contractAddr); err != nil {
		return err
	}

	// Validate funds if any
	if err := k.validateFunds(ctx, params.Caller, params.Funds); err != nil {
		return err
	}

	return nil
}

// validateGasLimits ensures gas limits are within acceptable bounds
func (k Keeper) validateGasLimits(gasLimit uint64) error {
	return k.gasConfig.ValidateGasLimit(gasLimit)
}

// validateTimeout ensures execution timeout is within acceptable bounds
func (k Keeper) validateTimeout(timeout time.Duration) error {
	return k.gasConfig.ValidateTimeout(timeout)
}

// validateCallerPermissions checks if caller has permission to execute contract
func (k Keeper) validateCallerPermissions(ctx context.Context, caller, contractAddr string) error {
	// For V1, we allow any valid address to call any contract
	// Future versions may implement more sophisticated permission systems

	if _, err := k.addressCodec.StringToBytes(caller); err != nil {
		return errorsmod.Wrapf(types.ErrInvalidInput, "invalid caller address: %s", err)
	}

	return nil
}

// validateFunds ensures caller has sufficient funds for the transaction
func (k Keeper) validateFunds(ctx context.Context, caller string, funds sdk.Coins) error {
	if funds.IsZero() {
		return nil // No funds to validate
	}

	// Check if funds are valid denominations
	if err := funds.Validate(); err != nil {
		return errorsmod.Wrapf(types.ErrInvalidInput, "invalid funds: %s", err)
	}

	// In production, we would check caller's actual balance here
	// For now, we just validate the coin format

	return nil
}

// CreateSecureExecutionEnvironment sets up a secure environment for contract execution
func (k Keeper) CreateSecureExecutionEnvironment(ctx context.Context, contractAddr string) (wasmvmtypes.Env, error) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)

	// Create block info with current chain state
	blockInfo := wasmvmtypes.BlockInfo{
		Height:  uint64(sdkCtx.BlockHeight()),
		Time:    wasmvmtypes.Uint64(sdkCtx.BlockTime().Unix()),
		ChainID: sdkCtx.ChainID(),
	}

	// Create contract environment
	env := wasmvmtypes.Env{
		Block: blockInfo,
		Contract: wasmvmtypes.ContractInfo{
			Address: contractAddr,
		},
		Transaction: &wasmvmtypes.TransactionInfo{
			Index: 0, // This would be set properly in full implementation
		},
	}

	return env, nil
}

// GasConversionFactor is how many WasmVM gas units equal one SDK gas unit
// WasmVM uses much larger gas units than Cosmos SDK
const GasConversionFactor = 100_000

// CreateGasMeter creates a gas meter for contract execution
func (k Keeper) CreateGasMeter(ctx context.Context, gasLimit uint64) wasmvmtypes.GasMeter {
	// Convert SDK gas limit to WasmVM gas units
	wasmGasLimit := gasLimit * GasConversionFactor

	sdkCtx := sdk.UnwrapSDKContext(ctx)
	logger := sdkCtx.Logger()
	logger.Info("⛽️ Creating WasmVM gas meter",
		"sdk_gas_limit", gasLimit,
		"conversion_factor", GasConversionFactor,
		"wasm_gas_limit", wasmGasLimit)

	return &WasmGasMeter{
		limit:    wasmGasLimit,
		consumed: 0,
		logger:   logger,
	}
}

// WasmGasMeter implements proper gas conversion between WasmVM and Cosmos SDK
type WasmGasMeter struct {
	limit    uint64 // In WasmVM gas units
	consumed uint64 // In WasmVM gas units
	logger   log.Logger
}

func (gm *WasmGasMeter) GasConsumed() uint64 {
	sdkGas := gm.consumed / GasConversionFactor
	gm.logger.Info("📊 Gas consumption check",
		"wasm_gas_consumed", gm.consumed,
		"sdk_gas_consumed", sdkGas)
	return gm.consumed
}

func (gm *WasmGasMeter) GasConsumedToLimit() uint64 {
	if gm.consumed > gm.limit {
		gm.logger.Info("📊 Gas consumption exceeds limit",
			"wasm_gas_consumed", gm.consumed,
			"wasm_gas_limit", gm.limit,
			"sdk_gas_consumed", gm.consumed/GasConversionFactor,
			"sdk_gas_limit", gm.limit/GasConversionFactor)
		return gm.limit
	}
	gm.logger.Info("📊 Gas consumption within limit",
		"wasm_gas_consumed", gm.consumed,
		"wasm_gas_limit", gm.limit,
		"sdk_gas_consumed", gm.consumed/GasConversionFactor,
		"sdk_gas_limit", gm.limit/GasConversionFactor)
	return gm.consumed
}

func (gm *WasmGasMeter) Limit() uint64 {
	gm.logger.Info("📊 Gas limit check",
		"wasm_gas_limit", gm.limit,
		"sdk_gas_limit", gm.limit/GasConversionFactor)
	return gm.limit
}

func (gm *WasmGasMeter) ConsumeGas(amount uint64, descriptor string) error {
	gm.logger.Info("⚡️ Gas consumption request",
		"operation", descriptor,
		"wasm_gas_requested", amount,
		"sdk_gas_requested", amount/GasConversionFactor,
		"wasm_gas_remaining", gm.limit-gm.consumed,
		"sdk_gas_remaining", (gm.limit-gm.consumed)/GasConversionFactor)

	if gm.consumed+amount > gm.limit {
		gm.logger.Error("❌ Out of gas",
			"operation", descriptor,
			"wasm_gas_consumed", gm.consumed,
			"wasm_gas_limit", gm.limit,
			"wasm_gas_requested", amount,
			"sdk_gas_consumed", gm.consumed/GasConversionFactor,
			"sdk_gas_limit", gm.limit/GasConversionFactor,
			"sdk_gas_requested", amount/GasConversionFactor)
		return fmt.Errorf("out of gas: consumed %d, limit %d, requested %d (sdk units: %d/%d)",
			gm.consumed,
			gm.limit,
			amount,
			gm.consumed/GasConversionFactor,
			gm.limit/GasConversionFactor)
	}

	gm.consumed += amount
	gm.logger.Info("✅ Gas consumed successfully",
		"operation", descriptor,
		"wasm_gas_consumed", amount,
		"wasm_gas_total", gm.consumed,
		"wasm_gas_remaining", gm.limit-gm.consumed,
		"sdk_gas_consumed", amount/GasConversionFactor,
		"sdk_gas_total", gm.consumed/GasConversionFactor,
		"sdk_gas_remaining", (gm.limit-gm.consumed)/GasConversionFactor)
	return nil
}

// IsolatedContractStore creates an isolated store for contract state
func (k Keeper) IsolatedContractStore(ctx context.Context, contractAddr string) corestore.KVStore {
	// Create an isolated store for this specific contract
	// This ensures contracts cannot access each other's state
	//
	// TODO: In production, this would implement proper store isolation
	// For now, we return the base store service as a placeholder
	// Each contract's state will be isolated through proper key prefixing
	// when we implement actual contract execution in future steps
	return k.storeService.OpenKVStore(ctx)
}

// ValidateContractCodeSize ensures contract code is within size limits
func (k Keeper) ValidateContractCodeSize(code []byte) error {
	return k.gasConfig.ValidateContractSize(code)
}

// ValidateMessageSize ensures messages are within size limits
func (k Keeper) ValidateMessageSize(msg []byte, msgType string) error {
	return k.gasConfig.ValidateMessageSize(msg, msgType)
}

// GetWasmVMCapabilities returns the security-configured capabilities for WasmVM
func (k Keeper) GetWasmVMCapabilities() string {
	// Return the capabilities that were configured during VM initialization
	// These are the same capabilities we set in keeper.go
	return "iterator,staking,stargate,cosmwasm_1_1,cosmwasm_1_2"
}

// LogContractExecution logs contract execution for monitoring and debugging
func (k Keeper) LogContractExecution(ctx context.Context, contractAddr, operation string, gasUsed uint64, success bool) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)

	eventType := "contract_execution"
	if !success {
		eventType = "contract_execution_failed"
	}

	sdkCtx.EventManager().EmitEvent(
		sdk.NewEvent(
			eventType,
			sdk.NewAttribute("contract_address", contractAddr),
			sdk.NewAttribute("operation", operation),
			sdk.NewAttribute("gas_used", fmt.Sprintf("%d", gasUsed)),
			sdk.NewAttribute("success", fmt.Sprintf("%t", success)),
		),
	)
}

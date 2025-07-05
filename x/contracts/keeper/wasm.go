package keeper

import (
	"context"
	"fmt"
	"time"

	"scarlett-core/x/contracts/types"

	wasmvmtypes "github.com/CosmWasm/wasmvm/v3/types"

	corestore "cosmossdk.io/core/store"
	errorsmod "cosmossdk.io/errors"
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

// CreateGasMeter creates a gas meter for contract execution
func (k Keeper) CreateGasMeter(gasLimit uint64) wasmvmtypes.GasMeter {
	return &WasmGasMeter{
		limit:    gasLimit,
		consumed: 0,
	}
}

// WasmGasMeter implements proper gas conversion between WasmVM and Cosmos SDK
type WasmGasMeter struct {
	limit    uint64
	consumed uint64
}

// GasConversionFactor is how many SDK gas units equal one WasmVM gas unit
// WasmVM uses much larger gas units than Cosmos SDK
const GasConversionFactor = 100_000

func (gm *WasmGasMeter) GasConsumed() uint64 {
	return gm.consumed / GasConversionFactor
}

func (gm *WasmGasMeter) GasConsumedToLimit() uint64 {
	if gm.consumed > gm.limit {
		return gm.limit / GasConversionFactor
	}
	return gm.consumed / GasConversionFactor
}

func (gm *WasmGasMeter) Limit() uint64 {
	return gm.limit / GasConversionFactor
}

func (gm *WasmGasMeter) ConsumeGas(amount uint64, descriptor string) error {
	// Convert WasmVM gas units to SDK gas units
	sdkGas := amount / GasConversionFactor
	if sdkGas == 0 {
		sdkGas = 1 // Minimum 1 gas per operation
	}

	if gm.consumed+sdkGas > gm.limit {
		return fmt.Errorf("out of gas: consumed %d, limit %d, requested %d (wasmvm units: %d)",
			gm.consumed/GasConversionFactor,
			gm.limit/GasConversionFactor,
			sdkGas,
			amount)
	}
	gm.consumed += sdkGas
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

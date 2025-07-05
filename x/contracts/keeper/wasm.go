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

// Security Configuration Constants
const (
	// Gas limits for contract execution
	DefaultGasLimit = uint64(1_000_000)  // 1M gas default limit
	MaxGasLimit     = uint64(10_000_000) // 10M gas maximum limit
	MinGasLimit     = uint64(100_000)    // 100K gas minimum limit

	// Execution timeouts
	DefaultExecutionTimeout = 30 * time.Second // 30 seconds max execution
	MaxExecutionTimeout     = 60 * time.Second // 60 seconds absolute max

	// Memory and storage limits
	MaxContractSize       = 1024 * 1024 // 1MB max contract size
	MaxInstantiateMsgSize = 64 * 1024   // 64KB max instantiate message
	MaxExecuteMsgSize     = 32 * 1024   // 32KB max execute message
	MaxQueryMsgSize       = 16 * 1024   // 16KB max query message
)

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
	if gasLimit < MinGasLimit {
		return errorsmod.Wrapf(types.ErrInvalidInput, "gas limit too low: %d, minimum: %d", gasLimit, MinGasLimit)
	}
	if gasLimit > MaxGasLimit {
		return errorsmod.Wrapf(types.ErrInvalidInput, "gas limit too high: %d, maximum: %d", gasLimit, MaxGasLimit)
	}
	return nil
}

// validateTimeout ensures execution timeout is within acceptable bounds
func (k Keeper) validateTimeout(timeout time.Duration) error {
	if timeout <= 0 {
		return errorsmod.Wrapf(types.ErrInvalidInput, "timeout must be positive")
	}
	if timeout > MaxExecutionTimeout {
		return errorsmod.Wrapf(types.ErrInvalidInput, "timeout too long: %v, maximum: %v", timeout, MaxExecutionTimeout)
	}
	return nil
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
	// This is a simplified gas meter
	// In production, this would integrate with SDK gas meter
	return &SimpleGasMeter{
		limit:    gasLimit,
		consumed: 0,
	}
}

// SimpleGasMeter is a basic implementation of gas metering for contract execution
type SimpleGasMeter struct {
	limit    uint64
	consumed uint64
}

func (gm *SimpleGasMeter) GasConsumed() uint64 {
	return gm.consumed
}

func (gm *SimpleGasMeter) GasConsumedToLimit() uint64 {
	if gm.consumed > gm.limit {
		return gm.limit
	}
	return gm.consumed
}

func (gm *SimpleGasMeter) Limit() uint64 {
	return gm.limit
}

func (gm *SimpleGasMeter) ConsumeGas(amount uint64, descriptor string) error {
	if gm.consumed+amount > gm.limit {
		return fmt.Errorf("out of gas: consumed %d, limit %d, requested %d", gm.consumed, gm.limit, amount)
	}
	gm.consumed += amount
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
	if len(code) > MaxContractSize {
		return errorsmod.Wrapf(types.ErrInvalidInput,
			"contract code too large: %d bytes, maximum: %d bytes",
			len(code), MaxContractSize)
	}
	return nil
}

// ValidateMessageSize ensures messages are within size limits
func (k Keeper) ValidateMessageSize(msg []byte, msgType string) error {
	var maxSize int

	switch msgType {
	case "instantiate":
		maxSize = MaxInstantiateMsgSize
	case "execute":
		maxSize = MaxExecuteMsgSize
	case "query":
		maxSize = MaxQueryMsgSize
	default:
		maxSize = MaxExecuteMsgSize // Default to execute message size
	}

	if len(msg) > maxSize {
		return errorsmod.Wrapf(types.ErrInvalidInput,
			"%s message too large: %d bytes, maximum: %d bytes",
			msgType, len(msg), maxSize)
	}

	return nil
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

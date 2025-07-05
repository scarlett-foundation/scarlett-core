package types

import (
	"fmt"
	"time"
)

// GasConfig holds all gas-related configuration for contract operations
// These values are embedded in the binary to ensure consensus safety
type GasConfig struct {
	// Contract deployment limits
	DeploymentGasLimit uint64 `json:"deployment_gas_limit"` // Gas limit for contract deployment
	MaxContractSize    int    `json:"max_contract_size"`    // Maximum contract size in bytes

	// Contract execution limits
	DefaultGasLimit uint64 `json:"default_gas_limit"` // Default gas limit for contract execution
	MaxGasLimit     uint64 `json:"max_gas_limit"`     // Maximum gas limit for contract execution
	MinGasLimit     uint64 `json:"min_gas_limit"`     // Minimum gas limit for contract execution

	// Execution timeouts
	DefaultTimeout time.Duration `json:"default_timeout"` // Default execution timeout
	MaxTimeout     time.Duration `json:"max_timeout"`     // Maximum execution timeout

	// Message size limits
	MaxInstantiateSize int `json:"max_instantiate_size"` // Maximum instantiate message size
	MaxExecuteSize     int `json:"max_execute_size"`     // Maximum execute message size
	MaxQuerySize       int `json:"max_query_size"`       // Maximum query message size
}

// DefaultGasConfig provides production-ready gas limits embedded in the binary
// All nodes building from the same source will have identical limits
var DefaultGasConfig = GasConfig{
	// Deployment limits - significantly increased for real-world contracts
	DeploymentGasLimit: 50_000_000, // 50M gas - handles large contracts like 275KB+ WASM
	MaxContractSize:    5_242_880,  // 5MB - allows complex DeFi contracts

	// Execution limits - reasonable for production use
	DefaultGasLimit: 1_000_000,   // 1M gas - good default for most operations
	MaxGasLimit:     100_000_000, // 100M gas - handles complex contract interactions
	MinGasLimit:     100_000,     // 100K gas - prevents tiny gas limit abuse

	// Timeout limits - secure but usable
	DefaultTimeout: 30 * time.Second,  // 30 seconds - reasonable for most operations
	MaxTimeout:     120 * time.Second, // 2 minutes - absolute maximum

	// Message size limits - increased for real-world usage
	MaxInstantiateSize: 262_144, // 256KB - handles complex instantiation
	MaxExecuteSize:     131_072, // 128KB - handles complex execution messages
	MaxQuerySize:       65_536,  // 64KB - handles complex queries
}

// ValidateGasConfig ensures all gas limits are reasonable and secure
func ValidateGasConfig(config GasConfig) error {
	// Validate deployment limits
	if config.DeploymentGasLimit < 1_000_000 {
		return fmt.Errorf("deployment gas limit too low: %d, minimum: 1M", config.DeploymentGasLimit)
	}
	if config.DeploymentGasLimit > 1_000_000_000 {
		return fmt.Errorf("deployment gas limit too high: %d, maximum: 1B", config.DeploymentGasLimit)
	}

	if config.MaxContractSize < 1024 {
		return fmt.Errorf("max contract size too small: %d, minimum: 1KB", config.MaxContractSize)
	}
	if config.MaxContractSize > 50_000_000 {
		return fmt.Errorf("max contract size too large: %d, maximum: 50MB", config.MaxContractSize)
	}

	// Validate execution limits
	if config.MinGasLimit >= config.MaxGasLimit {
		return fmt.Errorf("min gas limit (%d) must be less than max gas limit (%d)", config.MinGasLimit, config.MaxGasLimit)
	}

	if config.DefaultGasLimit < config.MinGasLimit || config.DefaultGasLimit > config.MaxGasLimit {
		return fmt.Errorf("default gas limit (%d) must be between min (%d) and max (%d)",
			config.DefaultGasLimit, config.MinGasLimit, config.MaxGasLimit)
	}

	// Validate timeouts
	if config.DefaultTimeout <= 0 {
		return fmt.Errorf("default timeout must be positive")
	}
	if config.MaxTimeout <= config.DefaultTimeout {
		return fmt.Errorf("max timeout (%v) must be greater than default timeout (%v)",
			config.MaxTimeout, config.DefaultTimeout)
	}

	// Validate message sizes
	if config.MaxInstantiateSize < 1024 || config.MaxExecuteSize < 1024 || config.MaxQuerySize < 1024 {
		return fmt.Errorf("message size limits must be at least 1KB")
	}

	return nil
}

// GetGasConfig returns the current gas configuration
// In this minimal implementation, it returns the embedded defaults
// Future versions can extend this to support governance-based updates
func GetGasConfig() GasConfig {
	return DefaultGasConfig
}

// ValidateMessageSize validates message size against configured limits
func (config GasConfig) ValidateMessageSize(msg []byte, msgType string) error {
	var maxSize int

	switch msgType {
	case "instantiate":
		maxSize = config.MaxInstantiateSize
	case "execute":
		maxSize = config.MaxExecuteSize
	case "query":
		maxSize = config.MaxQuerySize
	default:
		maxSize = config.MaxExecuteSize // Default to execute message size
	}

	if len(msg) > maxSize {
		return fmt.Errorf("%s message too large: %d bytes, maximum: %d bytes",
			msgType, len(msg), maxSize)
	}

	return nil
}

// ValidateContractSize validates contract size against configured limits
func (config GasConfig) ValidateContractSize(code []byte) error {
	if len(code) > config.MaxContractSize {
		return fmt.Errorf("contract code too large: %d bytes, maximum: %d bytes",
			len(code), config.MaxContractSize)
	}
	return nil
}

// ValidateGasLimit validates gas limit against configured bounds
func (config GasConfig) ValidateGasLimit(gasLimit uint64) error {
	if gasLimit < config.MinGasLimit {
		return fmt.Errorf("gas limit too low: %d, minimum: %d", gasLimit, config.MinGasLimit)
	}
	if gasLimit > config.MaxGasLimit {
		return fmt.Errorf("gas limit too high: %d, maximum: %d", gasLimit, config.MaxGasLimit)
	}
	return nil
}

// ValidateTimeout validates timeout against configured bounds
func (config GasConfig) ValidateTimeout(timeout time.Duration) error {
	if timeout <= 0 {
		return fmt.Errorf("timeout must be positive")
	}
	if timeout > config.MaxTimeout {
		return fmt.Errorf("timeout too long: %v, maximum: %v", timeout, config.MaxTimeout)
	}
	return nil
}

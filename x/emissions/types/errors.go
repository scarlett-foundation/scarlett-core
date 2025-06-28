package types

// DONTCOVER

import (
	"cosmossdk.io/errors"
)

// x/emissions module sentinel errors
var (
	ErrInvalidSigner               = errors.Register(ModuleName, 1100, "expected gov account as only signer for proposal message")
	ErrInvalidDestination          = errors.Register(ModuleName, 1101, "invalid emission destination")
	ErrInvalidWeight               = errors.Register(ModuleName, 1102, "invalid emission weight")
	ErrExceedsMaxDestinations      = errors.Register(ModuleName, 1103, "exceeds maximum destinations")
	ErrBelowMinValidatorReward     = errors.Register(ModuleName, 1104, "below minimum validator reward")
	ErrExceedsMaxSingleDestination = errors.Register(ModuleName, 1105, "exceeds maximum single destination")
	ErrEmissionsStopped            = errors.Register(ModuleName, 1106, "emissions are emergency stopped")
	ErrUnauthorized                = errors.Register(ModuleName, 1107, "unauthorized emission parameter change")
	ErrEmissionParamsNotFound      = errors.Register(ModuleName, 1108, "emission parameters not found")

	// Additional safety and validation errors
	ErrNoDestinations         = errors.Register(ModuleName, 1109, "no emission destinations configured")
	ErrDuplicateDestination   = errors.Register(ModuleName, 1110, "duplicate emission destination")
	ErrInvalidTotalWeight     = errors.Register(ModuleName, 1111, "total emission weight must equal 1.0")
	ErrBelowMinWeight         = errors.Register(ModuleName, 1112, "emission weight below minimum bound")
	ErrExceedsMaxWeight       = errors.Register(ModuleName, 1113, "emission weight exceeds maximum bound")
	ErrEmergencyActive        = errors.Register(ModuleName, 1114, "emergency controls are active")
	ErrInvalidEmergencyParams = errors.Register(ModuleName, 1115, "invalid emergency parameters")

	// Module registry errors
	ErrModuleAlreadyRegistered = errors.Register(ModuleName, 1116, "module is already registered")
	ErrModuleNotRegistered     = errors.Register(ModuleName, 1117, "module is not registered")
	ErrInvalidModuleName       = errors.Register(ModuleName, 1118, "invalid module name")
)

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
)

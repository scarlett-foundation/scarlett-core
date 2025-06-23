package types

// DONTCOVER

import (
	"cosmossdk.io/errors"
)

// x/scarlettcore module sentinel errors
var (
	ErrInvalidSigner = errors.Register(ModuleName, 1100, "expected gov account as only signer for proposal message")

	// Burn token errors
	ErrInsufficientBalance = errors.Register(ModuleName, 1101, "insufficient balance to burn")
	ErrInvalidAmount       = errors.Register(ModuleName, 1102, "invalid burn amount")
	ErrInvalidDenom        = errors.Register(ModuleName, 1103, "invalid denomination")
)

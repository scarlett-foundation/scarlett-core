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
	ErrInvalidAddress      = errors.Register(ModuleName, 1104, "invalid address")

	// Genesis burn errors
	ErrNotGenesisAddress      = errors.Register(ModuleName, 2001, "caller is not the configured genesis address")
	ErrAlreadyExecuted        = errors.Register(ModuleName, 2002, "genesis burn already executed")
	ErrInsufficientValidators = errors.Register(ModuleName, 2003, "insufficient active validators for safe execution")
	ErrNoStakeToBurn          = errors.Register(ModuleName, 2004, "genesis address has no stake to burn")
)

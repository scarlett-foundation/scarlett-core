package types

// DONTCOVER

import (
	"cosmossdk.io/errors"
)

// x/proofofdegen module sentinel errors
var (
	ErrInvalidSigner  = errors.Register(ModuleName, 1100, "expected gov account as only signer for proposal message")
	ErrUnauthorized   = errors.Register(ModuleName, 1101, "unauthorized")
	ErrNotEligible    = errors.Register(ModuleName, 1102, "address not eligible for airdrop")
	ErrAlreadyClaimed = errors.Register(ModuleName, 1103, "airdrop already claimed")
	ErrNoTokensTolaim = errors.Register(ModuleName, 1104, "no tokens available to claim")
	ErrInvalidProof   = errors.Register(ModuleName, 1105, "invalid merkle proof")
)

package types

// DONTCOVER

import (
	"cosmossdk.io/errors"
)

// x/contracts module sentinel errors
var (
	ErrInvalidSigner      = errors.Register(ModuleName, 1100, "expected gov account as only signer for proposal message")
	ErrWasmNotInitialized = errors.Register(ModuleName, 1101, "wasm keeper not initialized")
)

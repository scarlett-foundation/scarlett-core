package types

import "cosmossdk.io/collections"

const (
	// ModuleName defines the module name
	ModuleName = "contracts"

	// StoreKey defines the primary module store key
	StoreKey = ModuleName

	// GovModuleName duplicates the gov module's name to avoid a dependency with x/gov.
	// It should be synced with the gov module's name if it is ever changed.
	// See: https://github.com/cosmos/cosmos-sdk/blob/v0.52.0-beta.2/x/gov/types/keys.go#L9
	GovModuleName = "gov"
)

// Collection keys
var (
	// ParamsKey is the prefix to retrieve all Params
	ParamsKey = collections.NewPrefix("p_contracts")

	// ContractCodeKey is the prefix for storing contract bytecode
	ContractCodeKey = collections.NewPrefix("contract_code")

	// ContractInfoKey is the prefix for storing contract metadata
	ContractInfoKey = collections.NewPrefix("contract_info")

	// ContractSeqKey is the prefix for contract ID sequence
	ContractSeqKey = collections.NewPrefix("contract_seq")
)

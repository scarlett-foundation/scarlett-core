package types

import "cosmossdk.io/collections"

const (
	// ModuleName defines the module name
	ModuleName = "emissions"

	// StoreKey defines the primary module store key
	StoreKey = ModuleName

	// GovModuleName duplicates the gov module's name to avoid a dependency with x/gov.
	// It should be synced with the gov module's name if it is ever changed.
	// See: https://github.com/cosmos/cosmos-sdk/blob/v0.52.0-beta.2/x/gov/types/keys.go#L9
	GovModuleName = "gov"
)

// Collection keys for state management
var (
	// ParamsKey is the prefix to retrieve all Params
	ParamsKey = collections.NewPrefix("p_emissions")

	// EmissionParamsKey is the prefix for emission parameters
	EmissionParamsKey = collections.NewPrefix("emission_params")

	// EmissionHistoryPrefix is the prefix for emission parameter history
	EmissionHistoryPrefix = collections.NewPrefix("emission_history")

	// DestinationMetricsPrefix is the prefix for destination metrics
	DestinationMetricsPrefix = collections.NewPrefix("destination_metrics")
)

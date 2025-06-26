package types

import (
	"fmt"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// EmissionDestination represents a single destination for minted tokens
type EmissionDestination struct {
	ModuleName  string         `json:"module_name"`
	Weight      math.LegacyDec `json:"weight"` // Percentage weight (0.0 to 1.0)
	Description string         `json:"description"`
	Enabled     bool           `json:"enabled"`    // Can be disabled without removal
	MinWeight   math.LegacyDec `json:"min_weight"` // Governance safety bounds
	MaxWeight   math.LegacyDec `json:"max_weight"` // Governance safety bounds
}

// EmissionParams holds the governance-controlled emission configuration
type EmissionParams struct {
	Destinations []EmissionDestination `json:"destinations"`
	Enabled      bool                  `json:"enabled"`
	UpdatedBy    string                `json:"updated_by"` // Last governance proposal
	UpdatedAt    int64                 `json:"updated_at"` // Block height of last update
}

// DestinationStats tracks metrics for each emission destination
type DestinationStats struct {
	ModuleName    string    `json:"module_name"`
	TotalReceived sdk.Coins `json:"total_received"`
	LastAmount    sdk.Coins `json:"last_amount"`
	LastUpdate    int64     `json:"last_update"`
	ActiveSince   int64     `json:"active_since"`
}

// EmissionValidationRules defines safety bounds for governance proposals
type EmissionValidationRules struct {
	MinValidatorReward   math.LegacyDec // e.g., 0.30 (30% minimum to validators)
	MaxSingleDestination math.LegacyDec // e.g., 0.70 (70% maximum to any destination)
	MaxDestinations      uint32         // e.g., 10 maximum destinations
	MinProposalDeposit   sdk.Coins      // Minimum deposit for emission proposals
}

// Constants for module names
const (
	InferenceRewardsModuleName = "inferencerewards"
	FeeCollectorModuleName     = "fee_collector"
)

// NewParams creates a new Params instance.
func NewParams() Params {
	return Params{}
}

// DefaultParams returns a default set of parameters.
func DefaultParams() Params {
	return NewParams()
}

// DefaultEmissionParams returns the default 50/50 split configuration
func DefaultEmissionParams() EmissionParams {
	return EmissionParams{
		Enabled: true,
		Destinations: []EmissionDestination{
			{
				ModuleName:  FeeCollectorModuleName,
				Weight:      math.LegacyNewDecWithPrec(50, 2), // 0.50 (50%)
				Description: "Traditional staking rewards distributed to validators and delegators",
				Enabled:     true,
				MinWeight:   math.LegacyNewDecWithPrec(30, 2), // 30% minimum
				MaxWeight:   math.LegacyNewDecWithPrec(70, 2), // 70% maximum
			},
			{
				ModuleName:  InferenceRewardsModuleName,
				Weight:      math.LegacyNewDecWithPrec(50, 2), // 0.50 (50%)
				Description: "AI/ML computing rewards for miners and inference providers",
				Enabled:     true,
				MinWeight:   math.LegacyNewDecWithPrec(30, 2), // 30% minimum
				MaxWeight:   math.LegacyNewDecWithPrec(70, 2), // 70% maximum
			},
		},
		UpdatedBy: "genesis",
		UpdatedAt: 0,
	}
}

// DefaultValidationRules returns the default validation rules for emission proposals
func DefaultValidationRules() EmissionValidationRules {
	return EmissionValidationRules{
		MinValidatorReward:   math.LegacyNewDecWithPrec(30, 2),                         // 30% minimum to validators
		MaxSingleDestination: math.LegacyNewDecWithPrec(70, 2),                         // 70% maximum to any destination
		MaxDestinations:      10,                                                       // Maximum 10 destinations
		MinProposalDeposit:   sdk.NewCoins(sdk.NewCoin("sclt", math.NewInt(10000000))), // 10M sclt
	}
}

// Validate validates the set of params.
func (p Params) Validate() error {
	return nil
}

// Validate validates the emission parameters
func (ep EmissionParams) Validate() error {
	if !ep.Enabled {
		return nil // Skip validation if disabled
	}

	if len(ep.Destinations) == 0 {
		return fmt.Errorf("no emission destinations configured")
	}

	totalWeight := math.LegacyZeroDec()
	moduleNames := make(map[string]bool)
	rules := DefaultValidationRules()

	// Check maximum destinations
	if len(ep.Destinations) > int(rules.MaxDestinations) {
		return fmt.Errorf("too many destinations: %d, maximum allowed: %d", len(ep.Destinations), rules.MaxDestinations)
	}

	for i, dest := range ep.Destinations {
		// Check for duplicate module names
		if moduleNames[dest.ModuleName] {
			return fmt.Errorf("duplicate module name: %s", dest.ModuleName)
		}
		moduleNames[dest.ModuleName] = true

		// Validate weight
		if dest.Weight.IsNegative() {
			return fmt.Errorf("destination %d: weight cannot be negative", i)
		}
		if dest.Weight.GT(math.LegacyOneDec()) {
			return fmt.Errorf("destination %d: weight cannot exceed 1.0", i)
		}

		// Check safety bounds
		if dest.Weight.GT(rules.MaxSingleDestination) {
			return fmt.Errorf("destination %d: weight %s exceeds maximum allowed %s", i, dest.Weight.String(), rules.MaxSingleDestination.String())
		}

		// Check minimum validator reward (fee_collector)
		if dest.ModuleName == FeeCollectorModuleName && dest.Weight.LT(rules.MinValidatorReward) {
			return fmt.Errorf("validator reward weight %s below minimum required %s", dest.Weight.String(), rules.MinValidatorReward.String())
		}

		// Validate min/max weight bounds for destination
		if dest.MinWeight.IsNegative() || dest.MaxWeight.IsNegative() {
			return fmt.Errorf("destination %d: min/max weights cannot be negative", i)
		}
		if dest.MinWeight.GT(dest.MaxWeight) {
			return fmt.Errorf("destination %d: min weight cannot exceed max weight", i)
		}
		if dest.Weight.LT(dest.MinWeight) || dest.Weight.GT(dest.MaxWeight) {
			return fmt.Errorf("destination %d: weight %s outside bounds [%s, %s]", i, dest.Weight.String(), dest.MinWeight.String(), dest.MaxWeight.String())
		}

		totalWeight = totalWeight.Add(dest.Weight)
	}

	// Check total weight equals 1.0 (100%)
	if !totalWeight.Equal(math.LegacyOneDec()) {
		return fmt.Errorf("total weight must equal 1.0, got %s", totalWeight.String())
	}

	return nil
}

// ValidateDestination validates a single emission destination
func (dest EmissionDestination) Validate() error {
	if dest.ModuleName == "" {
		return fmt.Errorf("module name cannot be empty")
	}

	if dest.Weight.IsNegative() {
		return fmt.Errorf("weight cannot be negative")
	}

	if dest.Weight.GT(math.LegacyOneDec()) {
		return fmt.Errorf("weight cannot exceed 1.0")
	}

	if dest.MinWeight.IsNegative() || dest.MaxWeight.IsNegative() {
		return fmt.Errorf("min/max weights cannot be negative")
	}

	if dest.MinWeight.GT(dest.MaxWeight) {
		return fmt.Errorf("min weight cannot exceed max weight")
	}

	return nil
}

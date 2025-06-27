package emissions

import (
	"fmt"

	"cosmossdk.io/math"
)

// EmissionDestination represents a single destination for minted tokens
type EmissionDestination struct {
	ModuleName  string         `json:"module_name"`
	Weight      math.LegacyDec `json:"weight"` // Percentage weight (0.0 to 1.0)
	Description string         `json:"description"`
}

// EmissionsConfig holds the configuration for custom emission splitting
type EmissionsConfig struct {
	Destinations []EmissionDestination `json:"destinations"`
	Enabled      bool                  `json:"enabled"`
}

// Constants for module names
const (
	InferenceRewardsModuleName = "inferencerewards"
	FeeCollectorModuleName     = "fee_collector"
)

// DefaultEmissionsConfig returns the default 50/50 split configuration
func DefaultEmissionsConfig() EmissionsConfig {
	return EmissionsConfig{
		Enabled: true,
		Destinations: []EmissionDestination{
			{
				ModuleName:  FeeCollectorModuleName,
				Weight:      math.LegacyNewDecWithPrec(50, 2), // 0.50 (50%)
				Description: "Traditional staking rewards distributed to validators and delegators",
			},
			{
				ModuleName:  InferenceRewardsModuleName,
				Weight:      math.LegacyNewDecWithPrec(50, 2), // 0.50 (50%)
				Description: "AI/ML computing rewards for miners and inference providers",
			},
		},
	}
}

// Validate ensures the configuration is valid
func (c EmissionsConfig) Validate() error {
	if !c.Enabled {
		return nil
	}

	if len(c.Destinations) == 0 {
		return fmt.Errorf("no emission destinations configured")
	}

	totalWeight := math.LegacyZeroDec()
	moduleNames := make(map[string]bool)

	for i, dest := range c.Destinations {
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

		totalWeight = totalWeight.Add(dest.Weight)
	}

	// Check total weight equals 1.0 (100%)
	if !totalWeight.Equal(math.LegacyOneDec()) {
		return fmt.Errorf("total weight must equal 1.0, got %s", totalWeight.String())
	}

	return nil
}

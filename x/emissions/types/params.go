package types

import (
	"encoding/json"
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
	Destinations    []EmissionDestination `json:"destinations"`
	Enabled         bool                  `json:"enabled"`
	UpdatedBy       string                `json:"updated_by"`       // Last governance proposal ID
	UpdatedAt       int64                 `json:"updated_at"`       // Block height of last update
	EmergencyStop   bool                  `json:"emergency_stop"`   // Emergency halt flag
	FallbackEnabled bool                  `json:"fallback_enabled"` // Use fallback configuration
}

// DestinationStats tracks metrics for each emission destination
type DestinationStats struct {
	ModuleName    string    `json:"module_name"`
	TotalReceived sdk.Coins `json:"total_received"`
	LastAmount    sdk.Coins `json:"last_amount"`
	LastUpdate    int64     `json:"last_update"`
	ActiveSince   int64     `json:"active_since"`
	Distribution  int64     `json:"distribution_count"` // Number of distributions
}

// EmissionValidationRules defines economic safety bounds for emission parameters
type EmissionValidationRules struct {
	MinValidatorReward   math.LegacyDec `json:"min_validator_reward"`   // e.g., 0.30 (30% minimum to validators)
	MaxSingleDestination math.LegacyDec `json:"max_single_destination"` // e.g., 0.70 (70% maximum to any destination)
	MaxDestinations      uint32         `json:"max_destinations"`       // e.g., 10 maximum destinations
	MinProposalDeposit   sdk.Coins      `json:"min_proposal_deposit"`   // Minimum deposit for emission proposals
	RequireQuorum        bool           `json:"require_quorum"`         // Require quorum for emission changes
}

// EmergencyControls defines emergency control mechanisms
type EmergencyControls struct {
	EmergencyStop        bool   `json:"emergency_stop"`         // Halt all emissions
	FallbackToDefault    bool   `json:"fallback_to_default"`    // Revert to 50/50 split
	EmergencyAuthority   string `json:"emergency_authority"`    // Emergency halt authority address
	EmergencyReason      string `json:"emergency_reason"`       // Reason for emergency action
	EmergencyActivatedAt int64  `json:"emergency_activated_at"` // Block height when emergency activated
}

// Constants for module names
const (
	InferenceRewardsModuleName = "inferencerewards"
	FeeCollectorModuleName     = "fee_collector"
)

// DefaultParams returns a default set of parameters
func DefaultParams() Params {
	return NewParams()
}

// NewParams creates a new Params instance
func NewParams() Params {
	return Params{}
}

// Validate validates the set of params
func (p Params) Validate() error {
	return nil
}

// DefaultEmissionParams returns the default emission parameters with safety mechanisms
func DefaultEmissionParams() EmissionParams {
	return EmissionParams{
		Destinations: []EmissionDestination{
			{
				ModuleName:  FeeCollectorModuleName,
				Weight:      math.LegacyNewDecWithPrec(50, 2), // 0.50 = 50%
				Description: "Validator and delegator rewards",
				Enabled:     true,
				MinWeight:   math.LegacyNewDecWithPrec(30, 2), // 30% minimum
				MaxWeight:   math.LegacyNewDecWithPrec(70, 2), // 70% maximum
			},
			{
				ModuleName:  InferenceRewardsModuleName,
				Weight:      math.LegacyNewDecWithPrec(50, 2), // 0.50 = 50%
				Description: "AI inference provider rewards",
				Enabled:     true,
				MinWeight:   math.LegacyNewDecWithPrec(30, 2), // 30% minimum
				MaxWeight:   math.LegacyNewDecWithPrec(70, 2), // 70% maximum
			},
		},
		Enabled:         true,
		UpdatedBy:       "genesis",
		UpdatedAt:       0,
		EmergencyStop:   false,
		FallbackEnabled: false,
	}
}

// DefaultValidationRules returns the default validation rules with economic safety bounds
func DefaultValidationRules() EmissionValidationRules {
	return EmissionValidationRules{
		MinValidatorReward:   math.LegacyNewDecWithPrec(30, 2),                         // 30% minimum to validators
		MaxSingleDestination: math.LegacyNewDecWithPrec(70, 2),                         // 70% maximum to any single destination
		MaxDestinations:      10,                                                       // Maximum 10 destinations
		MinProposalDeposit:   sdk.NewCoins(sdk.NewCoin("sclt", math.NewInt(10000000))), // 10 SCLT minimum
		RequireQuorum:        true,                                                     // Require governance quorum
	}
}

// ValidateEmissionParams performs comprehensive validation of emission parameters
func ValidateEmissionParams(params EmissionParams) error {
	rules := DefaultValidationRules()

	// 1. Check if emergency stop is active
	if params.EmergencyStop {
		// Emergency stop overrides all other validations
		return nil
	}

	// 2. Validate destinations count
	if len(params.Destinations) == 0 {
		return ErrNoDestinations
	}
	if len(params.Destinations) > int(rules.MaxDestinations) {
		return ErrExceedsMaxDestinations
	}

	// 3. Check for duplicate destinations
	seen := make(map[string]bool)
	for _, dest := range params.Destinations {
		if seen[dest.ModuleName] {
			return ErrDuplicateDestination
		}
		seen[dest.ModuleName] = true
	}

	// 4. Validate each destination
	totalWeight := math.LegacyZeroDec()
	validatorWeight := math.LegacyZeroDec()

	for _, dest := range params.Destinations {
		if !dest.Enabled {
			continue // Skip disabled destinations
		}

		// Validate weight bounds
		if dest.Weight.IsNegative() || dest.Weight.GT(math.LegacyOneDec()) {
			return ErrInvalidWeight
		}

		// Check individual destination limits
		if dest.Weight.GT(rules.MaxSingleDestination) {
			return ErrExceedsMaxSingleDestination
		}

		// Check minimum weight constraints
		if !dest.MinWeight.IsZero() && dest.Weight.LT(dest.MinWeight) {
			return ErrBelowMinWeight
		}

		// Check maximum weight constraints
		if !dest.MaxWeight.IsZero() && dest.Weight.GT(dest.MaxWeight) {
			return ErrExceedsMaxWeight
		}

		// Track validator rewards (fee_collector)
		if dest.ModuleName == FeeCollectorModuleName {
			validatorWeight = dest.Weight
		}

		totalWeight = totalWeight.Add(dest.Weight)
	}

	// 5. Validate total weight equals 100%
	if !totalWeight.Equal(math.LegacyOneDec()) {
		return ErrInvalidTotalWeight
	}

	// 6. Enforce minimum validator reward safety bound
	if validatorWeight.LT(rules.MinValidatorReward) {
		return ErrBelowMinValidatorReward
	}

	// 7. Validate module names (basic validation)
	for _, dest := range params.Destinations {
		if dest.ModuleName == "" {
			return ErrInvalidDestination
		}
	}

	return nil
}

// ValidateDestination validates a single emission destination
func ValidateDestination(dest EmissionDestination) error {
	if dest.ModuleName == "" {
		return ErrInvalidDestination
	}

	if dest.Weight.IsNegative() || dest.Weight.GT(math.LegacyOneDec()) {
		return ErrInvalidWeight
	}

	if !dest.MinWeight.IsZero() && dest.Weight.LT(dest.MinWeight) {
		return ErrBelowMinWeight
	}

	if !dest.MaxWeight.IsZero() && dest.Weight.GT(dest.MaxWeight) {
		return ErrExceedsMaxWeight
	}

	return nil
}

// IsEmergencyActive checks if emergency controls are active
func (p EmissionParams) IsEmergencyActive() bool {
	return p.EmergencyStop
}

// IsValidatorRewardSufficient checks if validator rewards meet minimum requirements
func (p EmissionParams) IsValidatorRewardSufficient() bool {
	rules := DefaultValidationRules()
	validatorWeight := math.LegacyZeroDec()

	for _, dest := range p.Destinations {
		if dest.Enabled && dest.ModuleName == FeeCollectorModuleName {
			validatorWeight = dest.Weight
			break
		}
	}

	return validatorWeight.GTE(rules.MinValidatorReward)
}

// GetEnabledDestinations returns only enabled destinations
func (p EmissionParams) GetEnabledDestinations() []EmissionDestination {
	var enabled []EmissionDestination
	for _, dest := range p.Destinations {
		if dest.Enabled {
			enabled = append(enabled, dest)
		}
	}
	return enabled
}

// ToJSON converts EmissionParams to JSON string for storage
func (p EmissionParams) ToJSON() (string, error) {
	bytes, err := json.Marshal(p)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

// FromJSON converts JSON string to EmissionParams
func FromJSON(jsonStr string) (EmissionParams, error) {
	var params EmissionParams
	if err := json.Unmarshal([]byte(jsonStr), &params); err != nil {
		return EmissionParams{}, err
	}
	return params, nil
}

// CreateEmergencyParams creates emergency parameters that halt all emissions
func CreateEmergencyParams(reason string, authority string, blockHeight int64) EmissionParams {
	return EmissionParams{
		Destinations:    []EmissionDestination{}, // No destinations = all to fee_collector
		Enabled:         false,                   // Disable emissions
		UpdatedBy:       fmt.Sprintf("emergency_%d", blockHeight),
		UpdatedAt:       blockHeight,
		EmergencyStop:   true,
		FallbackEnabled: true,
	}
}

// CreateFallbackParams creates fallback parameters that revert to 50/50 split
func CreateFallbackParams(reason string, blockHeight int64) EmissionParams {
	params := DefaultEmissionParams()
	params.UpdatedBy = fmt.Sprintf("fallback_%d", blockHeight)
	params.UpdatedAt = blockHeight
	params.FallbackEnabled = true
	return params
}

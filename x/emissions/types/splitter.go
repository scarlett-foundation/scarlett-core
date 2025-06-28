package types

import (
	"fmt"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"
)

// EmissionSplitter handles the distribution of minted tokens according to configuration
type EmissionSplitter struct {
	config     EmissionsConfig
	bankKeeper bankkeeper.Keeper
}

// NewEmissionSplitter creates a new emission splitter with the given configuration
func NewEmissionSplitter(config EmissionsConfig, bankKeeper bankkeeper.Keeper) (*EmissionSplitter, error) {
	if err := config.Validate(); err != nil {
		return nil, fmt.Errorf("invalid emissions config: %w", err)
	}

	return &EmissionSplitter{
		config:     config,
		bankKeeper: bankKeeper,
	}, nil
}

// CalculateDistribution calculates how to split the given amount according to weights
func (es *EmissionSplitter) CalculateDistribution(totalAmount math.Int) ([]Distribution, error) {
	if !es.config.Enabled {
		// If disabled, send all to fee collector (traditional behavior)
		return []Distribution{
			{
				ModuleName:  FeeCollectorModuleName,
				Amount:      totalAmount,
				Description: "All emissions (traditional mint behavior)",
			},
		}, nil
	}

	distributions := make([]Distribution, 0, len(es.config.Destinations))
	allocatedAmount := math.ZeroInt()

	// Calculate weighted amounts for all destinations except the last
	for _, dest := range es.config.Destinations[:len(es.config.Destinations)-1] {
		// Amount = totalAmount * weight
		weightedAmount := math.LegacyNewDecFromInt(totalAmount).Mul(dest.Weight).TruncateInt()

		distributions = append(distributions, Distribution{
			ModuleName:  dest.ModuleName,
			Amount:      weightedAmount,
			Description: dest.Description,
		})

		allocatedAmount = allocatedAmount.Add(weightedAmount)
	}

	// Last destination gets the remainder to ensure exact conservation
	lastDest := es.config.Destinations[len(es.config.Destinations)-1]
	remainingAmount := totalAmount.Sub(allocatedAmount)

	distributions = append(distributions, Distribution{
		ModuleName:  lastDest.ModuleName,
		Amount:      remainingAmount,
		Description: lastDest.Description,
	})

	return distributions, nil
}

// DistributeTokens distributes tokens from the mint module to configured destinations
func (es *EmissionSplitter) DistributeTokens(ctx sdk.Context, mintDenom string, totalAmount math.Int) error {
	if totalAmount.IsZero() {
		return nil // Nothing to distribute
	}

	distributions, err := es.CalculateDistribution(totalAmount)
	if err != nil {
		return fmt.Errorf("failed to calculate distribution: %w", err)
	}

	// Execute transfers
	for _, dist := range distributions {
		if dist.Amount.IsZero() {
			continue // Skip zero transfers
		}

		coins := sdk.NewCoins(sdk.NewCoin(mintDenom, dist.Amount))
		if err := es.bankKeeper.SendCoinsFromModuleToModule(
			ctx,
			minttypes.ModuleName,
			dist.ModuleName,
			coins,
		); err != nil {
			return fmt.Errorf("failed to transfer %s to module %s: %w", coins.String(), dist.ModuleName, err)
		}
	}

	return nil
}

// GetConfig returns the current emissions configuration
func (es *EmissionSplitter) GetConfig() EmissionsConfig {
	return es.config
}

// UpdateConfig updates the emissions configuration with validation
func (es *EmissionSplitter) UpdateConfig(newConfig EmissionsConfig) error {
	if err := newConfig.Validate(); err != nil {
		return fmt.Errorf("invalid new config: %w", err)
	}

	es.config = newConfig
	return nil
}

// Distribution represents a calculated token distribution
type Distribution struct {
	ModuleName  string   `json:"module_name"`
	Amount      math.Int `json:"amount"`
	Description string   `json:"description"`
}

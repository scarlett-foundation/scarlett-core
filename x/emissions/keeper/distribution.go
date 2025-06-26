package keeper

import (
	"fmt"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	mintkeeper "github.com/cosmos/cosmos-sdk/x/mint/keeper"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"

	"scarlett-core/app/emissions"
	"scarlett-core/x/emissions/types"
)

// ProvideDynamicMintFn creates a dynamic mint function that uses governance-controlled emission parameters
func (k Keeper) ProvideDynamicMintFn() func(ctx sdk.Context, mintKeeper *mintkeeper.Keeper) error {
	return func(ctx sdk.Context, mintKeeper *mintkeeper.Keeper) error {
		return k.DynamicEmissionsMintFn(ctx)
	}
}

// DynamicEmissionsMintFn implements dynamic emission splitting using governance-controlled parameters
func (k Keeper) DynamicEmissionsMintFn(ctx sdk.Context) error {
	// 1. Get governance-controlled emission parameters
	emissionParams, err := k.GetEmissionParams(ctx)
	if err != nil {
		// Fallback to default parameters if not found
		emissionParams = types.DefaultEmissionParams()
		// Store default parameters for future use
		if err := k.SetEmissionParams(ctx, emissionParams); err != nil {
			return fmt.Errorf("failed to store default emission parameters: %w", err)
		}
	}

	// 2. Convert governance parameters to emissions config
	config := k.convertToEmissionsConfig(emissionParams)

	// 3. Create emission splitter using existing logic
	splitter, err := emissions.NewEmissionSplitter(config, k.bankKeeper)
	if err != nil {
		return fmt.Errorf("failed to create emission splitter: %w", err)
	}

	// 4. Get current minter state and parameters using our mint keeper
	minter, err := k.getMinterState(ctx)
	if err != nil {
		return fmt.Errorf("failed to get minter state: %w", err)
	}

	params, err := k.mintKeeper.Params.Get(ctx)
	if err != nil {
		return fmt.Errorf("failed to get mint params: %w", err)
	}

	// 5. Calculate inflation and provisions using standard mint logic
	if err := k.updateInflation(ctx, &minter, params); err != nil {
		return fmt.Errorf("failed to update inflation: %w", err)
	}

	// 6. Calculate and mint this block's provision
	blockProvisionCoin := minter.BlockProvision(params)
	if err := k.mintTokens(ctx, blockProvisionCoin); err != nil {
		return fmt.Errorf("failed to mint tokens: %w", err)
	}

	// 7. Distribute minted tokens using governance-controlled splitter
	if err := splitter.DistributeTokens(ctx, params.MintDenom, blockProvisionCoin.Amount); err != nil {
		return fmt.Errorf("failed to distribute tokens: %w", err)
	}

	// 8. Update minter state in the store
	if err := k.mintKeeper.Minter.Set(ctx, minter); err != nil {
		return fmt.Errorf("failed to update minter state: %w", err)
	}

	// 9. Update destination metrics
	if err := k.updateDestinationMetrics(ctx, emissionParams.Destinations, blockProvisionCoin.Amount, params.MintDenom); err != nil {
		return fmt.Errorf("failed to update destination metrics: %w", err)
	}

	// 10. Emit comprehensive events
	distributions, _ := splitter.CalculateDistribution(blockProvisionCoin.Amount)
	bondedRatio, _ := k.stakingKeeper.BondedRatio(ctx)

	k.emitDynamicEmissionEvent(
		ctx,
		blockProvisionCoin.Amount,
		distributions,
		emissionParams,
		minter,
		bondedRatio,
	)

	return nil
}

// convertToEmissionsConfig converts governance parameters to the existing emissions config format
func (k Keeper) convertToEmissionsConfig(params types.EmissionParams) emissions.EmissionsConfig {
	if !params.Enabled {
		// Return disabled config - all to fee collector
		return emissions.EmissionsConfig{
			Enabled:      false,
			Destinations: []emissions.EmissionDestination{},
		}
	}

	// Convert governance destinations to emissions destinations
	destinations := make([]emissions.EmissionDestination, len(params.Destinations))
	for i, dest := range params.Destinations {
		if !dest.Enabled {
			continue // Skip disabled destinations
		}

		destinations[i] = emissions.EmissionDestination{
			ModuleName:  dest.ModuleName,
			Weight:      dest.Weight,
			Description: dest.Description,
		}
	}

	return emissions.EmissionsConfig{
		Enabled:      true,
		Destinations: destinations,
	}
}

// getMinterState retrieves and validates the current minter state
func (k Keeper) getMinterState(ctx sdk.Context) (minttypes.Minter, error) {
	minter, err := k.mintKeeper.Minter.Get(ctx)
	if err != nil {
		return minttypes.Minter{}, fmt.Errorf("failed to get minter: %w", err)
	}
	return minter, nil
}

// updateInflation calculates and updates inflation rates and annual provisions
func (k Keeper) updateInflation(ctx sdk.Context, minter *minttypes.Minter, params minttypes.Params) error {
	// Get staking metrics
	stakingTokenSupply, err := k.mintKeeper.StakingTokenSupply(ctx)
	if err != nil {
		return fmt.Errorf("failed to get staking token supply: %w", err)
	}

	bondedRatio, err := k.stakingKeeper.BondedRatio(ctx)
	if err != nil {
		return fmt.Errorf("failed to get bonded ratio: %w", err)
	}

	// Calculate inflation using standard mint logic
	minter.Inflation = minter.NextInflationRate(params, bondedRatio)
	minter.AnnualProvisions = minter.NextAnnualProvisions(params, stakingTokenSupply)

	return nil
}

// mintTokens mints the specified amount of tokens to the mint module account
func (k Keeper) mintTokens(ctx sdk.Context, coin sdk.Coin) error {
	if coin.Amount.IsZero() {
		return nil // Nothing to mint
	}

	coins := sdk.NewCoins(coin)
	if err := k.mintKeeper.MintCoins(ctx, coins); err != nil {
		return fmt.Errorf("failed to mint %s: %w", coins.String(), err)
	}

	return nil
}

// updateDestinationMetrics updates statistics for each emission destination
func (k Keeper) updateDestinationMetrics(ctx sdk.Context, destinations []types.EmissionDestination, totalAmount math.Int, denom string) error {
	// Calculate distribution amounts
	for _, dest := range destinations {
		if !dest.Enabled {
			continue
		}

		// Calculate amount for this destination
		destAmount := math.LegacyNewDecFromInt(totalAmount).Mul(dest.Weight).TruncateInt()
		coins := sdk.NewCoins(sdk.NewCoin(denom, destAmount))

		// Update metrics (simplified for now - will implement full metrics later)
		// This is a placeholder for the full metrics implementation
		_ = coins // Use the calculated coins
	}

	return nil
}

// emitDynamicEmissionEvent emits a comprehensive event for dynamic emission distribution
func (k Keeper) emitDynamicEmissionEvent(
	ctx sdk.Context,
	totalAmount math.Int,
	distributions []emissions.Distribution,
	emissionParams types.EmissionParams,
	minter minttypes.Minter,
	bondedRatio math.LegacyDec,
) {
	// Build event attributes
	attributes := []sdk.Attribute{
		sdk.NewAttribute(minttypes.AttributeKeyBondedRatio, bondedRatio.String()),
		sdk.NewAttribute(minttypes.AttributeKeyInflation, minter.Inflation.String()),
		sdk.NewAttribute(minttypes.AttributeKeyAnnualProvisions, minter.AnnualProvisions.String()),
		sdk.NewAttribute(types.AttributeKeyAmount, totalAmount.String()),
		sdk.NewAttribute("num_destinations", fmt.Sprintf("%d", len(distributions))),
		sdk.NewAttribute("governance_controlled", "true"),
		sdk.NewAttribute("updated_by", emissionParams.UpdatedBy),
		sdk.NewAttribute("updated_at", fmt.Sprintf("%d", emissionParams.UpdatedAt)),
	}

	// Add distribution details
	for i, dist := range distributions {
		prefix := fmt.Sprintf("destination_%d", i)
		attributes = append(attributes,
			sdk.NewAttribute(prefix+"_module", dist.ModuleName),
			sdk.NewAttribute(prefix+"_amount", dist.Amount.String()),
			sdk.NewAttribute(prefix+"_description", dist.Description),
		)
	}

	// Emit the dynamic emission event
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(types.EventTypeEmissionDistributed, attributes...),
	)

	// Also emit the standard mint event for compatibility
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			minttypes.EventTypeMint,
			sdk.NewAttribute(minttypes.AttributeKeyBondedRatio, bondedRatio.String()),
			sdk.NewAttribute(minttypes.AttributeKeyInflation, minter.Inflation.String()),
			sdk.NewAttribute(minttypes.AttributeKeyAnnualProvisions, minter.AnnualProvisions.String()),
			sdk.NewAttribute(sdk.AttributeKeyAmount, totalAmount.String()),
		),
	)
}

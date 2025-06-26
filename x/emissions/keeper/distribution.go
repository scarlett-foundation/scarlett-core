package keeper

import (
	"encoding/json"
	"fmt"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	mintkeeper "github.com/cosmos/cosmos-sdk/x/mint/keeper"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"

	"scarlett-core/app/emissions"
	"scarlett-core/x/emissions/types"
)

// ProvideDynamicMintFn creates a dynamic mint function that uses governance-controlled emission parameters
func (k Keeper) ProvideDynamicMintFn() mintkeeper.MintFn {
	return func(ctx sdk.Context, mintKeeper *mintkeeper.Keeper) error {
		return k.DynamicEmissionsMintFn(ctx, mintKeeper)
	}
}

// DynamicEmissionsMintFn implements the MintFn interface for governance-controlled emissions
func (k Keeper) DynamicEmissionsMintFn(ctx sdk.Context, mintKeeper *mintkeeper.Keeper) error {
	// 1. Get governance-controlled parameters
	params, err := k.Params.Get(ctx)
	if err != nil || !params.Enabled {
		// Fallback to default behavior if governance params not set
		return k.applyFallbackConfiguration(ctx, mintKeeper)
	}

	// 2. Parse governance parameters
	var destinations []types.EmissionDestination
	if err := json.Unmarshal([]byte(params.EmissionDestinations), &destinations); err != nil {
		return k.activateEmergencyFallback(ctx, mintKeeper, "invalid governance parameters")
	}

	// 3. Convert to emissions config
	emissionsConfig := k.convertGovernanceToEmissionsConfig(destinations)

	// 4. Create emission splitter using existing logic
	splitter, err := emissions.NewEmissionSplitter(emissionsConfig, k.bankKeeper)
	if err != nil {
		return fmt.Errorf("failed to create emission splitter: %w", err)
	}

	// 5. Get minter state and mint tokens
	minter, err := k.getMinterState(ctx, mintKeeper)
	if err != nil {
		return err
	}

	mintParams, err := mintKeeper.Params.Get(ctx)
	if err != nil {
		return fmt.Errorf("failed to get mint params: %w", err)
	}

	// 6. Calculate inflation and provisions using standard mint logic
	if err := k.updateInflation(ctx, &minter, mintParams, mintKeeper); err != nil {
		return fmt.Errorf("failed to update inflation: %w", err)
	}

	// 7. Calculate and mint this block's provision
	blockProvisionCoin := minter.BlockProvision(mintParams)
	coins := sdk.NewCoins(blockProvisionCoin)
	if err := mintKeeper.MintCoins(ctx, coins); err != nil {
		return fmt.Errorf("failed to mint tokens: %w", err)
	}

	// 8. Distribute according to governance configuration
	if err := splitter.DistributeTokens(ctx, mintParams.MintDenom, blockProvisionCoin.Amount); err != nil {
		return fmt.Errorf("failed to distribute tokens: %w", err)
	}

	// 9. Update minter state in the store
	if err := mintKeeper.Minter.Set(ctx, minter); err != nil {
		return fmt.Errorf("failed to update minter state: %w", err)
	}

	// 10. Emit governance events
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeEmissionDistributed,
			sdk.NewAttribute(types.AttributeKeyAmount, blockProvisionCoin.String()),
			sdk.NewAttribute("governance_controlled", "true"),
			sdk.NewAttribute("destinations", params.EmissionDestinations),
		),
	)

	return nil
}

// convertGovernanceToEmissionsConfig converts governance parameters to emissions config
func (k Keeper) convertGovernanceToEmissionsConfig(destinations []types.EmissionDestination) emissions.EmissionsConfig {
	var emissionDestinations []emissions.EmissionDestination

	for _, dest := range destinations {
		if !dest.Enabled {
			continue
		}
		emissionDestinations = append(emissionDestinations, emissions.EmissionDestination{
			ModuleName:  dest.ModuleName,
			Weight:      dest.Weight,
			Description: dest.Description,
		})
	}

	return emissions.EmissionsConfig{
		Enabled:      true,
		Destinations: emissionDestinations,
	}
}

// getMinterState retrieves and validates the current minter state
func (k Keeper) getMinterState(ctx sdk.Context, mintKeeper *mintkeeper.Keeper) (minttypes.Minter, error) {
	minter, err := mintKeeper.Minter.Get(ctx)
	if err != nil {
		return minttypes.Minter{}, fmt.Errorf("failed to get minter: %w", err)
	}
	return minter, nil
}

// updateInflation calculates and updates inflation rates and annual provisions
func (k Keeper) updateInflation(ctx sdk.Context, minter *minttypes.Minter, params minttypes.Params, mintKeeper *mintkeeper.Keeper) error {
	// Get staking metrics
	stakingTokenSupply, err := mintKeeper.StakingTokenSupply(ctx)
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

// EMERGENCY CONTROL FUNCTIONS

// handleEmergencyStop processes emergency stop conditions
func (k Keeper) handleEmergencyStop(ctx sdk.Context, params types.EmissionParams) error {
	ctx.Logger().Info("Emergency stop is active, halting emissions",
		"reason", "emergency_stop", "height", ctx.BlockHeight())

	// Emit emergency stop event
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeEmergencyStop,
			sdk.NewAttribute("reason", "emergency_stop_active"),
			sdk.NewAttribute("block_height", fmt.Sprintf("%d", ctx.BlockHeight())),
			sdk.NewAttribute("updated_by", params.UpdatedBy),
		),
	)

	// During emergency stop, we don't mint or distribute any tokens
	// This effectively halts all emissions until governance resolves the emergency
	return nil
}

// applyFallbackConfiguration applies fallback configuration when enabled
func (k Keeper) applyFallbackConfiguration(ctx sdk.Context, mintKeeper *mintkeeper.Keeper) error {
	// If fallback is enabled, revert to safe default 50/50 configuration
	fallbackParams := types.DefaultEmissionParams()
	fallbackParams.FallbackEnabled = true
	fallbackParams.UpdatedBy = "automatic_fallback"
	fallbackParams.UpdatedAt = ctx.BlockHeight()

	// Store emergency parameters
	if err := k.SetEmissionParams(ctx, fallbackParams); err != nil {
		return fmt.Errorf("failed to store emergency parameters: %w", err)
	}

	// Store in history for audit trail
	if err := k.SetEmissionHistory(ctx, ctx.BlockHeight(), fallbackParams); err != nil {
		ctx.Logger().Error("Failed to store emergency history", "error", err)
	}

	// Emit emergency activation event
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeEmergencyStop,
			sdk.NewAttribute("reason", "automatic_fallback"),
			sdk.NewAttribute("emergency_type", "automatic_fallback"),
			sdk.NewAttribute("block_height", fmt.Sprintf("%d", ctx.BlockHeight())),
			sdk.NewAttribute("fallback_active", "true"),
		),
	)

	ctx.Logger().Error("Emergency fallback activated",
		"reason", "automatic_fallback", "height", ctx.BlockHeight())

	// Continue with fallback configuration
	// Note: This is a recursive call from within the mint function, so we need to
	// use the original mint keeper reference. Since we're in an emergency fallback,
	// we'll use the keeper's mint keeper reference.
	return k.DynamicEmissionsMintFn(ctx, mintKeeper)
}

// activateEmergencyFallback activates emergency fallback when critical errors occur
func (k Keeper) activateEmergencyFallback(ctx sdk.Context, mintKeeper *mintkeeper.Keeper, reason string) error {
	// Create emergency fallback parameters
	emergencyParams := types.CreateFallbackParams(reason, ctx.BlockHeight())

	// Store emergency parameters
	if err := k.SetEmissionParams(ctx, emergencyParams); err != nil {
		return fmt.Errorf("failed to store emergency parameters: %w", err)
	}

	// Store in history for audit trail
	if err := k.SetEmissionHistory(ctx, ctx.BlockHeight(), emergencyParams); err != nil {
		ctx.Logger().Error("Failed to store emergency history", "error", err)
	}

	// Emit emergency activation event
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeEmergencyStop,
			sdk.NewAttribute("reason", reason),
			sdk.NewAttribute("emergency_type", "automatic_fallback"),
			sdk.NewAttribute("block_height", fmt.Sprintf("%d", ctx.BlockHeight())),
			sdk.NewAttribute("fallback_active", "true"),
		),
	)

	ctx.Logger().Error("Emergency fallback activated",
		"reason", reason, "height", ctx.BlockHeight())

	// Continue with fallback configuration
	// Note: This is a recursive call from within the mint function, so we need to
	// use the original mint keeper reference. Since we're in an emergency fallback,
	// we'll use the keeper's mint keeper reference.
	return k.DynamicEmissionsMintFn(ctx, mintKeeper)
}

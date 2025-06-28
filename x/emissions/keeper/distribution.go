package keeper

import (
	"encoding/json"
	"fmt"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	mintkeeper "github.com/cosmos/cosmos-sdk/x/mint/keeper"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"

	"scarlett-core/x/emissions/types"
)

// MintAndDistributeEmissions is the consolidated autonomous emissions method
// that replaces all complex logic from both keeper and app layers
func (k Keeper) MintAndDistributeEmissions(ctx sdk.Context, mintKeeper *mintkeeper.Keeper) error {
	// EXACT PRESERVED LOGGING: Combine both entry patterns
	ctx.Logger().Info("🔥🔥🔥 GOVERNANCE-CONTROLLED MINT FUNCTION CALLED 🔥🔥🔥", "height", ctx.BlockHeight())
	ctx.Logger().Info("🔥🔥🔥 GOVERNANCE-CONTROLLED EMISSIONS ACTIVE 🔥🔥🔥",
		"block_height", ctx.BlockHeight(),
		"module", "emissions",
		"function", "MintAndDistributeEmissions")

	// 1. Get governance parameters (EXACT PRESERVED LOGGING)
	emissionParams, err := k.Params.Get(ctx)
	if err != nil || !emissionParams.Enabled {
		if err != nil {
			ctx.Logger().Info("Could not get governance parameters, falling back to default 50/50 split", "error", err)
		} else {
			ctx.Logger().Info("Governance emissions are explicitly disabled, falling back to default 50/50 split")
		}
		return k.applyFallbackConfiguration(ctx, mintKeeper)
	} else {
		ctx.Logger().Info("Using governance-defined emission parameters")
		ctx.Logger().Info("📋 Governance parameters retrieved",
			"enabled", emissionParams.Enabled,
			"destinations", emissionParams.EmissionDestinations)
	}

	// 2. Parse governance destinations (EXACT PRESERVED LOGGING)
	var destinations []struct {
		ModuleName  string `json:"module_name"`
		Weight      string `json:"weight"`
		Description string `json:"description"`
		Enabled     bool   `json:"enabled"`
	}

	if err := json.Unmarshal([]byte(emissionParams.EmissionDestinations), &destinations); err != nil {
		ctx.Logger().Error("Failed to parse governance emission destinations JSON, falling back to default", "error", err)
		return k.applyFallbackConfiguration(ctx, mintKeeper)
	}

	ctx.Logger().Info("✅ Parsed governance destinations",
		"num_destinations", len(destinations),
		"destinations", destinations)

	// 3. Convert to emissions config (EXACT PRESERVED LOGGING)
	var emissionDestinations []types.EmissionDestination
	for _, dest := range destinations {
		if !dest.Enabled {
			continue // Skip disabled destinations
		}
		weight, err := math.LegacyNewDecFromStr(dest.Weight)
		if err != nil {
			ctx.Logger().Error("Invalid weight in governance destination, falling back to default", "module", dest.ModuleName, "error", err)
			return k.applyFallbackConfiguration(ctx, mintKeeper)
		}
		emissionDestinations = append(emissionDestinations, types.EmissionDestination{
			ModuleName:  dest.ModuleName,
			Weight:      weight,
			Description: dest.Description,
			Enabled:     true,
			MinWeight:   math.LegacyNewDecWithPrec(30, 2), // 30% minimum
			MaxWeight:   math.LegacyNewDecWithPrec(70, 2), // 70% maximum
		})
	}

	config := k.convertGovernanceToEmissionsConfig(emissionDestinations)
	ctx.Logger().Info("📊 Using emission configuration",
		"enabled", config.Enabled,
		"destinations", len(config.Destinations))

	// 4. Create emission splitter (EXACT PRESERVED LOGGING)
	splitter, err := types.NewEmissionSplitter(config, k.bankKeeper)
	if err != nil {
		return fmt.Errorf("failed to create emission splitter: %w", err)
	}
	ctx.Logger().Info("✅ Emission splitter created successfully")

	// 5. Get minter state and parameters (EXACT PRESERVED LOGGING)
	minter, err := k.getMinterState(ctx, mintKeeper)
	if err != nil {
		return err
	}

	ctx.Logger().Info("📊 Minter state retrieved",
		"inflation", minter.Inflation.String(),
		"annual_provisions", minter.AnnualProvisions.String())

	mintParams, err := mintKeeper.Params.Get(ctx)
	if err != nil {
		return fmt.Errorf("failed to get mint params: %w", err)
	}

	ctx.Logger().Info("⚙️ Mint parameters retrieved",
		"mint_denom", mintParams.MintDenom,
		"inflation_rate_change", mintParams.InflationRateChange.String(),
		"inflation_max", mintParams.InflationMax.String(),
		"inflation_min", mintParams.InflationMin.String(),
		"goal_bonded", mintParams.GoalBonded.String(),
		"blocks_per_year", mintParams.BlocksPerYear)

	// 6. Calculate inflation and provisions (EXACT PRESERVED LOGGING)
	if err := k.updateInflation(ctx, &minter, mintParams, mintKeeper); err != nil {
		return fmt.Errorf("failed to update inflation: %w", err)
	}

	ctx.Logger().Info("📈 Inflation updated",
		"new_inflation", minter.Inflation.String(),
		"new_annual_provisions", minter.AnnualProvisions.String())

	// 7. Calculate and mint block provision (EXACT PRESERVED LOGGING)
	blockProvisionCoin := minter.BlockProvision(mintParams)
	ctx.Logger().Info("💰 Block provision calculated",
		"block_provision", blockProvisionCoin.String(),
		"amount", blockProvisionCoin.Amount.String(),
		"denom", blockProvisionCoin.Denom)

	if blockProvisionCoin.Amount.IsZero() {
		ctx.Logger().Info("⚠️ Zero block provision, skipping mint and distribution")
		return nil
	}

	coins := sdk.NewCoins(blockProvisionCoin)
	if err := mintKeeper.MintCoins(ctx, coins); err != nil {
		return fmt.Errorf("failed to mint tokens: %w", err)
	}

	ctx.Logger().Info("✅ Tokens minted successfully",
		"minted_coins", coins.String())
	ctx.Logger().Info("💰 Minted tokens for distribution",
		"amount", blockProvisionCoin.Amount.String(),
		"denom", blockProvisionCoin.Denom)

	// 8. Calculate distribution (EXACT PRESERVED LOGGING)
	distributions, err := splitter.CalculateDistribution(blockProvisionCoin.Amount)
	if err != nil {
		return fmt.Errorf("failed to calculate distribution: %w", err)
	}

	ctx.Logger().Info("📊 Distribution calculated",
		"num_distributions", len(distributions))

	// EXACT PRESERVED LOGGING: Per-destination details
	for i, dist := range distributions {
		ctx.Logger().Info("💸 Distribution detail",
			"index", i,
			"module", dist.ModuleName,
			"amount", dist.Amount.String(),
			"description", dist.Description)
	}

	// 9. Execute distribution (EXACT PRESERVED LOGGING)
	if err := splitter.DistributeTokens(ctx, mintParams.MintDenom, blockProvisionCoin.Amount); err != nil {
		return fmt.Errorf("failed to distribute tokens: %w", err)
	}

	ctx.Logger().Info("✅ Tokens distributed successfully")

	// 10. Update minter state (EXACT PRESERVED LOGGING)
	if err := mintKeeper.Minter.Set(ctx, minter); err != nil {
		return fmt.Errorf("failed to update minter state: %w", err)
	}

	ctx.Logger().Info("✅ Minter state updated successfully")

	// 11. Emit comprehensive events and final success (EXACT PRESERVED LOGGING)
	bondedRatio, _ := k.stakingKeeper.BondedRatio(ctx)
	types.EmitEmissionSplitEvent(ctx, blockProvisionCoin.Amount, distributions, config, minter, bondedRatio)

	// EXACT PRESERVED LOGGING: Both success completion patterns
	ctx.Logger().Info("✅ Governance-controlled emission distribution complete")
	ctx.Logger().Info("🎉 GOVERNANCE-CONTROLLED EMISSIONS COMPLETED SUCCESSFULLY 🎉",
		"block_height", ctx.BlockHeight(),
		"total_minted", blockProvisionCoin.String(),
		"num_distributions", len(distributions))

	return nil
}

// ProvideDynamicMintFn creates a dynamic mint function that uses governance-controlled emission parameters
func (k Keeper) ProvideDynamicMintFn() mintkeeper.MintFn {
	return func(ctx sdk.Context, mintKeeper *mintkeeper.Keeper) error {
		// CRITICAL DEBUG: This should appear if our function is being called
		ctx.Logger().Info("🚨🚨🚨 CUSTOM MINT FUNCTION CALLED 🚨🚨🚨", "height", ctx.BlockHeight())
		return k.DynamicEmissionsMintFn(ctx, mintKeeper)
	}
}

// DynamicEmissionsMintFn implements the MintFn interface for governance-controlled emissions
func (k Keeper) DynamicEmissionsMintFn(ctx sdk.Context, mintKeeper *mintkeeper.Keeper) error {
	// DEBUG: Log that our function is being called
	ctx.Logger().Info("🔥🔥🔥 GOVERNANCE-CONTROLLED EMISSIONS ACTIVE 🔥🔥🔥",
		"block_height", ctx.BlockHeight(),
		"module", "emissions",
		"function", "DynamicEmissionsMintFn")

	// 1. Get governance-controlled parameters
	params, err := k.Params.Get(ctx)
	if err != nil {
		ctx.Logger().Info("❌ Failed to get governance parameters, using fallback",
			"error", err.Error(),
			"fallback_reason", "params_get_error")
		return k.applyFallbackConfiguration(ctx, mintKeeper)
	}

	ctx.Logger().Info("📋 Governance parameters retrieved",
		"enabled", params.Enabled,
		"destinations", params.EmissionDestinations)

	if !params.Enabled {
		ctx.Logger().Info("⚠️ Governance parameters disabled, using fallback",
			"enabled", params.Enabled,
			"fallback_reason", "params_disabled")
		return k.applyFallbackConfiguration(ctx, mintKeeper)
	}

	// 2. Parse governance parameters
	var destinations []types.EmissionDestination
	if params.EmissionDestinations == "" {
		ctx.Logger().Info("⚠️ Empty emission destinations, using fallback",
			"destinations", params.EmissionDestinations,
			"fallback_reason", "empty_destinations")
		return k.applyFallbackConfiguration(ctx, mintKeeper)
	}

	if err := json.Unmarshal([]byte(params.EmissionDestinations), &destinations); err != nil {
		ctx.Logger().Error("❌ Failed to parse governance parameters, using emergency fallback",
			"error", err.Error(),
			"raw_destinations", params.EmissionDestinations,
			"fallback_reason", "json_parse_error")
		return k.activateEmergencyFallback(ctx, mintKeeper, "invalid governance parameters")
	}

	ctx.Logger().Info("✅ Parsed governance destinations",
		"num_destinations", len(destinations),
		"destinations", destinations)

	// 3. Convert to emissions config
	emissionsConfig := k.convertGovernanceToEmissionsConfig(destinations)
	ctx.Logger().Info("🔄 Converted to emissions config",
		"enabled", emissionsConfig.Enabled,
		"num_destinations", len(emissionsConfig.Destinations))

	// 4. Create emission splitter using existing logic
	splitter, err := types.NewEmissionSplitter(emissionsConfig, k.bankKeeper)
	if err != nil {
		ctx.Logger().Error("❌ Failed to create emission splitter",
			"error", err.Error())
		return fmt.Errorf("failed to create emission splitter: %w", err)
	}

	ctx.Logger().Info("✅ Emission splitter created successfully")

	// 5. Get minter state and mint tokens
	minter, err := k.getMinterState(ctx, mintKeeper)
	if err != nil {
		ctx.Logger().Error("❌ Failed to get minter state",
			"error", err.Error())
		return err
	}

	ctx.Logger().Info("📊 Minter state retrieved",
		"inflation", minter.Inflation.String(),
		"annual_provisions", minter.AnnualProvisions.String())

	mintParams, err := mintKeeper.Params.Get(ctx)
	if err != nil {
		ctx.Logger().Error("❌ Failed to get mint params",
			"error", err.Error())
		return fmt.Errorf("failed to get mint params: %w", err)
	}

	ctx.Logger().Info("⚙️ Mint parameters retrieved",
		"mint_denom", mintParams.MintDenom,
		"inflation_rate_change", mintParams.InflationRateChange.String(),
		"inflation_max", mintParams.InflationMax.String(),
		"inflation_min", mintParams.InflationMin.String(),
		"goal_bonded", mintParams.GoalBonded.String(),
		"blocks_per_year", mintParams.BlocksPerYear)

	// 6. Calculate inflation and provisions using standard mint logic
	if err := k.updateInflation(ctx, &minter, mintParams, mintKeeper); err != nil {
		ctx.Logger().Error("❌ Failed to update inflation",
			"error", err.Error())
		return fmt.Errorf("failed to update inflation: %w", err)
	}

	ctx.Logger().Info("📈 Inflation updated",
		"new_inflation", minter.Inflation.String(),
		"new_annual_provisions", minter.AnnualProvisions.String())

	// 7. Calculate and mint this block's provision
	blockProvisionCoin := minter.BlockProvision(mintParams)
	ctx.Logger().Info("💰 Block provision calculated",
		"block_provision", blockProvisionCoin.String(),
		"amount", blockProvisionCoin.Amount.String(),
		"denom", blockProvisionCoin.Denom)

	if blockProvisionCoin.Amount.IsZero() {
		ctx.Logger().Info("⚠️ Zero block provision, skipping mint and distribution")
		return nil
	}

	coins := sdk.NewCoins(blockProvisionCoin)
	if err := mintKeeper.MintCoins(ctx, coins); err != nil {
		ctx.Logger().Error("❌ Failed to mint tokens",
			"error", err.Error(),
			"coins", coins.String())
		return fmt.Errorf("failed to mint tokens: %w", err)
	}

	ctx.Logger().Info("✅ Tokens minted successfully",
		"minted_coins", coins.String())

	// 8. Calculate distribution before executing
	distributions, err := splitter.CalculateDistribution(blockProvisionCoin.Amount)
	if err != nil {
		ctx.Logger().Error("❌ Failed to calculate distribution",
			"error", err.Error())
		return fmt.Errorf("failed to calculate distribution: %w", err)
	}

	ctx.Logger().Info("📊 Distribution calculated",
		"num_distributions", len(distributions))

	for i, dist := range distributions {
		ctx.Logger().Info("💸 Distribution detail",
			"index", i,
			"module", dist.ModuleName,
			"amount", dist.Amount.String(),
			"description", dist.Description)
	}

	// 9. Distribute according to governance configuration
	if err := splitter.DistributeTokens(ctx, mintParams.MintDenom, blockProvisionCoin.Amount); err != nil {
		ctx.Logger().Error("❌ Failed to distribute tokens",
			"error", err.Error(),
			"total_amount", blockProvisionCoin.Amount.String(),
			"denom", mintParams.MintDenom)
		return fmt.Errorf("failed to distribute tokens: %w", err)
	}

	ctx.Logger().Info("✅ Tokens distributed successfully")

	// 10. Update minter state in the store
	if err := mintKeeper.Minter.Set(ctx, minter); err != nil {
		ctx.Logger().Error("❌ Failed to update minter state",
			"error", err.Error())
		return fmt.Errorf("failed to update minter state: %w", err)
	}

	ctx.Logger().Info("✅ Minter state updated successfully")

	// 11. Emit governance events
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeEmissionDistributed,
			sdk.NewAttribute(types.AttributeKeyAmount, blockProvisionCoin.String()),
			sdk.NewAttribute("governance_controlled", "true"),
			sdk.NewAttribute("destinations", params.EmissionDestinations),
			sdk.NewAttribute("num_distributions", fmt.Sprintf("%d", len(distributions))),
		),
	)

	ctx.Logger().Info("🎉 GOVERNANCE-CONTROLLED EMISSIONS COMPLETED SUCCESSFULLY 🎉",
		"block_height", ctx.BlockHeight(),
		"total_minted", blockProvisionCoin.String(),
		"num_distributions", len(distributions))

	return nil
}

// convertGovernanceToEmissionsConfig converts governance parameters to emissions config
func (k Keeper) convertGovernanceToEmissionsConfig(destinations []types.EmissionDestination) types.EmissionsConfig {
	var emissionDestinations []types.EmissionDestination

	for _, dest := range destinations {
		if !dest.Enabled {
			continue
		}
		emissionDestinations = append(emissionDestinations, types.EmissionDestination{
			ModuleName:  dest.ModuleName,
			Weight:      dest.Weight,
			Description: dest.Description,
			Enabled:     dest.Enabled,
			MinWeight:   dest.MinWeight,
			MaxWeight:   dest.MaxWeight,
		})
	}

	return types.EmissionsConfig{
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
	distributions []types.Distribution,
	emissionParams types.Params,
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
func (k Keeper) handleEmergencyStop(ctx sdk.Context, params types.Params) error {
	ctx.Logger().Info("🛑 Emergency stop is active, halting emissions",
		"reason", "emergency_stop", "height", ctx.BlockHeight())

	// Emit emergency stop event
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeEmergencyStop,
			sdk.NewAttribute("reason", "emergency_stop_active"),
			sdk.NewAttribute("block_height", fmt.Sprintf("%d", ctx.BlockHeight())),
		),
	)

	// During emergency stop, we don't mint or distribute any tokens
	// This effectively halts all emissions until governance resolves the emergency
	return nil
}

// applyFallbackConfiguration applies default 50/50 split when governance params are unavailable
func (k Keeper) applyFallbackConfiguration(ctx sdk.Context, mintKeeper *mintkeeper.Keeper) error {
	ctx.Logger().Info("🔄 APPLYING FALLBACK CONFIGURATION (50/50 SPLIT)",
		"block_height", ctx.BlockHeight(),
		"reason", "governance_params_unavailable")

	// 1. Use hardcoded default 50/50 configuration
	config := types.DefaultEmissionsConfig()
	splitter, err := types.NewEmissionSplitter(config, k.bankKeeper)
	if err != nil {
		ctx.Logger().Error("❌ Failed to create fallback emission splitter",
			"error", err.Error())
		return fmt.Errorf("failed to create fallback emission splitter: %w", err)
	}

	ctx.Logger().Info("✅ Fallback emission splitter created",
		"config", config)

	// 2. Get minter state and mint tokens
	minter, err := k.getMinterState(ctx, mintKeeper)
	if err != nil {
		ctx.Logger().Error("❌ Failed to get minter state in fallback",
			"error", err.Error())
		return err
	}

	ctx.Logger().Info("📊 Fallback minter state retrieved",
		"inflation", minter.Inflation.String(),
		"annual_provisions", minter.AnnualProvisions.String())

	mintParams, err := mintKeeper.Params.Get(ctx)
	if err != nil {
		ctx.Logger().Error("❌ Failed to get mint params in fallback",
			"error", err.Error())
		return fmt.Errorf("failed to get mint params: %w", err)
	}

	// 3. Calculate inflation and provisions using standard mint logic
	if err := k.updateInflation(ctx, &minter, mintParams, mintKeeper); err != nil {
		ctx.Logger().Error("❌ Failed to update inflation in fallback",
			"error", err.Error())
		return fmt.Errorf("failed to update inflation: %w", err)
	}

	ctx.Logger().Info("📈 Fallback inflation updated",
		"new_inflation", minter.Inflation.String(),
		"new_annual_provisions", minter.AnnualProvisions.String())

	// 4. Calculate and mint this block's provision
	blockProvisionCoin := minter.BlockProvision(mintParams)
	ctx.Logger().Info("💰 Fallback block provision calculated",
		"block_provision", blockProvisionCoin.String())

	if blockProvisionCoin.Amount.IsZero() {
		ctx.Logger().Info("⚠️ Zero block provision in fallback, skipping")
		return nil
	}

	coins := sdk.NewCoins(blockProvisionCoin)
	if err := mintKeeper.MintCoins(ctx, coins); err != nil {
		ctx.Logger().Error("❌ Failed to mint tokens in fallback",
			"error", err.Error())
		return fmt.Errorf("failed to mint tokens: %w", err)
	}

	ctx.Logger().Info("✅ Fallback tokens minted successfully",
		"minted_coins", coins.String())

	// 5. Calculate and log distribution
	distributions, err := splitter.CalculateDistribution(blockProvisionCoin.Amount)
	if err != nil {
		ctx.Logger().Error("❌ Failed to calculate fallback distribution",
			"error", err.Error())
		return fmt.Errorf("failed to calculate distribution: %w", err)
	}

	ctx.Logger().Info("📊 Fallback distribution calculated")
	for i, dist := range distributions {
		ctx.Logger().Info("💸 Fallback distribution detail",
			"index", i,
			"module", dist.ModuleName,
			"amount", dist.Amount.String(),
			"description", dist.Description)
	}

	// 6. Distribute tokens using fallback 50/50 split
	if err := splitter.DistributeTokens(ctx, mintParams.MintDenom, blockProvisionCoin.Amount); err != nil {
		ctx.Logger().Error("❌ Failed to distribute fallback tokens",
			"error", err.Error())
		return fmt.Errorf("failed to distribute tokens: %w", err)
	}

	ctx.Logger().Info("✅ Fallback tokens distributed successfully")

	// 7. Update minter state
	if err := mintKeeper.Minter.Set(ctx, minter); err != nil {
		ctx.Logger().Error("❌ Failed to update minter state in fallback",
			"error", err.Error())
		return fmt.Errorf("failed to update minter state: %w", err)
	}

	ctx.Logger().Info("✅ Fallback minter state updated successfully")

	// 8. Emit fallback events
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeEmissionDistributed,
			sdk.NewAttribute(types.AttributeKeyAmount, blockProvisionCoin.String()),
			sdk.NewAttribute("governance_controlled", "false"),
			sdk.NewAttribute("fallback_active", "true"),
			sdk.NewAttribute("split_type", "50_50_hardcoded"),
			sdk.NewAttribute("num_distributions", fmt.Sprintf("%d", len(distributions))),
		),
	)

	// Also emit standard events for compatibility
	bondedRatio, _ := k.stakingKeeper.BondedRatio(ctx)
	types.EmitEmissionSplitEvent(
		ctx,
		blockProvisionCoin.Amount,
		distributions,
		config,
		minter,
		bondedRatio,
	)

	ctx.Logger().Info("🎉 FALLBACK CONFIGURATION COMPLETED SUCCESSFULLY 🎉",
		"block_height", ctx.BlockHeight(),
		"total_minted", blockProvisionCoin.String(),
		"split_type", "50_50_hardcoded")

	return nil
}

// activateEmergencyFallback activates emergency fallback when critical errors occur
func (k Keeper) activateEmergencyFallback(ctx sdk.Context, mintKeeper *mintkeeper.Keeper, reason string) error {
	ctx.Logger().Error("🚨 ACTIVATING EMERGENCY FALLBACK 🚨",
		"reason", reason,
		"block_height", ctx.BlockHeight())

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

	// Use the same fallback logic as normal fallback to avoid recursion
	return k.applyFallbackConfiguration(ctx, mintKeeper)
}

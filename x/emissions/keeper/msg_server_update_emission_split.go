package keeper

import (
	"bytes"
	"context"
	"fmt"

	errorsmod "cosmossdk.io/errors"
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"scarlett-core/x/emissions/types"
)

func (k msgServer) UpdateEmissionSplit(goCtx context.Context, msg *types.MsgUpdateEmissionSplit) (*types.MsgUpdateEmissionSplitResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// 1. Authority validation - only governance can update emission parameters
	authority, err := k.addressCodec.StringToBytes(msg.Creator)
	if err != nil {
		return nil, errorsmod.Wrap(err, "invalid authority address")
	}

	if !bytes.Equal(k.GetAuthority(), authority) {
		expectedAuthorityStr, _ := k.addressCodec.BytesToString(k.GetAuthority())
		return nil, errorsmod.Wrapf(types.ErrUnauthorized, "invalid authority; expected %s, got %s", expectedAuthorityStr, msg.Creator)
	}

	// 2. Check for emergency stop conditions
	currentParams, err := k.GetEmissionParams(ctx)
	if err == nil && currentParams.IsEmergencyActive() {
		return nil, types.ErrEmergencyActive
	}

	// 3. Parse and validate destinations from governance proposal
	destinations, err := k.parseDestinations(goCtx, msg.Destinations, msg.Weights)
	if err != nil {
		return nil, fmt.Errorf("failed to parse destinations: %w", err)
	}

	// 4. Create new emission parameters
	newParams := types.EmissionParams{
		Destinations:    destinations,
		Enabled:         true,
		UpdatedBy:       fmt.Sprintf("gov_proposal_%s", msg.Creator),
		UpdatedAt:       ctx.BlockHeight(),
		EmergencyStop:   false,
		FallbackEnabled: false,
	}

	// 5. Comprehensive validation with safety bounds
	if err := types.ValidateEmissionParams(newParams); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	// 6. Additional governance-specific validations
	if err := k.validateGovernanceProposal(newParams, msg); err != nil {
		return nil, fmt.Errorf("governance validation failed: %w", err)
	}

	// 7. Store parameter history for audit trail
	if err := k.SetEmissionHistory(ctx, ctx.BlockHeight(), currentParams); err != nil {
		// Log but don't fail - history is important but not critical
		ctx.Logger().Error("Failed to store emission history", "error", err, "height", ctx.BlockHeight())
	}

	// 8. Update emission parameters
	if err := k.SetEmissionParams(ctx, newParams); err != nil {
		return nil, fmt.Errorf("failed to store emission parameters: %w", err)
	}

	// 9. Emit comprehensive governance event
	k.emitEmissionParamsUpdatedEvent(ctx, currentParams, newParams, msg.Reason)

	return &types.MsgUpdateEmissionSplitResponse{}, nil
}

// parseDestinations converts string arrays from governance proposal to EmissionDestination structs
func (k msgServer) parseDestinations(goCtx context.Context, moduleNames, weights []string) ([]types.EmissionDestination, error) {
	if len(moduleNames) != len(weights) {
		return nil, fmt.Errorf("mismatched destinations and weights count: %d vs %d", len(moduleNames), len(weights))
	}

	if len(moduleNames) == 0 {
		return nil, types.ErrNoDestinations
	}

	// Core modules that are exempt from registry requirements
	coreModules := map[string]bool{
		"fee_collector":  true, // Validator and delegator rewards
		"distribution":   true, // Distribution module
		"community_pool": true, // Community pool
	}

	destinations := make([]types.EmissionDestination, len(moduleNames))
	rules := types.DefaultValidationRules()

	for i, moduleName := range moduleNames {
		// 🔒 REGISTRY VALIDATION: Ensure all non-core modules are registered
		if !coreModules[moduleName] {
			if !k.Keeper.IsModuleRegistered(goCtx, moduleName) {
				return nil, types.ErrModuleNotRegistered.Wrapf("destination module '%s' is not registered - modules must be registered before being eligible for emissions", moduleName)
			}
		}

		// Parse weight from string
		weight, err := math.LegacyNewDecFromStr(weights[i])
		if err != nil {
			return nil, fmt.Errorf("invalid weight for %s: %w", moduleName, err)
		}

		// Create destination with safety bounds
		destinations[i] = types.EmissionDestination{
			ModuleName:  moduleName,
			Weight:      weight,
			Description: k.getModuleDescription(moduleName),
			Enabled:     true,
			MinWeight:   k.getMinWeightForModule(moduleName, rules),
			MaxWeight:   rules.MaxSingleDestination,
		}
	}

	return destinations, nil
}

// validateGovernanceProposal performs additional governance-specific validations
func (k msgServer) validateGovernanceProposal(params types.EmissionParams, msg *types.MsgUpdateEmissionSplit) error {
	rules := types.DefaultValidationRules()

	// 1. Check minimum proposal requirements
	if len(msg.Destinations) > int(rules.MaxDestinations) {
		return types.ErrExceedsMaxDestinations
	}

	// 2. Validate proposal reason is provided
	if msg.Reason == "" {
		return fmt.Errorf("governance proposal must include reason for emission changes")
	}

	// 3. Check for dangerous configurations
	if err := k.validateDangerousConfigurations(params); err != nil {
		return fmt.Errorf("dangerous configuration detected: %w", err)
	}

	// 4. Ensure validator rewards are sufficient
	if !params.IsValidatorRewardSufficient() {
		return types.ErrBelowMinValidatorReward
	}

	return nil
}

// validateDangerousConfigurations checks for potentially harmful emission configurations
func (k msgServer) validateDangerousConfigurations(params types.EmissionParams) error {
	rules := types.DefaultValidationRules()

	// 1. Check for concentration risk - no single destination should exceed safety threshold
	for _, dest := range params.Destinations {
		if dest.Weight.GT(rules.MaxSingleDestination) {
			return fmt.Errorf("destination %s exceeds maximum single destination limit: %s > %s",
				dest.ModuleName, dest.Weight.String(), rules.MaxSingleDestination.String())
		}
	}

	// 2. Ensure fee_collector (validator rewards) meets minimum threshold
	validatorWeight := math.LegacyZeroDec()
	for _, dest := range params.Destinations {
		if dest.ModuleName == "fee_collector" && dest.Enabled {
			validatorWeight = dest.Weight
			break
		}
	}

	if validatorWeight.LT(rules.MinValidatorReward) {
		return fmt.Errorf("validator rewards below safety minimum: %s < %s",
			validatorWeight.String(), rules.MinValidatorReward.String())
	}

	return nil
}

// getModuleDescription returns a human-readable description for known modules
func (k msgServer) getModuleDescription(moduleName string) string {
	descriptions := map[string]string{
		"fee_collector":    "Validator and delegator staking rewards",
		"inferencerewards": "AI inference provider rewards",
		"distribution":     "Distribution module rewards",
		"community_pool":   "Community pool funding",
	}

	if desc, exists := descriptions[moduleName]; exists {
		return desc
	}
	return fmt.Sprintf("Rewards for %s module", moduleName)
}

// getMinWeightForModule returns minimum weight requirements for specific modules
func (k msgServer) getMinWeightForModule(moduleName string, rules types.EmissionValidationRules) math.LegacyDec {
	// fee_collector (validators) have special minimum requirements
	if moduleName == "fee_collector" {
		return rules.MinValidatorReward
	}
	// Other modules have no specific minimum (can be zero)
	return math.LegacyZeroDec()
}

// emitEmissionParamsUpdatedEvent emits a comprehensive event for governance parameter updates
func (k msgServer) emitEmissionParamsUpdatedEvent(ctx sdk.Context, oldParams, newParams types.EmissionParams, reason string) {
	attributes := []sdk.Attribute{
		sdk.NewAttribute(types.AttributeKeyOldParams, k.paramsToString(oldParams)),
		sdk.NewAttribute(types.AttributeKeyNewParams, k.paramsToString(newParams)),
		sdk.NewAttribute(types.AttributeKeyReason, reason),
		sdk.NewAttribute("block_height", fmt.Sprintf("%d", ctx.BlockHeight())),
		sdk.NewAttribute("updated_by", newParams.UpdatedBy),
		sdk.NewAttribute("num_destinations", fmt.Sprintf("%d", len(newParams.Destinations))),
		sdk.NewAttribute("governance_controlled", "true"),
	}

	// Add individual destination details
	for i, dest := range newParams.Destinations {
		prefix := fmt.Sprintf("dest_%d", i)
		attributes = append(attributes,
			sdk.NewAttribute(prefix+"_module", dest.ModuleName),
			sdk.NewAttribute(prefix+"_weight", dest.Weight.String()),
			sdk.NewAttribute(prefix+"_enabled", fmt.Sprintf("%t", dest.Enabled)),
			sdk.NewAttribute(prefix+"_description", dest.Description),
		)
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(types.EventTypeEmissionParamsUpdated, attributes...),
	)
}

// paramsToString converts emission parameters to a compact string representation
func (k msgServer) paramsToString(params types.EmissionParams) string {
	if jsonStr, err := params.ToJSON(); err == nil {
		return jsonStr
	}
	return "invalid_params"
}

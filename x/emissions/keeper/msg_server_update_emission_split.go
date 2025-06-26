package keeper

import (
	"context"
	"fmt"
	"strings"

	"scarlett-core/x/emissions/types"

	errorsmod "cosmossdk.io/errors"
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k msgServer) UpdateEmissionSplit(ctx context.Context, msg *types.MsgUpdateEmissionSplit) (*types.MsgUpdateEmissionSplitResponse, error) {
	// Validate authority
	if msg.Creator != string(k.GetAuthority()) {
		return nil, errorsmod.Wrapf(types.ErrUnauthorized, "invalid authority; expected %s, got %s", k.GetAuthority(), msg.Creator)
	}

	// Parse destinations and weights from string arrays
	if len(msg.Destinations) != len(msg.Weights) {
		return nil, errorsmod.Wrap(types.ErrInvalidDestination, "destinations and weights arrays must have same length")
	}

	// Build emission destinations
	destinations := make([]types.EmissionDestination, len(msg.Destinations))
	for i, moduleName := range msg.Destinations {
		weight, err := math.LegacyNewDecFromStr(msg.Weights[i])
		if err != nil {
			return nil, errorsmod.Wrapf(types.ErrInvalidWeight, "invalid weight for destination %d: %v", i, err)
		}

		destinations[i] = types.EmissionDestination{
			ModuleName:  moduleName,
			Weight:      weight,
			Description: fmt.Sprintf("Governance-controlled rewards for %s module", moduleName),
			Enabled:     true,
			MinWeight:   math.LegacyNewDecWithPrec(1, 2),  // 1% minimum
			MaxWeight:   math.LegacyNewDecWithPrec(70, 2), // 70% maximum
		}
	}

	// Get current parameters for history tracking
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	currentParams, _ := k.EmissionParams.Get(ctx)

	// Create new emission parameters
	newParams := types.EmissionParams{
		Destinations: destinations,
		Enabled:      true,
		UpdatedBy:    msg.Reason,
		UpdatedAt:    sdkCtx.BlockHeight(),
	}

	// Validate new parameters
	if err := newParams.Validate(); err != nil {
		return nil, errorsmod.Wrap(types.ErrInvalidDestination, err.Error())
	}

	// Store new parameters
	if err := k.EmissionParams.Set(ctx, newParams); err != nil {
		return nil, errorsmod.Wrap(err, "failed to set emission parameters")
	}

	// Store in history for audit trail
	if err := k.EmissionHistory.Set(ctx, sdkCtx.BlockHeight(), newParams); err != nil {
		return nil, errorsmod.Wrap(err, "failed to store emission history")
	}

	// Emit governance event
	sdkCtx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeEmissionParamsUpdated,
			sdk.NewAttribute(types.AttributeKeyProposalID, msg.Reason),
			sdk.NewAttribute(types.AttributeKeyOldParams, formatParams(currentParams)),
			sdk.NewAttribute(types.AttributeKeyNewParams, formatParams(newParams)),
			sdk.NewAttribute(sdk.AttributeKeyModule, types.ModuleName),
			sdk.NewAttribute(sdk.AttributeKeyAction, "update_emission_split"),
		),
	)

	return &types.MsgUpdateEmissionSplitResponse{}, nil
}

// formatParams formats emission parameters for event attributes
func formatParams(params types.EmissionParams) string {
	if len(params.Destinations) == 0 {
		return "none"
	}

	var parts []string
	for _, dest := range params.Destinations {
		parts = append(parts, fmt.Sprintf("%s:%s", dest.ModuleName, dest.Weight.String()))
	}
	return strings.Join(parts, ",")
}

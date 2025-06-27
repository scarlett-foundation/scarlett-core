package emissions

import (
	"strconv"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"
)

// Event type constants
const (
	EventTypeEmissionSplit = "emission_split"

	AttributeKeyTotalAmount     = "total_amount"
	AttributeKeyNumDestinations = "num_destinations"
	AttributeKeyEnabled         = "enabled"
	AttributeKeyDestination     = "destination"
	AttributeKeyAmount          = "amount"
	AttributeKeyDescription     = "description"
)

// EmitEmissionSplitEvent emits a detailed event for emission splitting
func EmitEmissionSplitEvent(
	ctx sdk.Context,
	totalAmount math.Int,
	distributions []Distribution,
	config EmissionsConfig,
	minter minttypes.Minter,
	bondedRatio math.LegacyDec,
) {
	// Build event attributes
	attributes := []sdk.Attribute{
		sdk.NewAttribute(minttypes.AttributeKeyBondedRatio, bondedRatio.String()),
		sdk.NewAttribute(minttypes.AttributeKeyInflation, minter.Inflation.String()),
		sdk.NewAttribute(minttypes.AttributeKeyAnnualProvisions, minter.AnnualProvisions.String()),
		sdk.NewAttribute(AttributeKeyTotalAmount, totalAmount.String()),
		sdk.NewAttribute(AttributeKeyNumDestinations, strconv.Itoa(len(distributions))),
		sdk.NewAttribute(AttributeKeyEnabled, strconv.FormatBool(config.Enabled)),
	}

	// Add distribution details
	for i, dist := range distributions {
		prefix := AttributeKeyDestination + "_" + strconv.Itoa(i)
		attributes = append(attributes,
			sdk.NewAttribute(prefix+"_module", dist.ModuleName),
			sdk.NewAttribute(prefix+"_amount", dist.Amount.String()),
			sdk.NewAttribute(prefix+"_description", dist.Description),
		)
	}

	// Emit the event
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(EventTypeEmissionSplit, attributes...),
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

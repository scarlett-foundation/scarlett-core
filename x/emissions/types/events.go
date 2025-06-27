package types

// Event type constants for emissions module
const (
	EventTypeEmissionParamsUpdated = "emission_params_updated"
	EventTypeDestinationAdded      = "emission_destination_added"
	EventTypeDestinationRemoved    = "emission_destination_removed"
	EventTypeDestinationToggled    = "emission_destination_toggled"
	EventTypeEmissionDistributed   = "emission_distributed"
	EventTypeEmergencyStop         = "emission_emergency_stop"
)

// Event attribute keys
const (
	AttributeKeyProposalID      = "proposal_id"
	AttributeKeyOldParams       = "old_params"
	AttributeKeyNewParams       = "new_params"
	AttributeKeyDestination     = "destination"
	AttributeKeyAmount          = "amount"
	AttributeKeyReason          = "reason"
	AttributeKeyEmergencyReason = "emergency_reason"
	AttributeKeyEnabled         = "enabled"
	AttributeKeyWeight          = "weight"
	AttributeKeyDescription     = "description"
)

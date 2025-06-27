package types

// Event types for token burning
const (
	EventTypeBurnTokens = "burn_tokens"

	AttributeKeyBurner = "burner"
	AttributeKeyAmount = "amount"
	AttributeKeyDenom  = "denom"
)

// Event types for genesis burn
const (
	EventTypeBurnGenesisStake     = "burn_genesis_stake"
	EventTypeGenesisDecentralized = "genesis_decentralized"

	AttributeKeyGenesisAddress = "genesis_address"
	AttributeKeyUnbondedAmount = "unbonded_amount"
	AttributeKeyBurnedAmount   = "burned_amount"
	AttributeKeyCompletionTime = "completion_time"
	AttributeKeyUnbondingID    = "unbonding_id"
)

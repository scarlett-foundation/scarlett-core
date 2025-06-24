# Genesis Fair Launch Feature - Scarlett Core

## Overview

The Genesis Fair Launch feature implements a **Satoshi-style anonymous decentralization mechanism** where the founding validator (genesis account) can permanently burn their entire stake, removing themselves from the network and reducing the total token supply.

This creates **true decentralization** by:
- ✅ **Eliminating founder control** - Genesis keys become useless after execution
- ✅ **Reducing token supply** - All genesis tokens are burned (not redistributed)
- ✅ **One-time execution** - Cannot be repeated or reversed
- ✅ **Network safety** - Requires minimum validator diversity before execution
- ✅ **Transparent process** - All actions are recorded on-chain with events

## Architecture

### Implementation Pattern
- **Module**: Extends existing `x/scarlettcore` module (same pattern as tokenburner)
- **Message**: `MsgBurnGenesisStake` - Single transaction burns all genesis stake
- **Integration**: Uses existing burn functionality and staking module
- **Automation**: Automatic token burning after 21-day unbonding period

### Key Components

#### 1. Message Handler (`msg_server_burn_genesis_stake.go`)
- Validates caller is configured genesis address
- Prevents multiple executions (one-time only)
- Ensures minimum 3-validator network safety
- Unbonds all delegations from genesis validator
- Emits events for transparency

#### 2. EndBlock Processing (`genesis_burn.go`)
- Monitors completed unbonding delegations
- Automatically burns tokens when unbonding period expires
- Cleans up tracking state after burning

#### 3. Parameter Management
- `genesis_address` parameter configurable via governance
- Empty by default (feature disabled until configured)

## Network Configuration

### Current Test Setup (5 Validators)
```yaml
validators:
- name: alice          # Genesis address (100M sclt)
  bonded: 100000000sclt
- name: validator1     # 200M sclt
  bonded: 200000000sclt  
- name: validator2     # 100M sclt
  bonded: 100000000sclt
- name: validator3     # 100M sclt
  bonded: 100000000sclt
- name: validator4     # 100M sclt
  bonded: 100000000sclt
```

**Total Network Stake**: 600M sclt  
**Alice's Share**: 16.7% (100M / 600M)  
**Post-Burn Supply**: 500M sclt (16.7% reduction)

## Testing Steps

### Phase 1: Network Setup

1. **Clean Start**
   ```bash
   rm -rf ~/.scarlett-core
   ```

2. **Start Network**
   ```bash
   ignite chain serve --verbose
   ```

3. **Verify Validators**
   ```bash
   scarlett-cored query staking validators
   ```
   Should show 1 active validator (Alice).

4. **Verify All Accounts Created**
   ```bash
   scarlett-cored keys list
   ```
   Should show: alice, bob, validator1, validator2, validator3, validator4

### Phase 2: Create Additional Delegations (CRITICAL FOR TESTING)

**⚠️ IMPORTANT**: We need to delegate from other accounts to Alice's validator so that when Alice burns her own stake, the validator remains active with delegated stake from others.

5. **Get Alice's Validator Address**
   ```bash
   scarlett-cored query staking validators --output json | jq '.validators[0].operator_address'
   ```

6. **Create Delegations from Other Accounts**
   ```bash
   # Get Alice's validator operator address first
   ALICE_VALIDATOR=$(scarlett-cored query staking validators --output json | jq -r '.validators[0].operator_address')
   
   # Delegate from validator1 (200M sclt)
   scarlett-cored tx staking delegate $ALICE_VALIDATOR 200000000sclt --from validator1 --gas auto --gas-adjustment 1.5 --yes
   
   # Delegate from validator2 (100M sclt)  
   scarlett-cored tx staking delegate $ALICE_VALIDATOR 100000000sclt --from validator2 --gas auto --gas-adjustment 1.5 --yes
   
   # Delegate from validator3 (100M sclt)
   scarlett-cored tx staking delegate $ALICE_VALIDATOR 100000000sclt --from validator3 --gas auto --gas-adjustment 1.5 --yes
   
   # Delegate from validator4 (100M sclt)
   scarlett-cored tx staking delegate $ALICE_VALIDATOR 100000000sclt --from validator4 --gas auto --gas-adjustment 1.5 --yes
   ```

7. **Verify Total Delegations**
   ```bash
   scarlett-cored query staking delegations-to $ALICE_VALIDATOR --output json | jq '.delegation_responses[] | {delegator: .delegation.delegator_address, shares: .delegation.shares, balance: .balance.amount}'
   ```
   Should show:
   - Alice: 100M sclt (self-delegation)
   - validator1: 200M sclt
   - validator2: 100M sclt  
   - validator3: 100M sclt
   - validator4: 100M sclt
   - **Total: 600M sclt**

### Phase 3: Configure Genesis Address

8. **Update Proposal with Correct Address**
   ```bash
   # Get Alice's current address
   ALICE_ADDRESS=$(scarlett-cored keys show alice --address)
   echo "Alice's address: $ALICE_ADDRESS"
   
   # Update proposal.json with correct address
   ```

9. **Create Governance Proposal**
   ```bash
   scarlett-cored tx gov submit-proposal proposal.json --from alice --gas auto --gas-adjustment 1.5 --yes
   ```

10. **Vote on Proposal**
    ```bash
    scarlett-cored tx gov vote 1 yes --from alice --gas auto --gas-adjustment 1.5 --yes
    ```

11. **Wait for Proposal to Pass**
    ```bash
    # Wait for voting period (60 seconds)
    sleep 65
    
    # Check if genesis address is configured
    scarlett-cored query scarlettcore params --output json
    ```

### Phase 4: Execute Genesis Fair Launch

12. **Verify Pre-Burn State**
    ```bash
    # Check Alice's validator total stake (should be 600M)
    scarlett-cored query staking validators --output json | jq '.validators[0] | {moniker: .description.moniker, tokens: .tokens}'
    
    # Check Alice's self-delegation (should be 100M)
    scarlett-cored query staking delegations alice --output json | jq '.delegation_responses[] | {validator: .delegation.validator_address, shares: .delegation.shares, balance: .balance.amount}'
    ```

13. **Execute Genesis Burn**
    ```bash
    scarlett-cored tx scarlettcore burn-genesis-stake --from alice --gas auto --gas-adjustment 1.5 --yes
    ```

### Phase 5: Verify Results

14. **Check Transaction Success**
    ```bash
    # Look for debug logs in chain output:
    # - "DEBUG: Successfully unbonded Alice's self-delegation"
    # - Should show unbonded=100000000
    # - Should show completion_time for 21-day unbonding
    ```

15. **Verify Network Continues**
    ```bash
    # Network should continue running (not crash)
    # Alice's validator should still be active with reduced stake
    scarlett-cored query staking validators --output json | jq '.validators[0] | {moniker: .description.moniker, tokens: .tokens}'
    ```
    Should show Alice's validator with **500M sclt** (600M - 100M burned)

16. **Check Alice's Delegations**
    ```bash
    # Alice should have NO delegations (her self-delegation was unbonded)
    scarlett-cored query staking delegations alice --output json
    
    # But other delegators should still be delegated to Alice's validator
    scarlett-cored query staking delegations-to $ALICE_VALIDATOR --output json | jq '.delegation_responses[] | {delegator: .delegation.delegator_address, balance: .balance.amount}'
    ```

17. **Check Unbonding Delegation**
    ```bash
    scarlett-cored query staking unbonding-delegations alice --output json
    ```
    Should show Alice's 100M sclt unbonding with completion time ~21 days from now.

18. **Monitor Token Burning** (After 21 days)
    ```bash
    # Check total supply reduction
    scarlett-cored query bank total --output json
    ```

## Expected Results

### Immediate Effects (Block of Execution)
- ✅ **Alice's self-delegation unbonded**: 100M sclt enters 21-day unbonding
- ✅ **Network continues**: Alice's validator remains active with 500M delegated stake
- ✅ **Other delegators unaffected**: validator1-4 keep their delegations to Alice
- ✅ **Alice loses validator rewards**: No longer earns from her own stake
- ✅ **`genesis_stake_burn_initiated` event emitted**
- ✅ **One-time execution flag set** (prevents re-execution)

### Long-term Effects (After 21 Days)
- ✅ **100M sclt automatically burned** via EndBlock processing
- ✅ **Total supply reduced** from 600M to 500M sclt (16.7% deflation)
- ✅ **Alice's keys become useless** for validator rewards (true decentralization)
- ✅ **Network remains stable** with Alice's validator operated by delegated stake
- ✅ **`genesis_tokens_burned` event emitted**

## Key Testing Insights

### Why Delegations Are Critical
- **Without delegations**: Alice unbonds → Validator removed → Network crashes
- **With delegations**: Alice unbonds → Validator continues → Network stable
- **Result**: Only Alice's own stake is burned, others are unaffected

### Network State After Burn
```
Before: Alice's Validator = 600M sclt (100M self + 500M delegated)
After:  Alice's Validator = 500M sclt (0M self + 500M delegated)
Burned: 100M sclt (Alice's self-delegation only)
```

### Validator Ownership vs Stake
- **Alice remains validator operator** (can still propose blocks)
- **Alice loses staking rewards** from her own stake
- **Delegators continue earning** rewards from their stake
- **True decentralization**: Alice's financial incentive removed

## Safety Features

### Pre-execution Validation
- **Genesis Address Check**: Only configured genesis address can execute
- **One-time Execution**: Cannot be called multiple times
- **Validator Diversity**: Requires minimum 3 active validators
- **Stake Verification**: Validates genesis address has delegations to burn

### Network Protection
- **Gradual Process**: 21-day unbonding period prevents sudden supply shock
- **Automatic Execution**: No manual intervention required for token burning
- **Event Tracking**: All actions recorded on-chain for transparency
- **Governance Control**: Genesis address configurable via governance proposals

## Error Scenarios

### Common Error Messages
```bash
# Genesis address not configured
"genesis address not configured"

# Wrong caller
"caller is not the genesis address"

# Already executed
"genesis stake has already been burned"

# Insufficient validators
"insufficient active validators for safe execution"

# No stake to burn
"no delegations found for genesis address"
```

### Troubleshooting
1. **Proposal Rejected**: Check voting period and quorum requirements
2. **Transaction Failed**: Verify gas limits and account balances
3. **Wrong Address**: Ensure using exact genesis address from params
4. **Network Issues**: Check validator connectivity and consensus

## Development Notes

### File Structure
```
x/scarlettcore/
├── keeper/
│   ├── msg_server_burn_genesis_stake.go  # Main handler logic
│   └── genesis_burn.go                   # EndBlock processing
├── types/
│   ├── message_burn_genesis_stake.go     # Message definition
│   ├── params.go                         # Parameter definitions
│   └── events.go                         # Event constants
└── proto/scarlettcore/v1/
    ├── tx.proto                          # Message proto
    └── params.proto                      # Parameter proto
```

### Key Dependencies
- **StakingKeeper**: For unbonding delegations
- **BankKeeper**: For token burning (via existing scarlettcore burn)
- **Collections**: For state management and tracking

### Testing Commands
```bash
# Build and test
make build
make test

# Generate proto files
make proto-gen

# Run local chain
ignite chain serve --verbose

# Clean restart
rm -rf ~/.scarlett-core && ignite chain serve --verbose
```

## Future Enhancements

### Potential Improvements
- **Multiple Genesis Addresses**: Support for multiple founders
- **Partial Burning**: Burn percentage instead of all stake
- **Custom Schedules**: Configurable unbonding periods
- **Advanced Analytics**: Progress tracking and statistics
- **Governance Integration**: Enhanced proposal templates

### Production Considerations
- **Mainnet Testing**: Extensive testnet validation before mainnet
- **Emergency Procedures**: Governance halt mechanisms
- **Documentation**: User guides and validator instructions
- **Monitoring**: Real-time tracking and alerting systems

---

## Summary

The Genesis Fair Launch feature provides a **transparent, irreversible decentralization mechanism** that eliminates founder control while reducing token supply. The implementation follows Cosmos SDK best practices and integrates seamlessly with existing staking and governance modules.

**Key Achievement**: True Satoshi-style anonymity where genesis keys become permanently useless after execution, creating genuine decentralization with provable token scarcity. 
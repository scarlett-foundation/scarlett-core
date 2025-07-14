---
description: 
globs: 
alwaysApply: true
---
# Proofofdegen Merkle Tree Upgrade

## Overview
Add Merkle tree support to the **existing** proofofdegen module for scalable Genesis campaign airdrops. The core patience game mechanics and claim system already work - we just need to add **internal** proof validation and aggregate tracking for 10k+ wallet support.

**Current State**: ✅ Working module with claim mechanics, patience game, and individual wallet tracking  
**Target State**: ✅ Same functionality + internal Merkle proof validation + aggregate tracking for scale

**Key Design Principle**: Users should never see or handle Merkle proofs. The CLI interface remains identical.

## ⚠️ WHAT'S ALREADY IMPLEMENTED
- ✅ Full proofofdegen module with Campaign and EligibleWallet storage
- ✅ MsgClaim with self-claiming validation and patience mechanics  
- ✅ CalculateShare function with concentration effects
- ✅ EligibleAmount and CampaignInfo query endpoints
- ✅ Genesis state with test eligible wallets
- ✅ CLI commands and autocli integration
- ✅ Module registration and EndBlocker functionality

## 🎯 WHAT NEEDS TO BE ADDED

### Step 1: Add Merkle Fields to Proto
**Objective**: Extend existing structures with Merkle tree support (internal only)
**Action**: Add fields to existing proto files:
- `merkle_root` (string) to Campaign message  
- `claimed_count` and `claimed_weight` (uint64) to Campaign for aggregates
- **NO changes to MsgClaim** - keep current interface
**Files to modify**:
- `proto/scarlettcore/proofofdegen/v1/campaign.proto`
**Verification**: User runs `ignite generate proto-go` and builds successfully

### Step 2: Add Merkle Proof Storage
**Objective**: Store proofs in keeper for internal lookup
**Action**: Add proof storage to keeper:
- Add `MerkleProofs collections.Map[string, []string]` to keeper (address → proof array)
- Add proof lookup helper functions
- Keep existing MsgClaim interface unchanged
**Files to modify**:
- `x/proofofdegen/keeper/keeper.go`
- `x/proofofdegen/types/keys.go` (add MerkleProofsKey)
**Verification**: User builds and confirms new storage compiles

### Step 3: Add Internal Proof Validation
**Objective**: Implement automatic proof validation in claim handler
**Action**: Update claim handler to validate proofs internally:
- Look up proof automatically by address from MerkleProofs storage
- Hash leaf (address + weight) 
- Verify proof against campaign merkle_root
- Get weight from proof validation
- Continue with existing claim logic
**Files to modify**:
- `x/proofofdegen/keeper/msg_server_claim.go`
- `x/proofofdegen/types/errors.go` (add ErrInvalidProof)
- Add `keeper/merkle_helpers.go` for proof validation
**Verification**: User builds and tests proof validation with dummy data

### Step 4: Add Aggregate Tracking
**Objective**: Replace individual wallet iteration with aggregate counters
**Action**: Update claim handler and queries to use aggregate tracking:
- Update `claimed_count` and `claimed_weight` in Campaign during claims
- Modify `CalculateShare` to use aggregates instead of wallet iteration
- Update `CampaignInfo` query to use aggregate stats
**Files to modify**:
- `x/proofofdegen/keeper/msg_server_claim.go`
- `x/proofofdegen/keeper/query_campaign_info.go`
- Add `keeper/aggregate_helpers.go` for aggregate calculations
**Verification**: User builds and confirms aggregate calculations work

### Step 5: Create Merkle Genesis Tool
**Objective**: Build standalone tool for CSV-to-Merkle conversion with proof storage
**Action**: Create `cmd/merkle-genesis-tool/` with:
- Read CSV file (address, weight)
- Generate Merkle tree using standard library
- Output genesis JSON with merkle_root, aggregate totals, AND proof storage
- Store proofs in genesis state for keeper loading
**Files to create**:
- `cmd/merkle-genesis-tool/main.go`
- `cmd/merkle-genesis-tool/merkle.go`
**Verification**: User generates genesis from CSV and confirms tool works

### Step 6: Update Genesis State
**Objective**: Switch to Merkle-based genesis initialization with proof storage
**Action**: Update genesis to use merkle_root and proof storage:
- Modify `types/genesis.go` to include merkle_root and proof maps
- Update genesis validation to check merkle_root format
- Load both merkle_root and MerkleProofs into keeper
- Remove individual EligibleWallet entries (replaced by proofs)
**Files to modify**:
- `x/proofofdegen/types/genesis.go`
- `x/proofofdegen/keeper/genesis.go`
**Verification**: User builds with Merkle genesis and confirms it loads

## Implementation Notes

# Users still just do this (no change):
```
scarlett-cored tx proofofdegen claim [my-address]
scarlett-cored query proofofdegen eligible-amount [my-address]
```

### Existing Architecture to Preserve
- **Patience Game**: Keep exact same CalculateShare logic, just use aggregates
- **Self-Claiming**: Keep existing creator == address validation
- **Error Handling**: Keep existing error types, add ErrInvalidProof
- **Query Endpoints**: Keep same responses, just calculate from aggregates
- **CLI Commands**: Keep existing claim command interface IDENTICAL

### New Internal Merkle Claim Flow
1. User submits existing MsgClaim with just address (no proof needed)
2. Keeper automatically looks up proof by address from MerkleProofs storage
3. Keeper validates proof against campaign merkle_root internally
4. If valid, existing claim logic continues unchanged
5. Update campaign aggregates instead of individual wallet tracking

### User Experience (Unchanged)
```bash
# Users still just do this:
scarlett-cored tx proofofdegen claim [my-address]
scarlett-cored query proofofdegen eligible-amount [my-address]
```

### Build Verification Protocol
**After Each Step**:
1. Agent completes code changes
2. Agent asks: "Please run `ignite chain serve --reset-once --verbose` and confirm it builds successfully"
3. User confirms build status
4. If failure: Agent fixes before proceeding
5. If success: Agent commits and proceeds

## Success Criteria

### ✅ Functionality Preserved
- All existing claim mechanics work identically
- Patience game rewards still increase over time
- Self-claiming validation still enforced
- Query endpoints return same information
- **CLI interface completely unchanged**

### ✅ Scalability Added
- Internal Merkle proof validation for unlimited wallets
- Aggregate tracking eliminates O(n) operations
- Genesis state more efficient than individual wallet storage
- Proof validation is O(log n) per claim
- **Zero user-facing complexity**

### ✅ Tooling Complete
- merkle-genesis-tool generates production genesis with internal proofs
- Existing CLI commands work identically
- Clear migration path from current genesis
- **No additional user tools needed**

## Auto-Attachment
This rule attaches for:
- Proofofdegen module development
- Merkle tree implementation (internal)
- Airdrop scalability improvements
- Proof validation systems
- Aggregate tracking optimization 
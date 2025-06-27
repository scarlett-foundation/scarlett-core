# Open Module Registration & Funding - Scarlett Core

## Overview

The Open Module Registration feature transforms Scarlett Core into a **permissionless platform**, allowing any developer to build and register new modules on-chain. These registered modules can then be **funded with a share of token emissions through on-chain governance**.

This creates a powerful, decentralized flywheel for ecosystem growth:
- ✅ **Permissionless Innovation**: Any developer can build and register a module.
- ✅ **Community-Driven Funding**: Token holders have the final say on which modules receive token rewards.
- ✅ **Dynamic Economic Incentives**: The network can adapt to fund the most valuable new applications and services.
- ✅ **Transparent & Secure**: A two-step process (registration then governance) ensures both openness and security.
- ✅ **Discoverable Ecosystem**: Provides an on-chain registry of all community-built modules available for funding.

## Architecture & Workflow

The system introduces a two-phase process managed by the `x/emissions` module: a developer-led **Registration Phase** and a community-led **Activation Phase**.

### Implementation Pattern
- **Module**: Extends the existing `x/emissions` module.
- **New Message**: `MsgRegisterModule` - Allows developers to add their module to an on-chain registry.
- **Modified Message**: `MsgUpdateEmissionSplit` - The existing governance message is enhanced to validate proposals against the new registry.

### Key Components

#### 1. On-Chain Module Registry (`x/emissions/keeper`)
- A new `collections.Map` stores `RegisteredModule` objects.
- This registry acts as the single source of truth for all modules eligible for funding.

#### 2. Registration Message Handler (`msg_server_register_module.go`)
- Validates that the caller has provided all required information (module name, description).
- **Verifies that the module account exists on-chain** to prevent funding non-existent modules.
- Adds the new module to the registry with a `REGISTERED` status.

#### 3. Governance Validation Logic (`msg_server_update_emission_split.go`)
- When a governance proposal is submitted to update emission splits, it now checks each destination module against the registry.
- If a module in the proposal is not in the registry, the proposal is rejected.

---

## User & Developer Workflow

### Phase 1: Developer Registers a New Module

A developer builds a new Cosmos SDK module and wants it to be eligible for rewards.

1. **Build a Module**
   ```bash
   # Developer uses Ignite CLI to build their module
   ignite scaffold module mynewfeature --dep bank
   # ...developer adds logic to their module...
   ```

2. **Verify Module Account Creation**
   After the chain is running with the new module, the developer confirms the module's account was created.
   ```bash
   scarlett-cored q auth module-account mynewfeature
   ```

3. **Register the Module On-Chain**
   The developer submits a transaction to add their module to the official registry.
   ```bash
   # The creator is the developer's address
   scarlett-cored tx emissions register-module "mynewfeature" "This module rewards users for contributing AI training data." --from <developer-address> --gas auto --gas-adjustment 1.5 --yes
   ```

4. **Verify Registration**
   Anyone can now see the new module in the list of registered candidates.
   ```bash
   scarlett-cored query emissions list-registered-modules
   ```
   The output should now include `mynewfeature`.

### Phase 2: Community Activates the Module via Governance

Once a module is registered, the community can choose to fund it.

5. **Create a Governance Proposal**
   A community member creates a proposal to add the new module to the emissions split.
   ```json
   // file: proposal-fund-mynewfeature.json
   {
     "title": "Fund 'mynewfeature' Module with 10% of Emissions",
     "description": "This proposal activates the new AI training data module...",
     "messages": [
       {
         "@type": "/scarlettcore.emissions.v1.MsgUpdateEmissionSplit",
         "authority": "scarlett10d07y265gmmuvt4z0w9aw880jnsr700j4l5sjv",
         "destinations": [
           { "module_name": "fee_collector", "weight": "0.50" },
           { "module_name": "inferencerewards", "weight": "0.40" },
           { "module_name": "mynewfeature", "weight": "0.10" }
         ]
       }
     ],
     "deposit": "10000000sclt"
   }
   ```

6. **Submit and Pass the Proposal**
   ```bash
   # Submit the proposal
   scarlett-cored tx gov submit-proposal proposal-fund-mynewfeature.json --from <voter-address> --yes
   
   # Vote on the proposal
   scarlett-cored tx gov vote <proposal-id> yes --from <voter-address> --yes
   ```

### Phase 3: Verify Results

7. **Check Emission Parameters**
   After the proposal passes, the new parameters will be active on-chain.
   ```bash
   scarlett-cored query emissions params -o json | jq
   ```
   The output should now show `mynewfeature` as a destination with a weight of `0.10`.

8. **Monitor Module Account Balance**
   The `mynewfeature` module account will now start receiving 10% of all newly minted tokens every block.
   ```bash
   # Get module address
   MYNEWFEATURE_ADDR=$(scarlett-cored q auth module-account mynewfeature -o json | jq -r .value.address)
   
   # Check balance over time
   scarlett-cored q bank balances $MYNEWFEATURE_ADDR
   ```

## Safety Features

- **Module Account Verification**: The system prevents registration of modules whose accounts do not exist on the chain, preventing funds from being sent to a black hole.
- **Governance Gatekeeping**: Registration is permissionless, but funding is not. Only a successful governance vote can activate emissions for a module, ensuring community oversight.
- **Spam Prevention**: An optional registration fee can be added to the `MsgRegisterModule` message to prevent spamming the registry.
- **Immutability**: Once registered, a module's name cannot be changed, ensuring stability for governance proposals.

## Error Scenarios

### Common Error Messages
```bash
# Registering a module that is already registered
"module is already registered"

# Registering a module whose account doesn't exist on-chain
"module account '...' does not exist"

# Submitting a governance proposal with an unregistered module
"destination module '...' is not registered"
```

## Development Notes

### File Structure
```
x/emissions/
├── keeper/
│   ├── msg_server_register_module.go     # New handler for registration
│   └── msg_server_update_emission_split.go # Modified for validation
├── types/
│   ├── message_register_module.go        # New message definition
│   └── types.go                          # Definition for RegisteredModule
└── proto/scarlettcore/emissions/v1/
    └── tx.proto                          # New RPC for MsgRegisterModule
```

## Future Enhancements

- **Module Reputation**: A system to track a module's uptime or other performance metrics.
- **Automated Funding Tiers**: Governance could set funding tiers (e.g., "Gold", "Silver", "Bronze") that come with predefined emission shares.
- **Dynamic Fees**: Registration fees could be dynamic based on network usage.

---

## Summary

The Open Module Registration feature provides a **clear, secure, and transparent pathway for ecosystem development**. By separating the act of registration from the act of funding, it empowers developers to build permissionlessly while ensuring the community retains ultimate control over the network's economic incentives. This system is the foundation for a truly dynamic and community-owned AI application platform. 
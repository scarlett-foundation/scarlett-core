# CosmWasm Integration & Stake-to-Register - Scarlett Core

## Overview

This document outlines the integration of **CosmWasm**, a secure and interoperable smart contracting platform, into Scarlett Core. This feature elevates Scarlett from a purpose-built blockchain into a **fully permissionless smart contract platform**, enabling any developer to deploy their own applications.

To ensure long-term alignment and prevent spam, this system introduces a **Stake-to-Register** model:
- ✅ **Permissionless Deployment**: Any developer can deploy a CosmWasm smart contract at any time without asking for permission.
- ✅ **Stake-to-Register for Rewards**: To make their contract eligible for emission rewards, developers must **stake `sclt` tokens**, creating an economic bond to the network's success.
- ✅ **Community-Led Funding**: Staking only makes a contract *eligible*. The community retains ultimate control via governance to decide which contracts receive funding.
- ✅ **Enhanced Security**: Contracts run in a secure, sandboxed WebAssembly (Wasm) VM, preventing them from affecting the core protocol.

## Architecture & Workflow

The integration extends the existing `x/emissions` module and introduces the `x/wasmd` module to the runtime.

### Implementation Pattern
- **New Core Module**: The `x/wasmd` module is added to the application to enable Wasm smart contract functionality.
- **Modified Emissions Module**: The `x/emissions` registry is enhanced to support CosmWasm contract addresses as potential reward destinations.
- **New Staking Requirement**: The `MsgRegisterModule` handler will be updated to enforce the `sclt` staking requirement for contract registrations.

### Key Components

#### 1. CosmWasm Module (`x/wasmd`)
- The core engine that handles uploading, instantiating, and executing Wasm smart contracts.
- Provides a secure sandbox, gas metering, and a clean interface for contracts to interact with the underlying chain state.

#### 2. Enhanced Emissions Registry (`x/emissions/keeper`)
- The `RegisteredModule` type will be updated to include a `DestinationType` (`NATIVE_MODULE` or `COSMWASM_CONTRACT`) and a `ContractAddress`.
- The keeper can now store and differentiate between core protocol modules and user-deployed smart contracts.

#### 3. Stake-Aware Registration Handler (`msg_server_register_module.go`)
- The `RegisterModule` message handler will now perform a critical new check for `COSMWASM_CONTRACT` types.
- It will query the `staking` module to **verify that the registering developer has a sufficient amount of `sclt` tokens staked**.
- If the developer's stake is below the required threshold, the registration is rejected.

---

## User & Developer Workflow

### Phase 1: Developer Deploys a Smart Contract

1.  **Write and Compile a Contract**
    A developer writes their application logic in Rust and compiles it to a Wasm binary.
    ```rust
    // Example Rust contract structure
    #[entry_point]
    pub fn instantiate(...) -> Result<Response, ContractError> { ... }

    #[entry_point]
    pub fn execute(...) -> Result<Response, ContractError> { ... }
    ```

2.  **Deploy the Contract to the Chain**
    The developer uploads the Wasm binary and instantiates it on Scarlett Core. This is **fully permissionless**.
    ```bash
    # Upload the contract code
    scarlett-cored tx wasm store my_contract.wasm --from <dev-address> --gas auto --yes

    # Instantiate the contract to get a contract address
    scarlett-cored tx wasm instantiate <code-id> '{...}' --label "My AI App" --admin <dev-address> --from <dev-address> --yes
    ```
    This process returns a permanent contract address (e.g., `scarlett1...`).

### Phase 2: Stake & Register for Rewards

3.  **Stake Scarlett Tokens**
    To become eligible for rewards, the developer must first stake the required amount of `sclt` tokens.
    ```bash
    # Stake the required minimum (e.g., 10,000 sclt)
    scarlett-cored tx staking delegate <validator-address> 10000000sclt --from <dev-address> --yes
    ```

4.  **Register the Contract**
    With sufficient stake, the developer can now register their contract address in the emissions registry.
    ```bash
    # Register the deployed contract, linking it to their stake
    scarlett-cored tx emissions register-module "My AI App" "A contract that..." --contract-address <contract-address> --from <dev-address> --yes
    ```

### Phase 3: Community Activates Funding

5.  **Community Governance**
    The community can now see "My AI App" in the list of eligible destinations. A token holder can create a governance proposal to direct a portion of emissions to the contract's address.

6.  **Receive Funding**
    If the proposal passes, the developer's smart contract begins receiving a steady stream of `sclt` tokens every block, which it can use to fund its operations.

## Safety & Economic Security

- **Sandboxed Execution**: The Wasm VM ensures that even buggy or malicious contract code cannot halt the chain or access unauthorized parts of state.
- **Economic Stake**: The staking requirement creates a significant economic disincentive against deploying spam or harmful contracts. A developer who harms the network risks losing the value of their staked tokens.
- **Community Veto**: The community retains the ultimate power to decide what gets funded. A registered contract has no guarantee of receiving rewards; it must prove its value to the token holders.

## Development Notes

### File Structure Changes
```
app/
└── app.go                              # Add x/wasmd to the module manager
x/emissions/
├── keeper/
│   └── msg_server_register_module.go     # Add staking verification logic
└── types/
    └── types.go                          # Add DestinationType and ContractAddress
```

### Key Dependencies
- **WasmVM**: The binary must be built with the correct Wasm library linkage.
- **StakingKeeper**: The `emissions` keeper now requires a dependency on the `staking` keeper to verify developer stakes.

## Future Enhancements

- **Reputation-Based Staking**: The required stake could decrease as a developer's contracts prove to be reliable and valuable over time.
- **Contract-Funded Staking**: Allow contracts to stake tokens on behalf of their users to collectively meet the registration threshold.
- **Automated Slashing**: If a registered contract is found to be malicious by a governance vote, a portion of the developer's stake could be automatically slashed.

---

## Summary

By integrating CosmWasm with a novel **Stake-to-Register** mechanism, Scarlett Core can evolve into a premier platform for decentralized AI development. This model creates a robust, symbiotic relationship: developers are empowered to innovate permissionlessly, while the network is protected and its resources are directed by a community of economically-aligned token holders. 
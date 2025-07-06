# Contracts Module

The `contracts` module provides full CosmWasm smart contract functionality for the Scarlett blockchain, enabling deployment, instantiation, execution, and querying of Rust-based smart contracts.

## Overview

This module integrates [wasmd v0.61.0](https://github.com/CosmWasm/wasmd) to provide:

- **Contract Storage**: Upload WASM bytecode to the blockchain
- **Contract Instantiation**: Create contract instances with initial state
- **Contract Execution**: Call contract functions and modify state
- **Contract Queries**: Read contract state and metadata
- **Full CosmWasm Compatibility**: Support for all CosmWasm v1.x features

## Architecture

The module wraps wasmd's functionality with a custom keeper that:
- Handles VM initialization and lock file management
- Provides seamless integration with the Scarlett blockchain
- Maintains compatibility with the Cosmos SDK ecosystem

## Commands

### Store Contract Code

Upload WASM bytecode to the blockchain and get a code ID for instantiation.

```bash
# Store a contract from a WASM file
scarlett-cored tx contracts store-code <path-to-wasm-file> \
  --from <key-name> \
  --chain-id scarlettcore \
  --keyring-backend test \
  --gas auto \
  --gas-adjustment 1.3 \
  --yes

# Example with the counter contract
scarlett-cored tx contracts store-code examples/contracts/counter/target/wasm32-unknown-unknown/release/counter.wasm \
  --from alice \
  --chain-id scarlettcore \
  --keyring-backend test \
  --gas auto \
  --gas-adjustment 1.3 \
  --yes
```

**Response**: Returns a transaction with `code_id` in the events (e.g., `code_id: "1"`).

### Instantiate Contract

Create a contract instance from stored code with initial state.

```bash
# Instantiate a contract
scarlett-cored tx contracts instantiate-contract <code-id> <init-msg> <label> <funds> <admin> \
  --from <key-name> \
  --chain-id scarlettcore \
  --keyring-backend test \
  --gas auto \
  --gas-adjustment 1.3 \
  --yes

# Example with counter contract (using JSON file for complex messages)
echo '{"count":42}' > /tmp/init.json
scarlett-cored tx contracts instantiate-contract 1 /tmp/init.json "my-counter" "" "" \
  --from alice \
  --chain-id scarlettcore \
  --keyring-backend test \
  --gas auto \
  --gas-adjustment 1.3 \
  --yes
```

**Parameters**:
- `code-id`: The code ID from contract storage
- `init-msg`: JSON initialization message (can be file path or inline JSON)
- `label`: Human-readable label for the contract instance
- `funds`: Coins to send to the contract during instantiation (use `""` for none)
- `admin`: Admin address for contract upgrades (use `""` for no admin)

**Response**: Returns contract address in events (e.g., `_contract_address: "scarlett14hj2tavq8fpesdwxxcu44rty3hh90vhujrvcmstl4zr3txmfvw9s5rjw5p"`).

### Execute Contract

Call a contract function to modify its state.

```bash
# Execute a contract function
scarlett-cored tx contracts execute-contract <contract-address> <execute-msg> <funds> \
  --from <key-name> \
  --chain-id scarlettcore \
  --keyring-backend test \
  --gas auto \
  --gas-adjustment 1.3 \
  --yes

# Example: increment the counter
echo '{"increment":{}}' > /tmp/execute.json
scarlett-cored tx contracts execute-contract scarlett14hj2tavq8fpesdwxxcu44rty3hh90vhujrvcmstl4zr3txmfvw9s5rjw5p /tmp/execute.json "" \
  --from alice \
  --chain-id scarlettcore \
  --keyring-backend test \
  --gas auto \
  --gas-adjustment 1.3 \
  --yes

# Example: reset the counter to a specific value
echo '{"reset":{"count":100}}' > /tmp/reset.json
scarlett-cored tx contracts execute-contract scarlett14hj2tavq8fpesdwxxcu44rty3hh90vhujrvcmstl4zr3txmfvw9s5rjw5p /tmp/reset.json "" \
  --from alice \
  --chain-id scarlettcore \
  --keyring-backend test \
  --gas auto \
  --gas-adjustment 1.3 \
  --yes
```

**Parameters**:
- `contract-address`: The contract address from instantiation
- `execute-msg`: JSON message defining the function to call (can be file path or inline JSON)
- `funds`: Coins to send to the contract during execution (use `""` for none)

### Query Contract State

Read contract state and metadata without modifying it.

```bash
# Query contract state
scarlett-cored query contracts smart-contract-state <contract-address> <query-msg>

# Example: get current counter value
echo '{"get_count":{}}' > /tmp/query.json
scarlett-cored query contracts smart-contract-state scarlett14hj2tavq8fpesdwxxcu44rty3hh90vhujrvcmstl4zr3txmfvw9s5rjw5p /tmp/query.json

# The response will be base64 encoded, decode it:
# Response: data: eyJjb3VudCI6NDN9
echo "eyJjb3VudCI6NDN9" | base64 -d
# Output: {"count":43}
```

### Query Contract Info

Get metadata about a contract instance.

```bash
# Query contract information
scarlett-cored query contracts contract-info <contract-address>

# Example
scarlett-cored query contracts contract-info scarlett14hj2tavq8fpesdwxxcu44rty3hh90vhujrvcmstl4zr3txmfvw9s5rjw5p
```

**Response includes**:
- Contract address
- Code ID
- Creator address
- Admin address (if set)
- Label
- Creation block height and transaction index

## Complete Workflow Example

Here's a complete example using the counter contract:

### 1. Build the Contract (if needed)

```bash
cd examples/contracts/counter
cargo wasm
wasm-opt -Oz target/wasm32-unknown-unknown/release/counter.wasm -o target/wasm32-unknown-unknown/release/counter.wasm
cd ../../../
```

### 2. Store the Contract

```bash
scarlett-cored tx contracts store-code examples/contracts/counter/target/wasm32-unknown-unknown/release/counter.wasm \
  --from alice \
  --chain-id scarlettcore \
  --keyring-backend test \
  --gas auto \
  --gas-adjustment 1.3 \
  --yes
```

**Note the `code_id` from the transaction events (e.g., `"1"`).**

### 3. Instantiate the Contract

```bash
echo '{"count":42}' > /tmp/init.json
scarlett-cored tx contracts instantiate-contract 1 /tmp/init.json "my-counter" "" "" \
  --from alice \
  --chain-id scarlettcore \
  --keyring-backend test \
  --gas auto \
  --gas-adjustment 1.3 \
  --yes
```

**Note the contract address from the transaction events.**

### 4. Query Initial State

```bash
echo '{"get_count":{}}' > /tmp/query.json
scarlett-cored query contracts smart-contract-state <contract-address> /tmp/query.json
```

### 5. Execute Functions

```bash
# Increment the counter
echo '{"increment":{}}' > /tmp/increment.json
scarlett-cored tx contracts execute-contract <contract-address> /tmp/increment.json "" \
  --from alice \
  --chain-id scarlettcore \
  --keyring-backend test \
  --gas auto \
  --gas-adjustment 1.3 \
  --yes

# Reset to a specific value
echo '{"reset":{"count":100}}' > /tmp/reset.json
scarlett-cored tx contracts execute-contract <contract-address> /tmp/reset.json "" \
  --from alice \
  --chain-id scarlettcore \
  --keyring-backend test \
  --gas auto \
  --gas-adjustment 1.3 \
  --yes
```

### 6. Verify State Changes

```bash
scarlett-cored query contracts smart-contract-state <contract-address> /tmp/query.json
```

## Transaction Verification

After submitting any transaction, you can verify it was processed successfully:

```bash
# Query transaction by hash
scarlett-cored query tx <transaction-hash>

# Look for:
# - code: 0 (success)
# - events with relevant contract information
# - gas_used vs gas_wanted
```

## Error Handling

Common errors and solutions:

### Gas Estimation Issues
```bash
# If gas estimation fails, set manual gas limit
--gas 5000000
```

### JSON Parsing Issues
```bash
# Use files for complex JSON to avoid shell escaping issues
echo '{"complex":"message"}' > /tmp/msg.json
# Then use /tmp/msg.json in commands
```

### Contract Not Found
```bash
# Verify contract address is correct
scarlett-cored query contracts contract-info <address>
```

## Module Configuration

The contracts module is configured with:

- **Permissionless Code Upload**: Anyone can store contract code
- **Permissionless Instantiation**: Anyone can instantiate stored contracts
- **Full CosmWasm Feature Set**: Support for all wasmd v0.61.0 features
- **Automatic VM Management**: Lock files and VM instances are managed automatically

## Development

### Adding New Contract Types

1. Create your contract in Rust using CosmWasm
2. Build to WASM: `cargo wasm`
3. Optimize: `wasm-opt -Oz target/wasm32-unknown-unknown/release/contract.wasm -o optimized.wasm`
4. Store and instantiate using the commands above

### Testing Contracts

The module includes comprehensive testing for:
- Contract storage and instantiation
- State persistence and queries
- Error handling and edge cases
- Gas estimation and optimization

## Integration

This module integrates seamlessly with:
- **Cosmos SDK**: Full compatibility with standard Cosmos modules
- **IBC**: Contracts can interact with other chains via IBC
- **Governance**: Contract parameters can be managed via governance
- **Staking/Distribution**: Contracts can interact with native staking

## Security

- **Deterministic Execution**: All contract execution is deterministic
- **Gas Metering**: Prevents infinite loops and resource exhaustion
- **Sandboxed Environment**: Contracts run in isolated VM instances
- **Permissioned Upgrades**: Admin-controlled contract upgrades (if admin is set)

## Support

For issues or questions:
1. Check transaction logs for error details
2. Verify gas limits and JSON formatting
3. Ensure contract addresses and code IDs are correct
4. Review the CosmWasm documentation for contract-specific issues

## References

- [CosmWasm Documentation](https://docs.cosmwasm.com/)
- [wasmd Repository](https://github.com/CosmWasm/wasmd)
- [CosmWasm Book](https://book.cosmwasm.com/)
- [Cosmos SDK Documentation](https://docs.cosmos.network/) 
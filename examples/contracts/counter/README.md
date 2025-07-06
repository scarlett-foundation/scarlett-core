# Counter Contract

A simple CosmWasm counter contract example for scarlett-core.

## Features

- **Instantiate**: Initialize counter with a starting value
- **Execute**: Increment, decrement, or reset the counter
- **Query**: Get the current counter value
- **Ownership**: Only the owner can reset the counter

## Build

```bash
# Build the contract
cargo build --target wasm32-unknown-unknown --release

# Optimize the contract (optional)
docker run --rm -v "$(pwd)":/code \
  --mount type=volume,source="$(basename "$(pwd)")_cache",target=/code/target \
  --mount type=volume,source=registry_cache,target=/usr/local/cargo/registry \
  cosmwasm/workspace-optimizer:0.15.0
```

## Usage with scarlett-cored

### 1. Store Code
```bash
# Store the contract code
scarlett-cored tx contracts store-code target/wasm32-unknown-unknown/release/counter.wasm \
  --from validator \
  --gas auto \
  --gas-adjustment 1.3 \
  --broadcast-mode block \
  --yes
```

### 2. Instantiate Contract
```bash
# Instantiate with initial count of 0
scarlett-cored tx contracts instantiate-contract 1 \
  '{"count": 0}' \
  --from validator \
  --label "counter-example" \
  --gas auto \
  --gas-adjustment 1.3 \
  --broadcast-mode block \
  --yes
```

### 3. Execute Contract
```bash
# Increment counter
scarlett-cored tx contracts execute-contract <contract-address> \
  '{"increment": {}}' \
  --from validator \
  --gas auto \
  --gas-adjustment 1.3 \
  --broadcast-mode block \
  --yes

# Decrement counter
scarlett-cored tx contracts execute-contract <contract-address> \
  '{"decrement": {}}' \
  --from validator \
  --gas auto \
  --gas-adjustment 1.3 \
  --broadcast-mode block \
  --yes

# Reset counter (only owner)
scarlett-cored tx contracts execute-contract <contract-address> \
  '{"reset": {"count": 10}}' \
  --from validator \
  --gas auto \
  --gas-adjustment 1.3 \
  --broadcast-mode block \
  --yes
```

### 4. Query Contract
```bash
# Get current count
scarlett-cored query contracts smart-contract-state <contract-address> \
  '{"get_count": {}}'
```

## Messages

### InstantiateMsg
```json
{
  "count": 0
}
```

### ExecuteMsg
```json
// Increment
{"increment": {}}

// Decrement  
{"decrement": {}}

// Reset (owner only)
{"reset": {"count": 42}}
```

### QueryMsg
```json
{"get_count": {}}
```

## Response
```json
{
  "count": 42
}
```

## Testing

```bash
# Run unit tests
cargo test

# Run with logging
RUST_LOG=debug cargo test
``` 
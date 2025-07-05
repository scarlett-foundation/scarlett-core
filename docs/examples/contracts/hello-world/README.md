# Hello World Contract

This is a simple "Hello World" smart contract built with CosmWasm 3.0.1, demonstrating modern CosmWasm development patterns.

## Features

- **Modern CosmWasm 3.0.1**: Uses the latest CosmWasm version with updated patterns
- **Type-safe Storage**: Uses `cw-storage-plus` for efficient, type-safe storage
- **Proper Error Handling**: Custom error types with `thiserror` for better error messages
- **Owner-based Access Control**: Only the contract owner can update the greeting
- **Comprehensive Testing**: Unit tests demonstrating contract functionality
- **Schema Generation**: Automatic JSON schema generation for contract messages

## Contract Structure

```
src/
├── contract.rs     # Main contract logic with entry points
├── error.rs        # Custom error types
├── msg.rs          # Message definitions
├── state.rs        # State management with cw-storage-plus
└── lib.rs          # Library exports
examples/
└── schema.rs       # Schema generation script
```

## Messages

### InstantiateMsg
```json
{
  "greeting": "Hello, CosmWasm!"
}
```

### ExecuteMsg
```json
{
  "update_greeting": {
    "greeting": "New greeting message"
  }
}
```

### QueryMsg
```json
{
  "get_greeting": {}
}
```

Or:
```json
{
  "get_state": {}
}
```

## Building the Contract

### Prerequisites
- Rust 1.70+
- `wasm32-unknown-unknown` target

### Build Commands

```bash
# Build the contract
cargo build --release --target wasm32-unknown-unknown

# Optimize the contract (requires docker)
docker run --rm -v "$(pwd)":/code \
  --mount type=volume,source="$(basename "$(pwd)")_cache",target=/code/target \
  --mount type=volume,source=registry_cache,target=/usr/local/cargo/registry \
  cosmwasm/rust-optimizer:0.12.13

# Run tests
cargo test

# Generate schema
cargo run --example schema
```

## Usage Example

1. **Instantiate** the contract with an initial greeting
2. **Query** the greeting using `get_greeting`
3. **Update** the greeting (only owner can do this)
4. **Query** the full state using `get_state`

## Modern CosmWasm Features

This contract demonstrates several modern CosmWasm patterns:

- **Entry Points**: Uses `#[entry_point]` macro for cleaner entry point definitions
- **Type Safety**: Leverages `cw-storage-plus` for type-safe storage operations
- **Error Handling**: Custom error types with `thiserror` for better error messages
- **Testing**: Comprehensive unit tests with mocking
- **Schema Generation**: Automatic JSON schema generation for tooling support
- **Migration Support**: Built-in migration entry point for contract upgrades

## License

This contract is available under the MIT license. 
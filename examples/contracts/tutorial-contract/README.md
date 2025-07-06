# Tutorial Contract

This contract follows the official [CosmWasm Tutorial](https://cosmwasm.cosmos.network/tutorial) step by step.

## Overview

This example demonstrates how to build a CosmWasm smart contract from scratch, following the exact structure and approach outlined in the official documentation at https://cosmwasm.cosmos.network/tutorial/cw-contract/contract-creation.

## Tutorial Steps Covered

- ✅ **Contract Creation**: Set up a new Rust library with proper `Cargo.toml` configuration
- ✅ **Entry Points**: Implementation of instantiate and query functions
- ✅ **Building the Contract**: Compile to WebAssembly (245KB tutorial_contract.wasm)
- ✅ **Creating Queries**: Add query functionality with proper response types
- ✅ **Unit Testing**: Test query function with mock dependencies
- ✅ **MultiTest**: Integration testing with cw-multi-test framework
- ⏳ **Storing State**: State management with cw-storage-plus (next step)
- ⏳ **Execution Messages**: Handle execute messages
- ⏳ **Passing Events**: Event emission
- ⏳ **Handling Funds**: Payment processing
- ⏳ **Good Practices**: Following CosmWasm best practices

## Key Features

Following the tutorial exactly, this contract demonstrates:

- Proper project structure with `crate-type = ["cdylib"]`
- CosmWasm standard library integration
- Step-by-step learning approach
- Educational value for new developers

## Building

```bash
cargo build
```

## Testing

```bash
cargo test
```

## Differences from Other Examples

Unlike the `counter` example which was a standalone implementation, this `tutorial-contract` follows the official CosmWasm tutorial progression exactly, making it ideal for:

- Learning CosmWasm development
- Understanding best practices
- Following documented patterns
- Educational reference

## Next Steps

Continue following the tutorial at: https://cosmwasm.cosmos.network/tutorial/cw-contract/entry-points 
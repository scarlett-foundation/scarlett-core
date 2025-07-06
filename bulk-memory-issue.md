# Bulk Memory Validation Error with wasmvm v3.0.0 and wasmd v0.61.0

## Summary

We're experiencing a consistent "bulk memory support is not enabled" error when attempting to upload CosmWasm contracts to a chain using wasmd v0.61.0 with wasmvm v3.0.0. This occurs despite bulk memory operations being expected to be supported by default in wasmvm v3.0.0.

## Environment

- **wasmd version**: v0.61.0
- **wasmvm version**: v3.0.0  
- **Cosmos SDK version**: v0.53.x
- **CosmWasm contract version**: 2.1.4 (also tested with 1.5.0)
- **Rust version**: 1.88.0
- **Target**: wasm32-unknown-unknown
- **Operating System**: macOS (Darwin 24.1.0)

## Issue Description

When attempting to upload any WASM contract using the `store-code` command, we consistently receive the following error:

```
rpc error: code = Unknown desc = rpc error: code = Unknown desc = failed to execute message; message index: 0: failed to store wasm code: Error calling the VM: Error during static Wasm validation: Wasm bytecode could not be deserialized. Deserialization error: "bulk memory support is not enabled (at offset 0x3221)": create wasm contract failed [!cosm!wasm/wasmd@v0.61.0/x/wasm/keeper/keeper.go:195] with gas used: '3367598': unknown request
```

## Steps to Reproduce

1. Set up a chain using wasmd v0.61.0 with wasmvm v3.0.0
2. Create a simple CosmWasm contract (tested with both tutorial examples and basic counter contracts)
3. Compile contract to WASM: `cargo build --target wasm32-unknown-unknown --release`
4. Attempt to upload: `scarlett-cored tx contracts store-code contract.wasm --from alice --chain-id test --gas auto --gas-adjustment 1.5 --yes`
5. Observe the bulk memory validation error

## What We've Tried

### 1. Different Compilation Flags
```bash
# Attempted to disable bulk memory operations
RUSTFLAGS='-C target-feature=-bulk-memory' cargo build --target wasm32-unknown-unknown --release

# Attempted to disable multiple WebAssembly features
RUSTFLAGS='-C target-feature=-bulk-memory,-sign-ext,-simd128,-reference-types,-multi-value' cargo build --target wasm32-unknown-unknown --release
```

### 2. Different CosmWasm Versions
- Tested with cosmwasm-std 2.1.4 (latest)
- Downgraded to cosmwasm-std 1.5.0 for compatibility
- Both versions produce the same error

### 3. VM Configuration
Our current VM configuration:
```go
// Set available capabilities for wasmd v0.61.0 + wasmvm v3.0.0 compatibility
supportedFeatures := []string{
    "iterator", "staking", "stargate",
    "cosmwasm_1_1", "cosmwasm_1_2", "cosmwasm_1_3", "cosmwasm_1_4",
    "cosmwasm_2_0", "cosmwasm_2_1", "cosmwasm_2_2",
}

// Using empty VMConfig and NodeConfig to allow defaults
vmConfig := wasmtypes.VMConfig{}
nodeConfig := wasmtypes.NodeConfig{}
```

### 4. Multiple Contract Types
- Simple tutorial contracts with basic instantiate/query functions
- Counter contracts with state management
- All produce the same error at the same offset (0x3221)

## Expected Behavior

Based on the wasmvm v3.0.0 documentation and WebAssembly 2.0 support, bulk memory operations should be supported by default. We expect:

1. Contracts compiled with modern Rust toolchain should upload successfully
2. Bulk memory operations should be enabled by default in wasmvm v3.0.0
3. No additional configuration should be required for basic WebAssembly 2.0 features

## Questions

1. **Is additional VM configuration required** to enable bulk memory support in wasmvm v3.0.0?
2. **Are there specific VMConfig parameters** that need to be set for WebAssembly 2.0 features?
3. **Is this a known compatibility issue** between wasmd v0.61.0 and wasmvm v3.0.0?
4. **Should we be using different compilation flags** for contracts targeting wasmvm v3.0.0?

## Additional Context

- The error occurs consistently across different contract types and compilation approaches
- The same offset (0x3221) appears in all attempts, suggesting a systematic issue
- Our VM is properly initialized with lazy loading to avoid lock conflicts
- Chain functionality works correctly otherwise (queries, transactions, etc.)

## Potential Solutions or Requests

1. **Documentation** on proper VM configuration for bulk memory support
2. **Example configuration** showing how to enable WebAssembly 2.0 features
3. **Guidance** on contract compilation best practices for wasmvm v3.0.0
4. **Confirmation** of the expected compatibility matrix between wasmd/wasmvm versions

Any guidance on resolving this issue or proper configuration would be greatly appreciated. We're happy to provide additional debugging information or test potential solutions.

---

**Repository**: This issue is being filed in both wasmd and wasmvm repositories as it appears to be related to the integration between these components. 
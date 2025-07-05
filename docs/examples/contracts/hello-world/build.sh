#!/bin/bash

# Build script for Hello World CosmWasm Contract

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo -e "${YELLOW}🏗️  Building Hello World CosmWasm Contract...${NC}"

# Check if rust is installed
if ! command -v cargo &> /dev/null; then
    echo -e "${RED}❌ Rust is not installed. Please install Rust first.${NC}"
    exit 1
fi

# Check if wasm32 target is installed
if ! rustup target list --installed | grep -q "wasm32-unknown-unknown"; then
    echo -e "${YELLOW}📦 Installing wasm32-unknown-unknown target...${NC}"
    rustup target add wasm32-unknown-unknown
fi

# Clean previous builds
echo -e "${YELLOW}🧹 Cleaning previous builds...${NC}"
cargo clean

# Build optimized binary
RUSTFLAGS='-C link-arg=-s' cargo build --release --target wasm32-unknown-unknown --locked

# Create release directory if it doesn't exist
mkdir -p artifacts

# Optimize the Wasm binary
wasm-opt -Os ./target/wasm32-unknown-unknown/release/hello_world.wasm -o ./artifacts/hello_world.wasm

# Verify the output
ls -l ./artifacts/hello_world.wasm

# Check if the build was successful
if [ -f "target/wasm32-unknown-unknown/release/hello_world.wasm" ]; then
    echo -e "${GREEN}✅ Contract built successfully!${NC}"
    
    # Get file size
    size=$(ls -lh target/wasm32-unknown-unknown/release/hello_world.wasm | awk '{print $5}')
    echo -e "${GREEN}📁 Contract size: ${size}${NC}"
    
    # Copy to current directory for convenience
    cp target/wasm32-unknown-unknown/release/hello_world.wasm ./hello_world.wasm
    echo -e "${GREEN}📋 Contract copied to: ./hello_world.wasm${NC}"
else
    echo -e "${RED}❌ Build failed!${NC}"
    exit 1
fi

# Schema generation removed (optional feature)

# Run tests
echo -e "${YELLOW}🧪 Running tests...${NC}"
cargo test

echo -e "${GREEN}🎉 Build complete!${NC}"
echo -e "${GREEN}📂 Contract: ./hello_world.wasm${NC}"

echo -e "${YELLOW}💡 Next steps:${NC}"
echo -e "   1. Deploy using: scarlett-cored tx contracts deploy-contract"
echo -e "   2. Register for emissions funding (optional)"
echo -e "   3. Test the deployed contract" 
#!/bin/bash

# Define necessary variables
WILLCHAIN_CHAIN_ID="willchain-mainnet"
WILLCHAIN_NODE=http://localhost:26657
DEV_WALLET=alice

# Path to your wasm file (update this accordingly)
WASM_FILE="./wasm_artifacts/simple_option.wasm"

# Deploy the contract
echo "Deploying contract..."
RES=$(./build/wasmd tx wasm store $WASM_FILE --from $DEV_WALLET --gas auto --gas-adjustment 1.3 -y -b sync --output json --chain-id="$WILLCHAIN_CHAIN_ID" --node="$WILLCHAIN_NODE")
echo "$RES"

# Extract the Code ID from the response
# Ensure jq is installed or manually extract CODE_ID from the printed response
CODE_ID=$(echo $RES | jq -r '.logs[0].events[-1].attributes[-1].value')
echo "Code ID: $CODE_ID"

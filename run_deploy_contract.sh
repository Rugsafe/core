#!/bin/bash

# Define necessary variables
W3LL_CHAIN_ID="w3ll-chain"
W3LL_NODE=http://localhost:26657
DEV_WALLET=alice

# Path to your wasm file (update this accordingly)
WASM_FILE="./simple_option.wasm"

# Deploy the contract
echo "Deploying contract..."
RES=$(./build/wasmd tx wasm store $WASM_FILE --from $DEV_WALLET --gas auto --gas-adjustment 1.3 -y -b sync --output json --chain-id="$W3LL_CHAIN_ID" --node="$W3LL_NODE")
echo "$RES"

# Extract the Code ID from the response
# Ensure jq is installed or manually extract CODE_ID from the printed response
CODE_ID=$(echo $RES | jq -r '.logs[0].events[-1].attributes[-1].value')
echo "Code ID: $CODE_ID"

#!/bin/bash

# Define necessary variables
W3LL_CHAIN_ID="w3ll-mainnet"
W3LL_NODE=http://localhost:26657
DEV_WALLET=alice

# Path to your wasm file (update this accordingly)
WASM_FILE="./wasm_artifacts/ibc_tutorial.wasm"

# Deploy the contract
echo "Deploying IBC contract..."
RES=$(./build/wasmd tx wasm store $WASM_FILE --from $DEV_WALLET --gas auto --gas-adjustment 1.3 -y -b sync --output json --chain-id="$W3LL_CHAIN_ID" --node="$W3LL_NODE")
echo "$RES"

# Extract the Code ID from the response
# Ensure jq is installed or manually extract CODE_ID from the printed response
CODE_ID=$(echo $RES | jq -r '.logs[0].events[-1].attributes[-1].value')
echo "Code ID: $CODE_ID"

# check
./build/wasmd q wasm code-info $(CODE_ID)

#instantiate
./build/wasmd tx wasm instantiate $(CODE_ID) \
"{}" \
--amount="1w3ll" --no-admin --label "awesomwasm token" --from $DEV_WALLET --gas auto --gas-adjustment 1.3 -b block -y

# contract address
./build/wasmd q wasm list-contract-by-code $(CODE_ID) --output json | jq -r '.contracts[-1]'
CONTRACT_ADDRESS=$(shell cored-00 q wasm list-contract-by-code $(CODE_ID) --output json $(COREUM_NODE_ARGS) $(COREUM_CHAIN_ID_ARGS) | jq -r '.contracts[-1]')
echo $$CONTRACT_ADDRESS
#!/bin/bash

# Default chain name
CHAIN_NAME=${CHAINNAME:-"w3ll-mainnet"}

# Define paths
HOME_DIR="./private/.${CHAIN_NAME}"
MNEMONIC_FILE="seed.txt"

make build

echo "step 1"
./build/wasmd version

# Clean up previous runs
rm -rf ${HOME_DIR}
rm -rf ~/.${CHAIN_NAME}/

# Initialize with dynamic chain name
./build/wasmd init ${CHAIN_NAME} --home ${HOME_DIR} --chain-id ${CHAIN_NAME}

echo "step 3"
./build/wasmd keys list --home ${HOME_DIR}

# Check if mnemonic exists and import alice using it
if [ -f "$MNEMONIC_FILE" ]; then
    echo "Importing alice using existing mnemonic..."
    MNEMONIC=$(cat $MNEMONIC_FILE)
    echo "$MNEMONIC" | ./build/wasmd keys add alice --recover --home ${HOME_DIR} --keyring-backend test
else
    echo "Mnemonic file not found, generating new account for alice..."
    ./build/wasmd keys add alice --home ${HOME_DIR} --keyring-backend test
fi

echo "step 4.5"
grep bond_denom ${HOME_DIR}/config/genesis.json

# Add new account to genesis
echo "step 5"
./build/wasmd genesis add-genesis-account alice 111111111111stake --home ${HOME_DIR} --keyring-backend test
./build/wasmd genesis add-genesis-account alice 900000000000w3ll --home ${HOME_DIR} --keyring-backend test --append

# Generate a genesis tx carrying a self delegation
echo "step 6"
./build/wasmd genesis gentx alice 222222222stake --home ${HOME_DIR} --keyring-backend test --chain-id ${CHAIN_NAME}
echo "step 7"
./build/wasmd genesis gentx alice 700000000w3ll --home ${HOME_DIR} --keyring-backend test --chain-id ${CHAIN_NAME}
./build/wasmd genesis collect-gentxs --home ${HOME_DIR}

echo "step 8"
./build/wasmd start --home ${HOME_DIR} --rpc.laddr=tcp://0.0.0.0:26657 --grpc.address=0.0.0.0:9090

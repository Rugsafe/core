#!/bin/bash

# Define paths
HOME_DIR="./private/.wasmapp"
MNEMONIC_FILE="seed.txt"

make build

echo "step 1"
./build/wasmd version

# Clean up previous runs
rm -rf ${HOME_DIR}
rm -rf ~/.wasmapp/

# Initialize
./build/wasmd init w3ll-chain --home ${HOME_DIR} --chain-id w3ll-chain

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

# You can also save the newly generated mnemonic to a file if needed
# Note: This is not recommended for production use due to security implications

echo "step 4.5"
grep bond_denom ${HOME_DIR}/config/genesis.json

# Add new account to genesis
echo "step 5"
./build/wasmd genesis add-genesis-account alice 111111111111stake --home ${HOME_DIR} --keyring-backend test
./build/wasmd genesis add-genesis-account alice 900000000000w3ll --home ${HOME_DIR} --keyring-backend test --append

# Generate a genesis tx carrying a self delegation
echo "step 6"
./build/wasmd genesis gentx alice 222222222stake --home ${HOME_DIR} --keyring-backend test --chain-id w3ll-chain
echo "step 7"
./build/wasmd genesis gentx alice 700000000w3ll --home ${HOME_DIR} --keyring-backend test --chain-id w3ll-chain
./build/wasmd genesis collect-gentxs --home ${HOME_DIR}

echo "step 8"
./build/wasmd start --home ${HOME_DIR} --rpc.laddr=tcp://0.0.0.0:26657
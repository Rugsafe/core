#!/bin/bash

# Define constants
CHAIN_ID="willchain-mainnet"
DENOM="uwill"
AMOUNT="1000000$DENOM"
RECIPIENT_ADDRESS="<recipient-address>"  # Update with the correct address
WASMD="./build/wasmd"

# Checking and creating keys
echo "Checking existing keys..."
$WASMD keys show alice -a || $WASMD keys add alice

echo "Recreating keys for Bob and Charlie..."
$WASMD keys delete bob --keyring-backend test -y
$WASMD keys delete charlie --keyring-backend test -y
$WASMD keys add bob --keyring-backend test
$WASMD keys add charlie --keyring-backend test

# Creating a multisig account
echo "Creating a 3-of-2 multisig account..."
MULTISIG_NAME="multisig-abc"
$WASMD keys add $MULTISIG_NAME --multisig="alice,bob,charlie" --multisig-threshold=2 --keyring-backend test

# Preparing a transaction
echo "Preparing an unsigned transaction..."
MULTISIG_ADDRESS=$($WASMD keys show $MULTISIG_NAME -a --keyring-backend test)
$WASMD tx bank send $MULTISIG_ADDRESS $RECIPIENT_ADDRESS $AMOUNT --chain-id $CHAIN_ID --generate-only > unsignedTx.json

# Signing the transaction
echo "Signing the transaction with Alice and Bob..."
$WASMD tx sign unsignedTx.json --from=alice --chain-id=$CHAIN_ID --output-document=aliceSignedTx.json --keyring-backend test
$WASMD tx sign unsignedTx.json --from=bob --chain-id=$CHAIN_ID --output-document=bobSignedTx.json --keyring-backend test

# Combining signatures
echo "Combining signatures..."
$WASMD tx multisign unsignedTx.json $MULTISIG_NAME aliceSignedTx.json bobSignedTx.json --keyring-backend test > signedTx.json

# Broadcasting the transaction
echo "Broadcasting the transaction..."
$WASMD tx broadcast signedTx.json

# Cleaning up
echo "Cleaning up: deleting Bob and Charlie's keys..."
$WASMD keys delete bob --keyring-backend test -y
$WASMD keys delete charlie --keyring-backend test -y

echo "Transaction process complete."

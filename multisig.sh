#!/bin/bash

# Set variables
CHAIN_ID="willchain-mainnet"
DENOM="uwill"
AMOUNT="1000000$DENOM"
RECIPIENT_ADDRESS="will156mw28alhpenp4lknweat6432dux34uydx590v" # Change this to the recipient's address
WASMD="./build/wasmd"

# Check if Alice's key exists, otherwise create it
if ! $WASMD keys show alice -a; then
    echo "Alice's key does not exist, creating now..."
    $WASMD keys add alice
fi

# Create keys for Bob and Charlie, and delete them each time script runs
echo "Creating keys for Bob and Charlie..."
$WASMD keys add bob
$WASMD keys add charlie

# Create the multisig key
echo "Creating a 3-of-2 multisig account..."
$WASMD keys add multisig-abc --multisig="alice,bob,charlie" --multisig-threshold=2

# Get multisig address
MULTISIG_ADDRESS=$($WASMD keys show multisig-abc -a)

# Create an unsigned transaction
echo "Creating an unsigned transaction..."
$WASMD tx bank send $MULTISIG_ADDRESS $RECIPIENT_ADDRESS $AMOUNT --chain-id=$CHAIN_ID --generate-only > unsignedTx.json

# Sign the transaction with Alice's key
echo "Signing the transaction with Alice's key..."
$WASMD tx sign unsignedTx.json --from=alice --chain-id=$CHAIN_ID --output-document=aliceSignedTx.json

# Sign the transaction with Bob's key
echo "Signing the transaction with Bob's key..."
$WASMD tx sign unsignedTx.json --from=bob --chain-id=$CHAIN_ID --output-document=bobSignedTx.json

# Combine signatures
echo "Combining signatures from Alice and Bob..."
$WASMD tx multisign unsignedTx.json multisig-abc aliceSignedTx.json bobSignedTx.json > signedTx.json

# Broadcast the transaction
echo "Broadcasting the transaction..."
$WASMD tx broadcast signedTx.json

# Cleanup: delete Bob and Charlie's keys
echo "Cleaning up: deleting Bob and Charlie's keys..."
$WASMD keys delete bob -y
$WASMD keys delete charlie -y

echo "Transaction process complete."
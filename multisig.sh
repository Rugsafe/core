#!/bin/bash

# Define constants
CHAIN_ID="willchain-mainnet"
DENOM="uwill"
AMOUNT="1000000$DENOM"
RECIPIENT_ADDRESS="will156mw28alhpenp4lknweat6432dux34uydx590v"  # Update with the actual recipient's address
WASMD="./build/wasmd"
KEYRING_BACKEND="test"  # Using 'os' for persistent storage
HOME_DIR="./private/.${CHAIN_NAME}" # Custom directory for wasmd configurations and keyring

# Remove previous transaction files
echo "Removing previous transaction files..."
echo "rm -f unsignedTx.json aliceSignedTx.json bobSignedTx.json signedTx.json"
rm -f unsignedTx.json aliceSignedTx.json bobSignedTx.json signedTx.json

# Check if Alice exists
echo "Checking if Alice exists..."
echo "$WASMD keys show alice --home $HOME_DIR --keyring-backend $KEYRING_BACKEND -a || $WASMD keys add alice --home $HOME_DIR --keyring-backend $KEYRING_BACKEND"
$WASMD keys show alice --home $HOME_DIR --keyring-backend $KEYRING_BACKEND -a || $WASMD keys add alice --home $HOME_DIR --keyring-backend $KEYRING_BACKEND

# Handle keys for Bob and Charlie
for key in bob charlie; do
    echo "Handling key for $key..."
    echo "$WASMD keys show $key --home $HOME_DIR --keyring-backend $KEYRING_BACKEND -a || $WASMD keys add $key --home $HOME_DIR --keyring-backend $KEYRING_BACKEND"
    $WASMD keys show $key --home $HOME_DIR --keyring-backend $KEYRING_BACKEND -a || $WASMD keys add $key --home $HOME_DIR --keyring-backend $KEYRING_BACKEND
done

# Creating a multisig account from Alice, Bob, and Charlie
echo "Creating a 3-of-2 multisig account..."
ALICE_ADDRESS=$($WASMD keys show alice --home $HOME_DIR --keyring-backend $KEYRING_BACKEND -a)
BOB_ADDRESS=$($WASMD keys show bob --home $HOME_DIR --keyring-backend $KEYRING_BACKEND -a)
CHARLIE_ADDRESS=$($WASMD keys show charlie --home $HOME_DIR --keyring-backend $KEYRING_BACKEND -a)

ALICE_ADDRESS=will1p0k8gygawzpggzwftv7cv47zvgg8zaun7h2v28
BOB_ADDRESS=will1499sj4q9z0mkpcmfs0cgphp4xm27km723kc6t9
CHARLIE_ADDRESS=will19j3lsnklvxxff59lkgyn3vvj72vwcpsdf9gx8h

# ALICE_ADDRESS=alice
# BOB_ADDRESS=bob
# CHARLIE_ADDRESS=charlie
MULTISIG_ADDRESS="multisig-abc"

echo "$WASMD keys add $MULTISIG_ADDRESS --multisig='$ALICE_ADDRESS,$BOB_ADDRESS,$CHARLIE_ADDRESS' --multisig-threshold=2 --home $HOME_DIR --keyring-backend $KEYRING_BACKEND"
$WASMD keys add $MULTISIG_ADDRESS --multisig="$ALICE_ADDRESS,$BOB_ADDRESS,$CHARLIE_ADDRESS" --multisig-threshold=2 --home $HOME_DIR --keyring-backend $KEYRING_BACKEND
# $WASMD keys add $MULTISIG_ADDRESS --multisig="alice,bob,charlie" --multisig-threshold=2 --home $HOME_DIR --keyring-backend $KEYRING_BACKEND

# exit 1
# Preparing an unsigned transaction
echo "Preparing an unsigned transaction..."
MULTISIG_ACC_ADDRESS=$($WASMD keys show $MULTISIG_ADDRESS --home $HOME_DIR --keyring-backend $KEYRING_BACKEND -a)
echo "$WASMD tx bank send $MULTISIG_ACC_ADDRESS $RECIPIENT_ADDRESS $AMOUNT --chain-id $CHAIN_ID --generate-only > unsignedTx.json"
# $WASMD tx bank send $MULTISIG_ACC_ADDRESS $RECIPIENT_ADDRESS $AMOUNT --chain-id $CHAIN_ID --generate-only > unsignedTx.json
$WASMD tx bank send $ALICE_ADDRESS $RECIPIENT_ADDRESS $AMOUNT --chain-id $CHAIN_ID --from $ALICE_ADDRESS --generate-only > unsignedTx.json


# Signing the transaction with Alice's and Bob's keys
echo "Alice signing the transaction..."
echo "$WASMD tx sign unsignedTx.json --from alice --chain-id $CHAIN_ID --keyring-backend $KEYRING_BACKEND --output-document=aliceSignedTx.json"
# $WASMD tx sign unsignedTx.json --from alice --chain-id $CHAIN_ID --keyring-backend $KEYRING_BACKEND --output-document=aliceSignedTx.json
$WASMD tx sign unsignedTx.json --from alice --chain-id $CHAIN_ID --keyring-backend $KEYRING_BACKEND --output-document=aliceSignedTx.json
exit 1
echo "Bob signing the transaction..."
echo "$WASMD tx sign unsignedTx.json --from bob --chain-id $CHAIN_ID --keyring-backend $KEYRING_BACKEND --output-document=bobSignedTx.json"
$WASMD tx sign unsignedTx.json --from bob --chain-id $CHAIN_ID --keyring-backend $KEYRING_BACKEND --output-document=bobSignedTx.json

# Combining signatures
echo "Combining signatures..."
echo "$WASMD tx multisign unsignedTx.json $MULTISIG_ADDRESS aliceSignedTx.json bobSignedTx.json --keyring-backend $KEYRING_BACKEND > signedTx.json"
$WASMD tx multisign unsignedTx.json $MULTISIG_ADDRESS aliceSignedTx.json bobSignedTx.json --keyring-backend $KEYRING_BACKEND > signedTx.json

# Broadcasting the transaction
echo "Broadcasting the transaction..."
echo "$WASMD tx broadcast signedTx.json --keyring-backend $KEYRING_BACKEND"
$WASMD tx broadcast signedTx.json --keyring-backend $KEYRING_BACKEND

# Cleaning up
# echo "Cleaning up: deleting Bob and Charlie's keys..."
# echo "$WASMD keys delete bob --home $HOME_DIR --keyring-backend $KEYRING_BACKEND -y"
# $WASMD keys delete bob --home $HOME_DIR --keyring-backend $KEYRING_BACKEND -y
# echo "$WASMD keys delete charlie --home $HOME_DIR --keyring-backend $KEYRING_BACKEND -y"
# $WASMD keys delete charlie --home $HOME_DIR --keyring-backend $KEYRING_BACKEND -y

echo "Transaction process complete."

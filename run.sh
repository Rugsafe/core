#!/bin/bash



make build


echo "step 1"
./build/wasmd version



rm -rf ./private/.wasmapp
# sleep 3
rm -rf ~/.wasmapp/



./build/wasmd init w3ll-chain --home ./private/.wasmapp --chain-id w3ll-chain
# ./build/simd init demo --chain-id w3ll-chain


# echo "step 2"
#COPY OVER GENESIS
# cp ~/.simapp/config/genesis.json ./private/.simapp/config/genesis.json
# cp ~/.simapp/config/config.toml ./private/.simapp/config/config.toml


echo "step 3"
./build/wasmd keys list --home ./private/.wasmapp



echo "step 4"
./build/wasmd keys add alice --home ./private/.wasmapp --keyring-backend test
# ./build/simd keys show alice --keyring-backend test


echo "step 4.5"
grep bond_denom ./private/.wasmapp/config/genesis.json



# add a new account to genesis
echo "step 5"
./build/wasmd genesis add-genesis-account alice 111111111111stake --home ./private/.wasmapp --keyring-backend test
# give them w3ll
./build/wasmd genesis add-genesis-account alice 900000000000w3ll --home ./private/.wasmapp --keyring-backend test --append



#./build/simd add-genesis-account alice 100000000stake --home ./private/.simapp --keyring-backend test


#Generate a genesis tx carrying a self delegation
echo "step 6"
./build/wasmd genesis gentx alice 222222222stake --home ./private/.wasmapp --keyring-backend test --chain-id w3ll-chain
echo "step 7"
./build/simd genesis gentx alice 700000000w3ll --home ./private/.simapp --keyring-backend test --chain-id w3ll-chain
./build/wasmd genesis collect-gentxs --home ./private/.wasmapp




echo "step 8"
./build/wasmd start --home ./private/.wasmapp
#./build/simd start --home ~/.simapp


# export alice=$(./build/simd keys show alice --address --home ./private/.simapp --keyring-backend test)


# DEV

# ./build/simd tx bank send alice w3ll1udgssjwdyzeuadj9eed4vj6zj4uhlxdg6nwemd 100000000w3ll -b sync --chain-id="w3ll-chain"
# ./build/simd q bank balances w3ll1udgssjwdyzeuadj9eed4vj6zj4uhlxdg6nwemd
# ./build/simd keys list
# # from genesis
# ./build/simd keys add alice --recover 
# 0xjovis-MBP:cosmos-sdk 0xjovi$ ./build/simd keys delete alice -y
# ./build/simd tx wasm store ./simple_option.wasm --from dev-wallet --gas auto --gas-adjustment 1.3 -y -b block --output json --chain-id="w3ll-chain"
# docker run --rm -v "$(pwd)":/code \
#   --mount type=volume,source="$(basename "$(pwd)")_cache",target=/code/target \
#   --mount type=volume,source=registry_cache,target=/usr/local/cargo/registry \
#   cosmwasm/rust-optimizer:0.14.0

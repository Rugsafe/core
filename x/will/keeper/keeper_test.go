package keeper_test

import (
	"crypto/sha256"
	"fmt"
	"testing"

	// Import the tm-db package
	// dbm "github.com/tendermint/tm-db" // Import the tm-db package
	"github.com/bwesterb/go-ristretto"
	dbm "github.com/cosmos/cosmos-db"
	ibctransfertypes "github.com/cosmos/ibc-go/v8/modules/apps/transfer/types"
	ibcexported "github.com/cosmos/ibc-go/v8/modules/core/exported"

	// _proto "github.com/cosmos/gogoproto/proto"

	// "github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	// "cosmossdk.io/core/store"
	// corestoretypes "cosmossdk.io/core/store"
	"cosmossdk.io/log"
	corestore "cosmossdk.io/store"

	// corestoretypes "cosmossdk.io/core/store"
	storemetrics "cosmossdk.io/store/metrics"
	storetypes "cosmossdk.io/store/types"
	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"

	// codectypes "github.com/cosmos/cosmos-sdk/codec/types"

	"github.com/cosmos/cosmos-sdk/runtime"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/CosmWasm/wasmd/app"
	"github.com/CosmWasm/wasmd/x/will/keeper"
	"github.com/CosmWasm/wasmd/x/will/schemes/pedersen"
	"github.com/CosmWasm/wasmd/x/will/types"
)

func setupKeeper(t *testing.T) (*keeper.Keeper, sdk.Context) {
	// func setupKeeper(t *testing.T) *keeper.Keeper {
	// w3llApp, ctx := app.Setup(t)
	willchainApp := app.Setup(t)

	//config for native denom
	// cfg := sdk.GetConfig()
	// cfg.SetBech32PrefixForAccount("will", "willpub")

	mockedCodec := willchainApp.AppCodec()

	// channelKeeper := w3llApp.IBCKeeper.ChannelKeeper
	// channelKeeper := w3llApp.WillKeeper.ChannelKeeper
	channelKeeper := willchainApp.GetIBCKeeper().ChannelKeeper
	// scopedKeeper := w3llApp.ScopedIBCKeeper
	// scopedKeeper := w3llApp.WillKeeper.ScopedKeeper
	scopedKeeper := willchainApp.ScopedIBCKeeper
	fmt.Println(channelKeeper)
	fmt.Println(scopedKeeper)
	// Initialize DB and store
	memDB := dbm.NewMemDB()
	// ms := corestore.NewCommitMultiStore(memDB) // Initialize the MultiStore with the in-memory DB
	ms := corestore.NewCommitMultiStore(memDB, log.NewTestLogger(t), storemetrics.NewNoOpMetrics())
	// keyWill := storetypes.NewKVStoreKey(types.StoreKey)
	storeKey := storetypes.NewKVStoreKey(types.StoreKey)
	ibcStoreKey := storetypes.NewKVStoreKey(ibctransfertypes.StoreKey)    // IBC store key
	ibcExportedStoreKey := storetypes.NewKVStoreKey(ibcexported.StoreKey) // IBC store key

	storeservice := runtime.NewKVStoreService(storeKey)

	// Create and mount store keys
	// ms.MountStoreWithDB(keyWill, storetypes.StoreTypeIAVL, memDB)
	ms.MountStoreWithDB(storeKey, storetypes.StoreTypeIAVL, memDB)
	ms.MountStoreWithDB(ibcStoreKey, storetypes.StoreTypeIAVL, memDB)         // Mount the IBC store
	ms.MountStoreWithDB(ibcExportedStoreKey, storetypes.StoreTypeIAVL, memDB) // Mount the IBC store

	//dev
	bankStoreKey := storetypes.NewKVStoreKey("acc")
	ms.MountStoreWithDB(bankStoreKey, storetypes.StoreTypeIAVL, memDB)
	bankStoreKey2 := storetypes.NewKVStoreKey("bank")
	ms.MountStoreWithDB(bankStoreKey2, storetypes.StoreTypeIAVL, memDB)

	// ms.MountStoreWithDB(keyAcc, storetypes.StoreTypeIAVL, memDB)
	// ms.MountStoreWithDB(string("acc"), storetypes.StoreTypeIAVL, memDB)

	require.NoError(t, ms.LoadLatestVersion())

	// Create context
	ctx := sdk.NewContext(ms, tmproto.Header{}, true, log.NewNopLogger())
	// clientKeeper := clientkeeper.NewKeeper(cdc, key, paramSpace, stakingKeeper, upgradeKeeper)
	// connectionKeeper := connectionkeeper.NewKeeper(cdc, key, paramSpace, clientKeeper)
	// portKeeper := portkeeper.NewKeeper(scopedKeeper)
	// channelKeeper := channelkeeper.NewKeeper(cdc, key, clientKeeper, connectionKeeper, &portKeeper, scopedKeeper)

	// Initialize keeper with the store key
	// TODO: FIX
	// k := keeper.NewKeeper(mockedCodec, storeservice, nil, channelKeeper, scopedKeeper)
	// w3llApp.PermissionedWasmKeeper = *wasmkeeper.NewDefaultPermissionKeeper(app.WasmKeeper)
	k := keeper.NewKeeper(
		mockedCodec,
		storeservice,
		nil,
		channelKeeper,
		willchainApp.GetIBCKeeper().PortKeeper,
		willchainApp.ScopedWillKeeper,
		willchainApp.ScopedIBCKeeper,
		*willchainApp.CapabilityKeeper,
		willchainApp.WasmKeeper,
		willchainApp.GetBankKeeper(),
		willchainApp.PermissionedWasmKeeper,
		willchainApp.GetAccountKeeper(),
	)
	return &k, ctx
	// return &k
}

// /*
func TestKeeperCreateWill(t *testing.T) {
	kpr, ctx := setupKeeper(t)

	msg := &types.MsgCreateWillRequest{
		Creator:     "creator-address",
		Name:        "Test Will",
		Beneficiary: "beneficiary-address",
		Height:      1,
		Components:  []*types.ExecutionComponent{},
	}

	// Assuming CreateWill method is correctly implemented to handle mocks
	will, err := kpr.CreateWill(sdk.UnwrapSDKContext(ctx), msg)
	require.NoError(t, err)
	assert.NotNil(t, will)

	// Assertions to verify the contents of the will
	assert.Equal(t, msg.Creator, will.Creator, "will creator should match the request")
	assert.Equal(t, msg.Name, will.Name, "will name should match the request")
	assert.Equal(t, msg.Beneficiary, will.Beneficiary, "will beneficiary should match the request")
	assert.Equal(t, msg.Height, will.Height, "will height should match the request")

	assert.Len(t, will.Components, len(msg.Components), "number of will components should match the request")

	// Retrieve the will by ID
	retrievedWill, err := kpr.GetWillByID(sdk.UnwrapSDKContext(ctx), will.ID)
	require.NoError(t, err, "failed to retrieve will by ID")
	assert.NotNil(t, retrievedWill, "retrieved will should not be nil")

	// Compare the retrieved will with the created will
	assert.Equal(t, will.ID, retrievedWill.ID, "retrieved will ID should match the created will ID")
	assert.Equal(t, will.Creator, retrievedWill.Creator, "retrieved will creator should match")
	assert.Equal(t, will.Name, retrievedWill.Name, "retrieved will name should match")
	assert.Equal(t, will.Beneficiary, retrievedWill.Beneficiary, "retrieved will beneficiary should match")
	assert.Equal(t, will.Height, retrievedWill.Height, "retrieved will height should match")
	assert.Equal(t, will.Status, retrievedWill.Status, "retrieved will status should match")
	// Add more assertions as needed to compare other fields
}

// TODO: write test for will execution transfer component

func TestKeeperClaimWithSchnorrSignature(t *testing.T) {
	kpr, ctx := setupKeeper(t)

	// Hardcoded values from your Schnorr signature example
	publicKeyHex := "2320a2da28561875cedbb0c25ae458e0a1d087834ae49b96a3f93cec79a8190c"
	signatureHex := "7ab0edb9b0929b5bb4b47dfb927d071ecc5de75985662032bb52ef3c5ace640b165c6df5ea8911a6c0195a3140be5119a5b882e91b34cbcdd31ef3f5b0035b06"

	message := "message-2b-signed"
	creator := "creator-address"

	msg := &types.MsgCreateWillRequest{
		Creator:     creator,
		Name:        "Test Will",
		Beneficiary: "beneficiary-address",
		Height:      2,
		Components: []*types.ExecutionComponent{
			{
				Name:   "SchnorrSignatureComponent",
				Id:     "abc",
				Status: "inactive", // TODO: should be set by keeper upon createWill
				ComponentType: &types.ExecutionComponent_Claim{
					Claim: &types.ClaimComponent{
						SchemeType: &types.ClaimComponent_Schnorr{
							Schnorr: &types.SchnorrSignature{
								PublicKey: []byte(publicKeyHex),
								// Signature: []byte(signatureHex),
								Message: message,
							},
						},
					},
				},
			},
		},
	}

	// CreateWill method is correctly implemented to handle mocks
	will, err := kpr.CreateWill(sdk.UnwrapSDKContext(ctx), msg)
	require.NoError(t, err)
	assert.NotNil(t, will)

	fmt.Println("====[TEST]HERE IS THE PRACTICE WILL[TEST]=====")
	fmt.Println((will))

	// a will has already been created and you have its ID
	willID := will.ID // Replace with the actual will ID
	componentID := will.Components[0].Id

	// verify will upon creation time status is live
	require.Equal(t, will.Status, "live")

	// verify the will claimable component is inactive
	require.Equal(t, will.Components[0].Status, "inactive")

	// roll height forward
	ctx_future := sdk.UnwrapSDKContext(ctx).WithBlockHeight(2)

	// run the begin blocker with the updated block height for the will to execute.
	kpr.BeginBlocker(ctx_future)

	// verify will is expired now, and claimable
	will_for_claimable_check, err_claimable_check := kpr.GetWillByID(ctx, will.ID)
	require.NoError(t, err_claimable_check)
	fmt.Println("WILL FOR CLAIMABLE CHECK:")
	fmt.Println(will_for_claimable_check)
	// verify will is now expired, since after expiry
	require.Equal(t, will_for_claimable_check.Status, "expired")

	// verify the will's claimable component is now active
	require.Equal(t, will_for_claimable_check.Components[0].Status, "active")

	// Construct the claim request with the Schnorr claim
	claimMsg := &types.MsgClaimRequest{
		WillId:      willID,
		Claimer:     creator,
		ComponentId: componentID,
		ClaimType: &types.MsgClaimRequest_SchnorrClaim{
			SchnorrClaim: &types.SchnorrClaim{
				PublicKey: []byte(publicKeyHex),
				Signature: []byte(signatureHex),
				Message:   message,
			},
		},
	}

	// Process the claim
	err = kpr.Claim(sdk.UnwrapSDKContext(ctx_future), claimMsg)
	require.NoError(t, err, "processing Schnorr claim should not produce an error")
	// TODO: verify will components status is not active anymore after successful claim

	will_for_status_check, err_status_check := kpr.GetWillByID(ctx, will.ID)
	require.NoError(t, err_status_check)

	// verify the will's claimable component's status is now claimed
	require.Equal(t, will_for_status_check.Components[0].Status, "claimed")
}

func TestKeeperClaimWithPedersenCommitment(t *testing.T) {
	kpr, ctx := setupKeeper(t) // Initialize your test environment
	creator := "creator-address"

	// Generate random scalars for value and blinding factor
	var valueScalar, blindingFactorScalar ristretto.Scalar
	valueScalar.Rand()
	blindingFactorScalar.Rand()

	var H ristretto.Point
	H.Rand()

	// Commit to the random value using the random blinding factor
	originalCommitment := pedersen.CommitTo(&H, &blindingFactorScalar, &valueScalar)

	addedCommitment := kpr.AddCommitments(originalCommitment, originalCommitment)

	// Create a will including the Pedersen commitment
	msg := &types.MsgCreateWillRequest{
		Creator:     creator,
		Name:        "Test Will",
		Beneficiary: "beneficiary-address",
		Height:      2,
		Components: []*types.ExecutionComponent{
			{
				Name:   "PedersenCommitmentComponent",
				Id:     "component-id",
				Status: "inactive",
				ComponentType: &types.ExecutionComponent_Claim{
					Claim: &types.ClaimComponent{
						SchemeType: &types.ClaimComponent_Pedersen{
							Pedersen: &types.PedersenCommitment{
								Commitment:       originalCommitment.Bytes(),
								TargetCommitment: addedCommitment.Bytes(),
							},
						},
					},
				},
			},
		},
	}

	will, err := kpr.CreateWill(sdk.UnwrapSDKContext(ctx), msg)
	require.NoError(t, err)
	require.NotNil(t, will)

	// Advance the block height to simulate time passage
	ctx = sdk.UnwrapSDKContext(ctx).WithBlockHeight(2)
	kpr.BeginBlocker(ctx)

	// Simulate a claim with the same commitment
	claimMsg := &types.MsgClaimRequest{
		WillId:      will.ID,
		Claimer:     creator,
		ComponentId: "component-id",
		ClaimType: &types.MsgClaimRequest_PedersenClaim{
			PedersenClaim: &types.PedersenClaim{
				Commitment: originalCommitment.Bytes(),
			},
		},
	}

	err = kpr.Claim(sdk.UnwrapSDKContext(ctx), claimMsg)
	require.NoError(t, err)

	updatedWill, err := kpr.GetWillByID(sdk.UnwrapSDKContext(ctx), will.ID)
	require.NoError(t, err)
	require.Equal(t, "claimed", updatedWill.Components[0].Status)
}

// Converts a string to a scalar value using SHA256 hash.
func stringToScalar(data string) ristretto.Scalar {
	var scalar ristretto.Scalar
	hash := sha256.Sum256([]byte(data)) // hash is a [32]byte array
	fmt.Println("String to scaler: ", data, "and hash: ", hash)

	bytes := scalar.SetBytes(&hash)
	fmt.Println("bytes: ", bytes, "and scalar: ", scalar)
	return scalar
}

// // Test with deterministic scalars derived from strings.
func TestKeeperClaimWithConstantPedersenCommitment(t *testing.T) {
	kpr, ctx := setupKeeper(t) // setupKeeper needs to be defined according to your context setup.
	creator := "will1p0k8gygawzpggzwftv7cv47zvgg8zaun5xucxz"

	// Define challenge and answer strings
	challengeString := "foo"
	answerString := "bar"

	// Convert strings to deterministic scalars
	challengeScalar := stringToScalar(challengeString)
	answerScalar := stringToScalar(answerString)
	fmt.Println("Challenge: ", challengeScalar, " Answer: ", answerScalar)

	var blindingFactorOriginal, blindingFactorClaim ristretto.Scalar
	blindingFactorOriginal.Rand()
	blindingFactorClaim.Rand()

	// Create Pedersen commitment using these scalars
	H := ristretto.Point{}
	H.Rand() // Base point for commitment (publicly known and constant)

	// Commitment from the creator (original commitment)
	originalCommitment := pedersen.CommitTo(&H, &blindingFactorOriginal, &challengeScalar)

	// For simplicity, using the same challenge-answer pair for the claim
	// In real scenarios, this might be a different pair or some operation based on the original
	claimCommitment := pedersen.CommitTo(&H, &blindingFactorClaim, &answerScalar)

	// will creator would have to add them ahead of time? or encrypt the blinding factor?
	addedCommitment := kpr.AddCommitments(originalCommitment, claimCommitment)

	// Create a will with Pedersen commitment component including a target commitment
	msg := &types.MsgCreateWillRequest{
		Creator:     creator,
		Name:        "Test Will",
		Beneficiary: "beneficiary-address",
		Height:      2,
		Components: []*types.ExecutionComponent{
			{
				Name:   "PedersenCommitmentComponent",
				Id:     "component-id",
				Status: "inactive",
				ComponentType: &types.ExecutionComponent_Claim{
					Claim: &types.ClaimComponent{
						SchemeType: &types.ClaimComponent_Pedersen{
							Pedersen: &types.PedersenCommitment{
								Commitment:       originalCommitment.Bytes(),
								TargetCommitment: addedCommitment.Bytes(),
							},
						},
					},
				},
			},
		},
	}

	// Simulated creation of the will
	will, err := kpr.CreateWill(sdk.UnwrapSDKContext(ctx), msg)
	require.NoError(t, err)
	require.NotNil(t, will)

	// Advance the block height to simulate time passage
	ctx = sdk.UnwrapSDKContext(ctx).WithBlockHeight(2)
	kpr.BeginBlocker(ctx)

	// Simulate the claiming process
	claimMsg := &types.MsgClaimRequest{
		WillId:      will.ID,
		Claimer:     creator,
		ComponentId: "component-id",
		ClaimType: &types.MsgClaimRequest_PedersenClaim{
			PedersenClaim: &types.PedersenClaim{
				Commitment: claimCommitment.Bytes(),
			},
		},
	}

	// Process the claim
	err = kpr.Claim(sdk.UnwrapSDKContext(ctx), claimMsg)
	require.NoError(t, err)

	// Verify the status of the commitment after processing the claim
	updatedWill, err := kpr.GetWillByID(sdk.UnwrapSDKContext(ctx), will.ID)
	require.NoError(t, err)
	require.Equal(t, "claimed", updatedWill.Components[0].Status)
}

///////////////////////////////////////////
///////////////////////////////////////////
///////////////////////////////////////////
///////////////////////////////////////////
///////////////////////////////////////////
///////////////////////////////////////////
///////////////////////////////////////////
///////////////////////////////////////////
///////////////////////////////////////////
///////////////////////////////////////////
///////////////////////////////////////////
///////////////////////////////////////////
///////////////////////////////////////////
///////////////////////////////////////////
///////////////////////////////////////////
///////////////////////////////////////////

// */

/*
func setupWithFundedAccount(t *testing.T, ctx sdk.Context, kpr *keeper.Keeper, addr sdk.AccAddress, coins sdk.Coins) {
	// Ensure the account exists
	fmt.Println("setupWithFundedAccount")
	fmt.Println("Address:", addr)
	fmt.Println("Coins:", coins)
	fmt.Println("prefix: ", sdk.GetConfig().GetBech32AccountAddrPrefix())

	acc := kpr.GetAccountKeeper().NewAccountWithAddress(ctx, addr)

	kpr.GetAccountKeeper().SetAccount(ctx, acc)

	// Fund the account
	err := kpr.GetBankKeeper().MintCoins(ctx, minttypes.ModuleName, coins)
	require.NoError(t, err)

	err = kpr.GetBankKeeper().SendCoinsFromModuleToAccount(ctx, minttypes.ModuleName, addr, coins)
	require.NoError(t, err)

	// Verify the balance
	balance := kpr.GetBankKeeper().GetBalance(ctx, addr, coins[0].Denom)
	require.True(t, balance.IsEqual(coins[0]))
}
*/

// func setupWithFundedAccount(t *testing.T, ctx sdk.Context, k *keeper.Keeper, addr sdk.AccAddress, coins sdk.Coins) {
func setupWithFundedAccount(t *testing.T, willchainApp *app.WasmApp, ctx sdk.Context, k *keeper.Keeper, addr sdk.AccAddress, coins sdk.Coins) {
	t.Helper()

	// Ensure the account exists
	fmt.Println("setupWithFundedAccount")
	fmt.Println("Address:", addr)
	fmt.Println("Coins:", coins)
	fmt.Println("prefix: ", sdk.GetConfig().GetBech32AccountAddrPrefix())
	fmt.Println("Address:", addr.String()) // Use .String() to get Bech32 format
	fmt.Println(ctx)                       // Use .String() to get Bech32 format
	fmt.Println(sdk.UnwrapSDKContext(ctx)) // Use .String() to get Bech32 format

	// Ensure the account keeper and bank keeper are not nil
	require.NotNil(t, k.GetAccountKeeper(), "accountKeeper is nil")
	require.NotNil(t, k.GetBankKeeper(), "bankKeeper is nil")

	// Ensure address is not empty
	require.NotEmpty(t, addr, "address is empty")

	// Validate coins
	require.NotEmpty(t, coins, "coins are empty")
	require.True(t, coins.IsValid(), "invalid coins")

	// Create and set account
	fmt.Println("k.GetAccountKeeper()")
	fmt.Println(k.GetAccountKeeper())
	fmt.Println("==================")
	fmt.Println(k.GetAccountKeeper().GetModuleAddress("will"))
	fmt.Println(k.GetAccountKeeper().Accounts)
	fmt.Println(k.GetAccountKeeper().AccountNumber)
	// for acc_i := range k.GetAccountKeeper().Accounts.Indexes.IndexesList() {
	// 	fmt.Println("acc_i")
	// 	fmt.Println(acc_i)
	// }
	for i, v := range k.GetAccountKeeper().Accounts.Indexes.IndexesList() {
		fmt.Println("acc_i")
		fmt.Println(i)
		fmt.Println(v)
	}

	// fmt.Println(ctx)
	fmt.Println("==================")

	// fmt.Println(k.GetAccountKeeper().GetAllAccounts(ctx))
	fmt.Println(k.GetAccountKeeper().GetAllAccounts(ctx))

	account_num := k.GetAccountKeeper().NextAccountNumber(ctx)
	fmt.Println("account_num: ", account_num)

	acc := k.GetAccountKeeper().NewAccountWithAddress(ctx, addr)
	require.NotNil(t, acc, "failed to create new account")
	k.GetAccountKeeper().SetAccount(ctx, acc)

	// Fund account
	err := k.GetBankKeeper().MintCoins(ctx, minttypes.ModuleName, coins)
	require.NoError(t, err, "failed to mint coins")
	err = k.GetBankKeeper().SendCoinsFromModuleToAccount(ctx, minttypes.ModuleName, addr, coins)
	require.NoError(t, err, "failed to send coins from module to account")

	// Verify the account's balance
	balance := k.GetBankKeeper().GetBalance(ctx, addr, coins.GetDenomByIndex(0))
	require.Equal(t, coins.AmountOf(balance.Denom), balance.Amount, "balance does not match expected amount")
}

/*
ctx
{{{}} 0x140001c6a00 {{0 0}  0 {0 0 <nil>} {[] {0 []}} [] [] [] [] [] [] [] [] []} []  [] {} [] 0x140031d11b0 <nil> true false 0 [] {<nil> <nil> <nil> <nil> <nil>} 0x140031b34b8 0 {1000 1000 1000 3 2000 30 30} {100 100 100 0 200 3 3} {[] false} <nil> {0 [] {0 0 <nil>}  []}}
{{{}} 0x140001c6a00 {{0 0}  0 {0 0 <nil>} {[] {0 []}} [] [] [] [] [] [] [] [] []} []  [] {} [] 0x140031d11b0 <nil> true false 0 [] {<nil> <nil> <nil> <nil> <nil>} 0x140031b34b8 0 {1000 1000 1000 3 2000 30 30} {100 100 100 0 200 3 3} {[] false} <nil> {0 [] {0 0 <nil>}  []}}
*/
// func TestExecuteTransfer(t *testing.T) {
// 	// Setup the keeper and context
// 	willchainApp := app.Setup(t)

// 	kpr, ctx := setupKeeper(t)
// 	// kpr := &w3llApp.WillKeeper
// 	// ctx := w3llApp.NewUncachedContext(true, tmproto.Header{})

// 	// Create test accounts
// 	// fromAddr := sdk.AccAddress([]byte("w3ll1p0k8gygawzpggzwftv7cv47zvgg8zaun5xucxz"))
// 	// toAddr := sdk.AccAddress([]byte("w3ll1p0k8gygawzpggzwftv7cv47zvgg8zaun5xucxz"))
// 	var err error
// 	// Define the addresses in Bech32 format
// 	// fromAddrStr := "w3ll1p0k8gygawzpggzwftv7cv47zvgg8zaun5xucxz"
// 	// toAddrStr := "w3ll1p0k8gygawzpggzwftv7cv47zvgg8zaun5xucxz"
// 	// fromAddrStr := "cosmos1p0k8gygawzpggzwftv7cv47zvgg8zaunupyu5j"
// 	// toAddrStr := "cosmos1p0k8gygawzpggzwftv7cv47zvgg8zaunupyu5j"
// 	fromAddrStr := "cosmos146s93ufqqq9p909ykanau83p8zdct2tgnlzmvv"
// 	toAddrStr := "cosmos1fl48vsnmsdzcv85q5d2q4z5ajdha8yu34mf0eh"

// 	// Convert string addresses to sdk.AccAddress
// 	fromAddr, err := sdk.AccAddressFromBech32(fromAddrStr)
// 	require.NoError(t, err, "Invalid 'from' address")
// 	toAddr, err := sdk.AccAddressFromBech32(toAddrStr)
// 	require.NoError(t, err, "Invalid 'to' address")

// 	// sdk.AccAddress{}
// 	// Define the amount to transfer
// 	// amount := sdk.NewCoins(sdk.NewInt64Coin("w3ll", 100))
// 	// amount := sdk.NewInt64Coin("w3ll", 100)
// 	transferAmount := sdk.NewInt64Coin("atom", 100) // This is the amount to be transferred

// 	// fundAmount := sdk.NewCoins(sdk.NewInt64Coin("w3ll", 1000)) // This is the amount to fund the account with

// 	// Ensure the "from" account exists and fund it
// 	// setupWithFundedAccount(t, ctx, kpr, fromAddr, fundAmount)

// 	// Define a nominal fund amount for the "to" account
// 	nominalFundAmount := sdk.NewCoins(sdk.NewInt64Coin("atom", 1))

// 	/////////???????????/
// 	// fmt.Println(w3llApp.AccountKeeper.GetAllAccounts(ctx))

// 	// panic(1)
// 	// Ensure the "to" account exists (funding not necessary but account must exist)
// 	// setupWithFundedAccount(t, ctx, kpr, toAddr, sdk.NewCoins())
// 	// setupWithFundedAccount(t, sdk.UnwrapSDKContext(ctx), kpr, toAddr, nominalFundAmount)
// 	setupWithFundedAccount(t, willchainApp, ctx, kpr, toAddr, nominalFundAmount)
// 	// setupWithFundedAccount(t, kpr, toAddr, nominalFundAmount)

// 	// Prepare the transfer component
// 	component := types.ExecutionComponent{
// 		ComponentType: &types.ExecutionComponent_Transfer{
// 			Transfer: &types.TransferComponent{
// 				From:   fromAddr.String(),
// 				To:     toAddr.String(),
// 				Amount: &transferAmount,
// 			},
// 		},
// 	}

// 	// 	fmt.Println("component")
// 	// 	fmt.Printf("component: %+v\n", component)
// 	// 	fmt.Println(component.ComponentType)

// 	// Create a dummy will object (if needed)
// 	will := types.Will{
// 		ID:          "will-1",
// 		Creator:     fromAddr.String(),
// 		Name:        "test will",
// 		Beneficiary: toAddr.String(),
// 		Height:      2,
// 	}

// 	// Execute the transfer
// 	err_transfer := kpr.ExecuteTransfer(ctx, &component, will)
// 	require.NoError(t, err_transfer, "ExecuteTransfer should not return an error")

// 	// Verify the transfer by checking balances
// 	fromBalance := kpr.GetBankKeeper().GetBalance(ctx, fromAddr, "w3ll")
// 	toBalance := kpr.GetBankKeeper().GetBalance(ctx, toAddr, "w3ll")

// 	require.Equal(t, int64(0), fromBalance.Amount.Int64(), "from account balance should decrease by the transfer amount")
// 	require.Equal(t, int64(100), toBalance.Amount.Int64(), "to account balance should increase by the transfer amount")

// 	// Optionally, verify events emitted during the transfer
// 	events := ctx.EventManager().Events()
// 	transferEventFound := false
// 	for _, event := range events {
// 		if event.Type == "transfer" { // Adjust event type as necessary
// 			transferEventFound = true
// 			// Further checks on event attributes can be done here
// 			break
// 		}
// 	}
// 	require.True(t, transferEventFound, "transfer event should be emitted")
// }

// func TestExecuteTransfer(t *testing.T) {
// 	// Setup the keeper and context
// 	kpr, ctx := setupKeeper(t)

// 	// Create test accounts and fund the sender account
// 	fromAddr := sdk.AccAddress([]byte("w3ll1p0k8gygawzpggzwftv7cv47zvgg8zaun5xucxz"))
// 	toAddr := sdk.AccAddress([]byte("w3ll1p0k8gygawzpggzwftv7cv47zvgg8zaun5xucxz"))
// 	// amount := sdk.NewCoins(sdk.NewInt64Coin("w3ll", 100))
// 	amount := sdk.NewInt64Coin("w3ll", 100)

// 	// Fund the sender account
// 	// err := kpr.GetBankKeeper().AddCoins(ctx, fromAddr, sdk.NewCoins(sdk.NewInt64Coin("w3ll", 1000)))
// 	// require.NoError(t, err)

// 	// Prepare the transfer component
// 	component := types.ExecutionComponent{
// 		ComponentType: &types.ExecutionComponent_Transfer{
// 			Transfer: &types.TransferComponent{

// 				// From:   fromAddr.String(),
// 				// To:     toAddr.String(),

// 				From: "w3ll1p0k8gygawzpggzwftv7cv47zvgg8zaun5xucxz",
// 				To:   "w3ll1p0k8gygawzpggzwftv7cv47zvgg8zaun5xucxz",

// 				Amount: &amount,
// 			},
// 		},
// 	}
// 	fmt.Println("component")
// 	fmt.Printf("component: %+v\n", component)
// 	fmt.Println(component.ComponentType)

// 	// Create a dummy will object (if needed)
// 	will := types.Will{
// 		ID:          "will-1",
// 		Creator:     "w3ll1p0k8gygawzpggzwftv7cv47zvgg8zaun5xucxz",
// 		Name:        "test will",
// 		Beneficiary: "w3ll1p0k8gygawzpggzwftv7cv47zvgg8zaun5xucxz",
// 		Height:      2,
// 	}

// 	// Execute the transfer
// 	err := kpr.ExecuteTransfer(ctx, &component, will)
// 	require.NoError(t, err, "ExecuteTransfer should not return an error")

// 	// Verify the transfer by checking balances
// 	fromBalance := kpr.GetBankKeeper().GetBalance(ctx, fromAddr, "w3ll")
// 	toBalance := kpr.GetBankKeeper().GetBalance(ctx, toAddr, "w3ll")

// 	require.Equal(t, int64(900), fromBalance.Amount.Int64(), "from account balance should decrease by the transfer amount")
// 	require.Equal(t, int64(100), toBalance.Amount.Int64(), "to account balance should increase by the transfer amount")

// 	// Optionally, verify events emitted during the transfer
// 	events := ctx.EventManager().Events()
// 	transferEventFound := false
// 	for _, event := range events {
// 		if event.Type == "transfer" { // Adjust event type as necessary
// 			transferEventFound = true
// 			// Further checks on event attributes can be done here
// 			break
// 		}
// 	}
// 	require.True(t, transferEventFound, "transfer event should be emitted")
// }

// func TestKeeperContractCall(t *testing.T) {
// 	kpr, ctx := setupKeeper(t)

// 	// Mock data for the test
// 	contractAddr := "cosmos1contractaddress" // Mock contract address
// 	callerAddr := "cosmos1calleraddress"     // Mock caller address
// 	contractData := []byte("contract execution data")
// 	// coinInt := sdl//sdk.NewInt(100)
// 	coinInt := math.NewInt(100)
// 	var coin sdk.Coin = sdk.NewCoin("testcoin", coinInt) // Mock execution data
// 	coins := sdk.NewCoins(coin)                          // Mock coins to send with the contract call
// 	fmt.Println("coins")
// 	fmt.Println(coins)
// 	// Convert the address strings to sdk.AccAddress
// 	contractSDKAddr, err := sdk.AccAddressFromBech32(contractAddr)
// 	require.NoError(t, err, "invalid contract address")
// 	fmt.Println("contractSDKAddr")
// 	fmt.Println(contractSDKAddr)

// 	callerSDKAddr, err := sdk.AccAddressFromBech32(callerAddr)
// 	fmt.Println(callerSDKAddr)
// 	fmt.Println("callerSDKAddr")
// 	require.NoError(t, err, "invalid caller address")

// 	// Assuming the contract exists and the keeper is correctly set up to call it.
// 	// You might need to mock or pre-set the state of your keeper to simulate a real contract existing at `contractAddr`.

// 	// Simulate calling the contract
// 	result, err := kpr.ExecuteContract(ctx, &types.ExecutionComponent_Contract{
// 		Contract: &types.ContractComponent{
// 			Address: contractAddr,
// 			Data:    contractData,
// 		},
// 	})

// 	// Validate the execution result
// 	require.NoError(t, err, "ExecuteContract should not return an error")
// 	assert.NotNil(t, result, "ExecuteContract should return a result")

// 	// Additional assertions depending on what the contract execution is expected to do or return
// 	// For example, if you expect a specific state change, query the state and assert the expected changes
// 	// If the contract execution returns a result, you can decode and assert specific values in the result
// }

// func TestSendIBCMessage(t *testing.T) {
// 	kpr, ctx := setupKeeper(t)

// 	// Define test data
// 	channelID := "testChannelID"
// 	portID := "testPortID"
// 	data := []byte("testData")
// 	fmt.Println(ctx)
// 	fmt.Println(sdk.UnwrapSDKContext(ctx))
// 	// Assume we have a method in our keeper to abstract the sending and sequence management
// 	err := kpr.SendIBCMessage(sdk.UnwrapSDKContext(ctx), channelID, portID, data)
// 	require.NoError(t, err, "Sending IBC message should not result in an error")

// 	events := ctx.EventManager().ABCIEvents()
// 	require.Len(t, events, 1, "Expected one event to be emitted")
// 	event := events[0]
// 	require.Equal(t, event.Type, "ibc_message_sent", "Expected event type 'ibc_message_sent'")
// 	require.Equal(t, event.Attributes[0].Key, "channel_id", "Expected attribute 'channel_id'")
// 	require.Equal(t, event.Attributes[0].Value, channelID, "Expected channel ID matches")
// 	require.Equal(t, event.Attributes[1].Key, "port_id", "Expected attribute 'port_id'")
// 	require.Equal(t, event.Attributes[1].Value, portID, "Expected port ID matches")
// }

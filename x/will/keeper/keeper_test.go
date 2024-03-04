package keeper_test

import (
	// "context"

	"encoding/hex"
	"fmt"
	"testing"

	// _proto "github.com/cosmos/gogoproto/proto"
	"github.com/stretchr/testify/assert"
	// "github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	// "google.golang.org/protobuf/proto"

	// "cosmossdk.io/core/store"
	// corestoretypes "cosmossdk.io/core/store"
	"cosmossdk.io/log"
	corestore "cosmossdk.io/store"

	// corestore "github.com/cosmos/cosmos-sdk/store"

	// corestoretypes "cosmossdk.io/core/store"
	storemetrics "cosmossdk.io/store/metrics"
	storetypes "cosmossdk.io/store/types"

	// codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/runtime"
	sdk "github.com/cosmos/cosmos-sdk/types"

	// Import the tm-db package
	// dbm "github.com/tendermint/tm-db" // Import the tm-db package
	dbm "github.com/cosmos/cosmos-db"

	// tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
	"github.com/CosmWasm/wasmd/app"
	"github.com/CosmWasm/wasmd/x/will/keeper"
	"github.com/CosmWasm/wasmd/x/will/types"
	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"
)

func setupKeeper(t *testing.T) (*keeper.Keeper, sdk.Context) {
	w3llApp := app.Setup(t)
	mockedCodec := w3llApp.AppCodec()

	// Initialize DB and store
	memDB := dbm.NewMemDB()
	// ms := corestore.NewCommitMultiStore(memDB) // Initialize the MultiStore with the in-memory DB
	ms := corestore.NewCommitMultiStore(memDB, log.NewTestLogger(t), storemetrics.NewNoOpMetrics())
	// keyWill := storetypes.NewKVStoreKey(types.StoreKey)
	storeKey := storetypes.NewKVStoreKey(types.StoreKey)
	storeservice := runtime.NewKVStoreService(storeKey)

	// Create and mount store keys
	// ms.MountStoreWithDB(keyWill, storetypes.StoreTypeIAVL, memDB)
	ms.MountStoreWithDB(storeKey, storetypes.StoreTypeIAVL, memDB)
	require.NoError(t, ms.LoadLatestVersion())

	// Create context
	ctx := sdk.NewContext(ms, tmproto.Header{}, false, log.NewNopLogger())

	// Initialize keeper with the store key
	// k := keeper.NewKeeper(mockedCodec, storeKey, nil)
	k := keeper.NewKeeper(mockedCodec, storeservice, nil)

	return &k, ctx
}

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

	// If you have specific expectations for the Components, verify those as well
	// This example assumes you want to check the length of the components slice
	assert.Len(t, will.Components, len(msg.Components), "number of will components should match the request")
}

func TestKeeperClaimWithSchnorrSignature(t *testing.T) {
	kpr, ctx := setupKeeper(t)

	// Hardcoded values from your Schnorr signature example
	// publicKeyHex := "2320a2da28561875cedbb0c25ae458e0a1d087834ae49b96a3f93cec79a8190c"
	// signatureRHex := "7ab0edb9b0929b5bb4b47dfb927d071ecc5de75985662032bb52ef3c5ace640b"
	// signatureSHex := "165c6df5ea8911a6c0195a3140be5119a5b882e91b34cbcdd31ef3f5b0035b06"

	// v2
	// publicKeyHex := "6f2de2f173efcbd7fc1fdec2d2939040575a248759d6d2373eaf775b1eef3a6e"
	// signatureRHex := "4143f859db4b5fd2e97aea3c332eb78497d4785cdfe682d9954036ab9f63fc34"
	// signatureSHex := "22c819840897cffc4936ce576e21c4bc6712bd8475b371c7897ea267b83d180e"

	// v3
	publicKeyHex := "6f2de2f173efcbd7fc1fdec2d2939040575a248759d6d2373eaf775b1eef3a6e"
	signatureRHex := "4143f859db4b5fd2e97aea3c332eb78497d4785cdfe682d9954036ab9f63fc34"
	signatureSHex := "22c819840897cffc4936ce576e21c4bc6712bd8475b371c7897ea267b83d180e"

	message := "message-2b-signed"
	creator := "creator-address"
	// Convert hexadecimal strings to bytes
	publicKeyBytes, err := hex.DecodeString(publicKeyHex)
	require.NoError(t, err)
	signatureRBytes, err := hex.DecodeString(signatureRHex)
	require.NoError(t, err)
	signatureSBytes, err := hex.DecodeString(signatureSHex)
	require.NoError(t, err)
	// messageBytes := []byte(message)
	fmt.Println("===DEBUG====")
	fmt.Println(publicKeyHex)
	fmt.Println(publicKeyBytes)
	fmt.Println(signatureRHex)
	fmt.Println(signatureRBytes)
	fmt.Println(signatureSHex)
	fmt.Println(signatureSBytes)
	// Assuming the signature is the concatenation of R and S components
	signatureBytes := append(signatureRBytes, signatureSBytes...)

	msg := &types.MsgCreateWillRequest{
		Creator:     creator,
		Name:        "Test Will",
		Beneficiary: "beneficiary-address",
		Height:      2,
		Components: []*types.ExecutionComponent{
			{
				Name: "SchnorrSignatureComponent",
				ComponentType: &types.ExecutionComponent_Claim{
					Claim: &types.ClaimComponent{
						SchemeType: &types.ClaimComponent_Schnorr{
							Schnorr: &types.SchnorrSignature{
								PublicKey: publicKeyBytes,
								Signature: signatureBytes,
								Message:   message,
							},
						},
					},
				},
			},
		},
	}

	// Assuming CreateWill method is correctly implemented to handle mocks
	will, err := kpr.CreateWill(sdk.UnwrapSDKContext(ctx), msg)
	require.NoError(t, err)
	assert.NotNil(t, will)

	fmt.Println("====[TEST]HERE IS THE PRACTICE WILL[TEST]=====")
	fmt.Println((will))

	// Assuming a will has already been created and you have its ID
	willID := will.ID // Replace with the actual will ID
	componentID := will.Components[0].Id

	// roll height forward
	ctx_future := sdk.UnwrapSDKContext(ctx).WithBlockHeight(2)

	// run the begin blocker with the updated block height for the will to execute.
	kpr.BeginBlocker(ctx_future)

	// Construct the claim request with the Schnorr claim
	claimMsg := &types.MsgClaimRequest{
		WillId:      willID,
		Claimer:     creator,
		ComponentId: componentID,
		ClaimType: &types.MsgClaimRequest_SchnorrClaim{
			SchnorrClaim: &types.SchnorrClaim{
				PublicKey: publicKeyBytes,
				Signature: signatureBytes,
				Message:   []byte(message),
			},
		},
	}

	// Process the claim
	err = kpr.Claim(sdk.UnwrapSDKContext(ctx_future), claimMsg)
	require.NoError(t, err, "processing Schnorr claim should not produce an error")

}

package keeper_test

import (
	"fmt"
	"testing"

	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"
	// Import the tm-db package
	// dbm "github.com/tendermint/tm-db" // Import the tm-db package
	dbm "github.com/cosmos/cosmos-db"
	// _proto "github.com/cosmos/gogoproto/proto"
	"github.com/stretchr/testify/assert"
	// "github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	// "cosmossdk.io/core/store"
	// corestoretypes "cosmossdk.io/core/store"
	"cosmossdk.io/log"
	corestore "cosmossdk.io/store"

	// corestoretypes "cosmossdk.io/core/store"
	storemetrics "cosmossdk.io/store/metrics"
	storetypes "cosmossdk.io/store/types"

	// codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/runtime"
	sdk "github.com/cosmos/cosmos-sdk/types"

	// tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
	"github.com/CosmWasm/wasmd/app"
	"github.com/CosmWasm/wasmd/x/will/keeper"
	"github.com/CosmWasm/wasmd/x/will/types"
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
								Signature: []byte(signatureHex),
								Message:   message,
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

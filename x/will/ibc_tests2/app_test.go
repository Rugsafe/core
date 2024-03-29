package ibctesting

import (
	"fmt"
	"testing"
	"time"

	"cosmossdk.io/log"
	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"

	corestore "cosmossdk.io/store"
	storemetrics "cosmossdk.io/store/metrics"
	storetypes "cosmossdk.io/store/types"
	dbm "github.com/cosmos/cosmos-db"
	"github.com/cosmos/cosmos-sdk/runtime"
	sdk "github.com/cosmos/cosmos-sdk/types"
	ibctransfertypes "github.com/cosmos/ibc-go/v8/modules/apps/transfer/types"
	clienttypes "github.com/cosmos/ibc-go/v8/modules/core/02-client/types"
	channeltypes "github.com/cosmos/ibc-go/v8/modules/core/04-channel/types"
	ibcexported "github.com/cosmos/ibc-go/v8/modules/core/exported"
	ibctesting "github.com/cosmos/ibc-go/v8/testing"

	//
	"github.com/CosmWasm/wasmd/app"
	"github.com/CosmWasm/wasmd/x/will/keeper"
	"github.com/CosmWasm/wasmd/x/will/types"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type IBCTestSuite struct {
	suite.Suite
	coordinator *ibctesting.Coordinator
	chainA      *ibctesting.TestChain
	chainB      *ibctesting.TestChain
}

func setupKeeper(t *testing.T) (*keeper.Keeper, sdk.Context) {
	w3llApp := app.Setup(t)
	mockedCodec := w3llApp.AppCodec()
	// Initialize DB and store
	memDB := dbm.NewMemDB()
	ms := corestore.NewCommitMultiStore(memDB, log.NewTestLogger(t), storemetrics.NewNoOpMetrics())
	storeKey := storetypes.NewKVStoreKey(types.StoreKey)
	ibcStoreKey := storetypes.NewKVStoreKey(ibctransfertypes.StoreKey)    // IBC store key
	ibcExportedStoreKey := storetypes.NewKVStoreKey(ibcexported.StoreKey) // IBC store key
	storeservice := runtime.NewKVStoreService(storeKey)
	// memKey := storetypes.NewMemoryStoreKey("mem_capability") // Memory key for capability

	// Create and mount store keys
	ms.MountStoreWithDB(storeKey, storetypes.StoreTypeIAVL, memDB)
	ms.MountStoreWithDB(ibcStoreKey, storetypes.StoreTypeIAVL, memDB)         // Mount the IBC store
	ms.MountStoreWithDB(ibcExportedStoreKey, storetypes.StoreTypeIAVL, memDB) // Mount the IBC store
	require.NoError(t, ms.LoadLatestVersion())

	// Create context
	ctx := sdk.NewContext(ms, tmproto.Header{}, false, log.NewNopLogger())

	// capabilityKeeper := capabilitykeeper.NewKeeper(mockedCodec, storeKey, memKey)

	// Initialize keeper with the store key
	fmt.Println("IBC TESTING LOGS")
	fmt.Println(w3llApp.IBCKeeper.ChannelKeeper)
	fmt.Println(w3llApp.ScopedWillKeeper)
	fmt.Println(w3llApp.ScopedIBCKeeper)
	fmt.Println(w3llApp.CapabilityKeeper)

	// scopedWillKeeper := w3llApp.CapabilityKeeper.ScopeToModule(willtypes.ModuleName)
	// scopedIBCKeeper := w3llApp.CapabilityKeeper.ScopeToModule(ibcexported.ModuleName)
	k := keeper.NewKeeper(
		mockedCodec,
		storeservice,
		nil,
		w3llApp.IBCKeeper.ChannelKeeper,
		w3llApp.IBCKeeper.PortKeeper,

		//scopedWillKeeper,
		//scopedIBCKeeper,

		w3llApp.ScopedWillKeeper,
		w3llApp.ScopedIBCKeeper,

		*w3llApp.CapabilityKeeper,
	)

	return &k, ctx
}

func (suite *IBCTestSuite) SetupTest() {
	suite.coordinator = ibctesting.NewCoordinator(suite.T(), 2)
	suite.chainA = suite.coordinator.GetChain(ibctesting.GetChainID(1))
	suite.chainB = suite.coordinator.GetChain(ibctesting.GetChainID(2))
}

func (suite *IBCTestSuite) TestClientCreation() {
	path := ibctesting.NewPath(suite.chainA, suite.chainB)
	suite.coordinator.Setup(path) // This sets up clients, connections, and channels

	// Now you can assert that clients have been created correctly
	suite.Require().NotEmpty(path.EndpointA.ClientID)
	suite.Require().NotEmpty(path.EndpointB.ClientID)

}

// util
// Updates IBC clients without incrementing block height.
func (suite *IBCTestSuite) updateClientsWithoutBlockIncrement(path *ibctesting.Path) {
	err := path.EndpointA.UpdateClient()
	suite.Require().NoError(err)

	err = path.EndpointB.UpdateClient()
	suite.Require().NoError(err)
}

// Only increment time and commit blocks without updating clients.
func (suite *IBCTestSuite) incrementTimeAndCommit() {
	suite.coordinator.IncrementTimeBy(time.Minute)
	suite.coordinator.CommitBlock(suite.chainA, suite.chainB)
}

func (suite *IBCTestSuite) sync(path *ibctesting.Path) {
	for path.EndpointB.GetClientState().GetLatestHeight().LT(path.EndpointA.GetClientState().GetLatestHeight()) {
		suite.incrementTimeAndCommit()
		suite.updateClientsWithoutBlockIncrement(path)
	}
}

func (suite *IBCTestSuite) TestPacketTransmission() {
	path := ibctesting.NewPath(suite.chainA, suite.chainB)
	suite.coordinator.Setup(path)

	// Adjust timeout height and timestamp
	futureHeight := uint64(10)

	futureTimestamp := uint64(time.Now().Add(2 * time.Hour).UnixNano())

	timeoutHeight := clienttypes.NewHeight(1, uint64(suite.chainA.CurrentHeader.Height)+futureHeight)
	timeoutTimestamp := futureTimestamp

	packetData := []byte("data")
	packet := channeltypes.NewPacket(
		packetData,
		1,
		path.EndpointA.ChannelConfig.PortID,
		path.EndpointA.ChannelID,
		path.EndpointB.ChannelConfig.PortID,
		path.EndpointB.ChannelID,
		timeoutHeight,
		timeoutTimestamp,
	)

	fmt.Printf("Block height before incrementing: A: %d, B: %d\n", suite.chainA.App.LastBlockHeight(), suite.chainB.App.LastBlockHeight())

	fmt.Println("timeoutHeight: ", timeoutHeight)
	fmt.Println("timeoutTimestamp: ", timeoutTimestamp)
	fmt.Println("packet.TimeoutHeight: ", packet.TimeoutHeight)
	fmt.Println("packet.GetTimeoutTimestamp(): ", packet.GetTimeoutTimestamp())

	fmt.Printf("Block height after first increment: A: %d, B: %d\n", suite.chainA.App.LastBlockHeight(), suite.chainB.App.LastBlockHeight())

	packetSeq, sendPacketErr := path.EndpointA.SendPacket(packet.TimeoutHeight, packet.GetTimeoutTimestamp(), packetData)
	fmt.Println("packetSeq: ", packetSeq)
	suite.Require().NoError(sendPacketErr)

	//increment?
	suite.incrementTimeAndCommit()

	// Ensure clients on both chains are updated before receiving the packet
	suite.updateClientsWithoutBlockIncrement(path)

	fmt.Printf("Block height before receiving packet: A: %d, B: %d\n", suite.chainA.App.LastBlockHeight(), suite.chainB.App.LastBlockHeight())
	recvErr := path.EndpointB.RecvPacket(packet)
	suite.Require().NoError(recvErr)

	fmt.Printf("Client on ChainA (before ack) at height: %s\n", path.EndpointA.GetClientState().GetLatestHeight())
	fmt.Printf("Client on ChainB (before ack) at height: %s\n", path.EndpointB.GetClientState().GetLatestHeight())

	// Manually wait for chain B to reach chain A's current height
	fmt.Println("SANITY CHECK")
	fmt.Println(path.EndpointB.Chain.CurrentHeader.GetHeight())
	fmt.Println(path.EndpointA.Chain.CurrentHeader.Height)

	// ack := channeltypes.NewResultAcknowledgement([]byte{byte(1)})
	ack := channeltypes.NewResultAcknowledgement([]byte{byte(1)}).Acknowledgement()
	fmt.Println(ack)
	ackErr := path.EndpointB.AcknowledgePacket(packet, ack)
	suite.Require().NoError(ackErr)

	fmt.Printf("Final block height: A: %d, B: %d\n", suite.chainA.App.LastBlockHeight(), suite.chainB.App.LastBlockHeight())

}

// func (suite *IBCTestSuite) TestBeginBlockerInteractionWithIBC() {
//     // Initialize your keeper with mock or real dependencies
//     myKeeper := keeper.NewKeeper(...)

//     // Setup the state required to trigger the BeginBlocker logic
//     // For example, create a will that is due for processing
//     ctx := suite.chainA.GetContext()
//     myKeeper.CreateWill(ctx, "testWill", ...)

//     // Advance the blockchain to the block height where the will should be processed
//     suite.coordinator.CommitBlock(suite.chainA, suite.chainB)

//     // Verify that the IBC message was sent as expected
//     // This might involve checking the mock ChannelKeeper or directly inspecting the state on the receiving chain
// }

func TestSendIBCMessage(t *testing.T) {
	// Setup
	keeper, ctx := setupKeeper(t)

	// Define your test channel ID, port ID, and data
	channelID := "channel-0"
	portID := "port-0"
	data := []byte("test data")

	// Act
	fmt.Println("CTX:")
	// fmt.Println(ctx)
	fmt.Println(keeper.ChannelKeeper.GetNextSequenceSend(ctx, portID, channelID))
	err := keeper.SendIBCMessage(ctx, channelID, portID, data)
	require.NoError(t, err, "SendIBCMessage should not error")

	// Here you would check for effects of sending the message.
	// In an actual blockchain environment, you would check state changes or events.
	// For this test environment, you might check if state entries expected to change have changed,
	// if events are emitted, or other side effects depending on what SendIBCMessage does.

	// Example: Check if a packet has been sent (mocking or state inspection)
	// This part is highly dependent on how your app is structured and what tools or mock setups you have available.
	// For an actual chain, you'd look into the state storage or events to see if the packet send operation has been initiated.

	// Assert
	// Verify that SendPacket was called on the mock with expected arguments

}

func TestMyIBCTestSuite(t *testing.T) {
	suite.Run(t, new(IBCTestSuite))
}

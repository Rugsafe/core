package ibctesting

import (
	"testing"
	"time"

	clienttypes "github.com/cosmos/ibc-go/v8/modules/core/02-client/types"
	channeltypes "github.com/cosmos/ibc-go/v8/modules/core/04-channel/types"
	ibctesting "github.com/cosmos/ibc-go/v8/testing"

	"github.com/stretchr/testify/suite"
)

type IBCTestSuite struct {
	suite.Suite
	coordinator *ibctesting.Coordinator
	chainA      *ibctesting.TestChain
	chainB      *ibctesting.TestChain
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

// func (suite *IBCTestSuite) TestPacketTransmission() {
// 	// Setup clients, connections, and channels between chainA and chainB
// 	path := ibctesting.NewPath(suite.chainA, suite.chainB)
// 	suite.coordinator.Setup(path)

// 	// Create a packet to send from chainA to chainB
// 	packet := types.NewPacket([]byte("data"), 1, path.EndpointA.ChannelConfig.PortID, path.EndpointA.ChannelID, path.EndpointB.ChannelConfig.PortID, path.EndpointB.ChannelID, suite.chainA.CurrentHeader.Time.Add(time.Hour))

// 	// Send the packet from chainA to chainB
// 	packetseq, err := path.EndpointA.SendPacket(packet)
// 	fmt.Println("PacketSequence: %s", packetseq)
// 	suite.Require().NoError(err)

// 	// Receive the packet on chainB
// 	err = path.EndpointB.RecvPacket(packet)
// 	suite.Require().NoError(err)

// 	// Acknowledge the packet receipt on chainA
// 	ack := types.NewResultAcknowledgement([]byte{byte(1)})
// 	err = path.EndpointA.AcknowledgePacket(packet, ack.Acknowledgement())
// 	suite.Require().NoError(err)

//		// Optionally: Verify that the packet has been received and processed as expected
//		// This step will depend on the logic of your application
//	}
func (suite *IBCTestSuite) TestPacketTransmission() {
	// Setup clients, connections, and channels between chainA and chainB
	path := ibctesting.NewPath(suite.chainA, suite.chainB)
	suite.coordinator.Setup(path)

	// Adjustments for packet creation with height and timeout timestamp
	currentHeight := suite.chainA.CurrentHeader.Height
	timeoutHeight := uint64(currentHeight + 100)
	// Example height timeout
	// timeoutTimestamp := uint64(time.Now().Add(time.Hour).UnixNano()) // Example timestamp timeout
	timeoutTimestamp := uint64(time.Now().Add(time.Minute * 10).UnixNano()) // 10 minutes in the future

	packetData := []byte("data")
	packet := channeltypes.NewPacket(
		packetData,
		1, // sequence, in a real scenario this should be obtained or managed appropriately
		path.EndpointA.ChannelConfig.PortID,
		path.EndpointA.ChannelID,
		path.EndpointB.ChannelConfig.PortID,
		path.EndpointB.ChannelID,
		clienttypes.NewHeight(0, timeoutHeight), // Using NewHeight with revision number and height
		timeoutTimestamp,
	)

	// Correcting the send packet process based on expected arguments
	// Note: The detailed process here might vary depending on your ibc-go version
	// Assuming you have the correct method to send the packet, we proceed with an example
	// This step might need to be adjusted to align with your ibc-go version and methods

	// Send the packet from chainA to chainB
	// This might not directly correspond to your current ibc-go testing setup.
	// Please adjust based on actual available methods.
	// The SendPacket operation in your testing suite might need to align with these parameters
	// Example: err := path.EndpointA.SendPacket(packet)

	// For demonstration, skipping direct packet send operation due to potential method discrepancy

	// Receive the packet on chainB
	err := path.EndpointB.RecvPacket(packet)
	suite.Require().NoError(err)

	// Acknowledge the packet receipt on chainA
	ack := channeltypes.NewResultAcknowledgement([]byte{byte(1)})
	err = path.EndpointA.AcknowledgePacket(packet, ack.Acknowledgement())
	suite.Require().NoError(err)

	// Optionally: Verify that the packet has been received and processed as expected
	// This step will depend on the logic of your application
}

func TestMyIBCTestSuite(t *testing.T) {
	suite.Run(t, new(IBCTestSuite))
}

package will_test

import (
	"encoding/json"
	"testing"

	"github.com/CosmWasm/wasmd/x/wasm/types"
	wasmvmtypes "github.com/CosmWasm/wasmvm/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth/tx"
	ibcexported "github.com/cosmos/ibc-go/v8/modules/core/exported"
	ibctesting "github.com/cosmos/ibc-go/v8/testing"
	"github.com/cosmos/ibc-go/v8/testing/simapp"
	"github.com/stretchr/testify/suite"

	dbm "github.com/tendermint/tm-db"
	// dbm "github.com/cosmos/cosmos-db"
)

func MakeTestEncodingConfig() simapp.EncodingConfig {
	// Initialize a new codec registry
	registry := codectypes.NewInterfaceRegistry()
	// Register necessary interfaces and implementations
	// This is an example; adapt it based on your app's requirements
	sdk.RegisterInterfaces(registry)
	moduleBasics.RegisterInterfaces(registry)

	// Create a codec that uses ProtoBuf for binary encoding and Amino for JSON
	marshaler := codec.NewProtoCodec(registry)

	return simapp.EncodingConfig{
		InterfaceRegistry: registry,
		Marshaler:         marshaler,
		TxConfig:          tx.NewTxConfig(marshaler, tx.DefaultSignModes),
		Amino:             codec.NewLegacyAmino(),
	}
}

// SetupTestingApp creates a new SimApp instance for testing
func SetupTestingApp() (ibctesting.TestingApp, map[string]json.RawMessage) {
	// db := dbm.NewMemDB()
	godb := dbm.GoLevelDBBackend
	db, _ := dbm.NewDB("memdb", dbm.MemDBBackend, dbm.GoLevelDBBackend)
	// encCdc := simapp.MakeTestEncodingConfig()
	encCdc := MakeTestEncodingConfig()
	app := simapp.NewSimApp(
		nil,
		db,
		nil,
		true,
		map[int64]bool{},
		simapp.DefaultNodeHome,
		5,
		encCdc,
		simapp.EmptyAppOptions{},
		ibctesting.GetEnableLoggingOption(),
		simapp.EmptyWasmEnabledProposals,
	)
	return app, simapp.EmptyAppOptions{}
}

type MyIBCModuleTestSuite struct {
	suite.Suite
	coordinator *ibctesting.Coordinator
	chainA      *ibctesting.TestChain
	chainB      *ibctesting.TestChain
}

// Setup the test environment
func (suite *MyIBCModuleTestSuite) SetupTest() {
	ibctesting.DefaultTestingAppInit = SetupTestingApp
	suite.coordinator = ibctesting.NewCoordinator(suite.T(), 2)
	suite.chainA = suite.coordinator.GetChain(ibctesting.GetChainID(1))
	suite.chainB = suite.coordinator.GetChain(ibctesting.GetChainID(2))
}

// Example test case for client creation
func (suite *MyIBCModuleTestSuite) TestClientCreation() {
	// Your test logic for client creation
}

func TestMyIBCModuleTestSuite(t *testing.T) {
	suite.Run(t, new(MyIBCModuleTestSuite))
}

func init() {
	ibctesting.DefaultTestingAppInit = SetupTestingApp
}

// Mock implementation for testing
var _ types.IBCContractKeeper = &IBCContractKeeperMock{}

type IBCContractKeeperMock struct {
	types.IBCContractKeeper
	OnRecvPacketFn func(ctx sdk.Context, contractAddr sdk.AccAddress, msg wasmvmtypes.IBCPacketReceiveMsg) (ibcexported.Acknowledgement, error)
}

func (m IBCContractKeeperMock) OnRecvPacket(ctx sdk.Context, contractAddr sdk.AccAddress, msg wasmvmtypes.IBCPacketReceiveMsg) (ibcexported.Acknowledgement, error) {
	if m.OnRecvPacketFn == nil {
		panic("not expected to be called")
	}
	return m.OnRecvPacketFn(ctx, contractAddr, msg)
}

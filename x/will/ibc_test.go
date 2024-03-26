package will_test

import (
	"encoding/json"
	"testing"

	// Import your app and other necessary packages
	// "github.com/cosmos/ibc-go/testing"
	// ibctesting "github.com/cosmos/ibc-go/v4/testing"
	// ibctesting "github.com/cosmos/ibc-go/testing"
	"cosmossdk.io/log"
	simapp "github.com/cosmos/cosmos-sdk/simapp"

	// simapp "cosmossdk.io/simapp"

	// ibctesting "github.com/cosmos/ibc-go/testing"
	// "github.com/cosmos/ibc-go/testing/simapp"
	ibctesting "github.com/cosmos/ibc-go/v8/testing"
	dbm "github.com/tendermint/tm-db"

	"github.com/stretchr/testify/suite"
)

// Define your app-specific testing setup if needed
// SetupTestingApp creates a new SimApp instance for testing
func SetupTestingApp() (ibctesting.TestingApp, map[string]json.RawMessage) {
	db := dbm.NewMemDB()
	encCdc := simapp.MakeTestEncodingConfig()
	app := simapp.NewSimApp(log.NewNopLogger(), db, nil, true, map[int64]bool{}, simapp.DefaultNodeHome, 5, encCdc, simapp.EmptyAppOptions{})
	return app, simapp.NewDefaultGenesisState(encCdc.Marshaler)
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

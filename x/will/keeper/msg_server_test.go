package keeper

import (
	"context"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/CosmWasm/wasmd/x/will/types"
)

// MockKeeper is a mock of Keeper interface
type MockKeeper struct {
	mock.Mock
}

// CreateWill mocks the CreateWill method in the IKeeper interface
func (m *MockKeeper) CreateWill(ctx context.Context, msg *types.MsgCreateWillRequest) (*types.Will, error) {
	args := m.Called(ctx, msg)
	return args.Get(0).(*types.Will), args.Error(1)
}

// GetWillByID mocks the GetWillByID method in the IKeeper interface
func (m *MockKeeper) GetWillByID(ctx context.Context, id string) (*types.Will, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*types.Will), args.Error(1)
}

// ListWillsByAddress mocks the ListWillsByAddress method in the IKeeper interface
func (m *MockKeeper) ListWillsByAddress(ctx context.Context, address string) ([]*types.Will, error) {
	args := m.Called(ctx, address)
	return args.Get(0).([]*types.Will), args.Error(1)
}

func (mk *MockKeeper) Claim(ctx context.Context, msg *types.MsgClaimRequest) error {
	args := mk.Called(ctx, msg)
	return args.Error(0)
}

func TestCreateWill(t *testing.T) {
	// Context for the tests
	ctx := context.TODO()

	// Create a mock keeper
	mockKeeper := new(MockKeeper)
	msgServer := NewMsgServerImpl(mockKeeper)

	// Define your test case
	testCases := []struct {
		name    string
		request *types.MsgCreateWillRequest
		setup   func()
		check   func(r *types.MsgCreateWillResponse, err error)
	}{
		{
			name: "successful will creation",
			request: &types.MsgCreateWillRequest{
				// Populate your request fields
				Creator:     "fakeid",
				Name:        "fakename",
				Beneficiary: "fakebeneficiary",
			},
			setup: func() {
				// Setup expectations and return values for your mock
				mockKeeper.On("CreateWill", mock.Anything, mock.Anything).Return(&types.Will{ID: "1"}, nil)
			},
			check: func(r *types.MsgCreateWillResponse, err error) {
				// Assert expectations
				require.NoError(t, err)
				require.NotNil(t, r)
				require.Equal(t, "1", r.Id)
			},
		},
		// Add more test cases for different scenarios, including error handling
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.setup() // Setup the mock expectations

			response, err := msgServer.CreateWill(ctx, tc.request)

			tc.check(response, err) // Check the outcome
			mockKeeper.AssertExpectations(t)
		})
	}
}

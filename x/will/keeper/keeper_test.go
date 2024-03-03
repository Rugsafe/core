package keeper_test

import (
	"context"
	"encoding/hex"
	"fmt"
	"testing"

	_proto "github.com/cosmos/gogoproto/proto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/proto"

	"cosmossdk.io/core/store"
	corestoretypes "cosmossdk.io/core/store"
	"cosmossdk.io/log"

	// corestoretypes "cosmossdk.io/core/store"
	storetypes "cosmossdk.io/store/types"

	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"

	// tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
	"github.com/CosmWasm/wasmd/app"
	"github.com/CosmWasm/wasmd/x/will/keeper"
	"github.com/CosmWasm/wasmd/x/will/types"
	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"
)

type MockKVStore struct {
	mock      mock.Mock
	mockStore corestoretypes.KVStore // Define the mock store
}

func (m *MockKVStore) Delete(key []byte) error {
	// Implement your custom logic for Delete using core.KVStore
	return nil
}

func (m *MockKVStore) Get(key []byte) ([]byte, error) {
	// Implement your logic for getting a value using your internal store
	return nil, nil // Replace with your actual implementation
}

func (m *MockKVStore) Has(key []byte) (bool, error) {
	// Mock the Has method. Return true or false based on your test needs
	return false, nil
}

func (m *MockKVStore) Iterator(start, end []byte) (storetypes.Iterator, error) {
	// Mock the Iterator method. Return a mock iterator based on your test needs.
	// This is a simplistic example. You might need a more complex setup depending on your test cases.
	return NewMockIterator()
}

func (m *MockKVStore) ReverseIterator(start, end []byte) (storetypes.Iterator, error) {
	// Mock the ReverseIterator method based on your test needs.
	// You can return a separate mock iterator for reverse iteration.
	return NewMockReverseIterator() // Example using a separate mock
}

func (m *MockKVStore) Set(key, value []byte) error {
	// Implement your mock logic for setting a value
	return nil
}

func NewMockIterator() (storetypes.Iterator, error) {
	// Implement a mock iterator that suits your test needs
	return (&MockIterator{}), nil // Example using a struct
}

func NewMockReverseIterator() (storetypes.Iterator, error) {
	// Example implementation. Adjust according to the actual interface requirements.
	return &MockIterator{}, nil
}

// Example MockIterator struct
type MockIterator struct {
	// Implement methods for Next, Key, Value, etc.
}

func (it *MockIterator) Close() error {
	// Implement the Close method if necessary. For a mock, it might do nothing.
	return nil
}

func (it *MockIterator) Domain() ([]byte, []byte) {
	// Return dummy or specific start and end keys as required by your tests
	return []byte("start"), []byte("end")
}

func (it *MockIterator) Error() error {
	// Return nil or an actual error depending on your test scenario
	return nil
}

func (it *MockIterator) Key() []byte {
	// Return a dummy key or a specific key based on your test needs
	return []byte("key")
}

func (it *MockIterator) Next() {
}

func (it *MockIterator) Valid() bool {
	return false
}

func (it *MockIterator) Value() []byte {
	// Return a dummy value or a specific value based on your test needs
	return []byte("value")
}

// Implement other required methods of types.KVStore

// func (m *MockKVStoreService) OpenKVStore(ctx context.Context) store.KVStore {
// 	return &MockKVStore{m.mockStore} // Return your mock KVStore
// }

// MockCodec is a mock for the Codec interface
type MockCodec struct {
	mock.Mock
}

func (m *MockCodec) InterfaceRegistry() codectypes.InterfaceRegistry {
	// Implement your mock logic for InterfaceRegistry
	return nil // Replace with your actual implementation
}

func (m *MockCodec) MarshalBinaryBare(msg proto.Message) ([]byte, error) {
	args := m.Called(msg)
	return args.Get(0).([]byte), args.Error(1)
}

func (m *MockCodec) UnmarshalBinaryBare(bz []byte, ptr proto.Message) error {
	args := m.Called(bz, ptr)
	return args.Error(0)
}

func (m *MockCodec) GetMsgAnySigners(msg *codectypes.Any) ([][]byte, proto.Message, error) {
	args := m.Called(msg)
	return args.Get(0).([][]byte), args.Get(1).(proto.Message), args.Error(2)
}

func (m *MockCodec) GetMsgV1Signers(msg _proto.Message) ([][]byte, proto.Message, error) {
	args := m.Called(msg)
	return args.Get(0).([][]byte), args.Get(1).(proto.Message), args.Error(2)
}

func (m *MockCodec) GetMsgV2Signers(msg proto.Message) ([][]byte, error) {
	args := m.Called(msg)
	// Assuming the mock setup correctly returns a slice of byte slices for the first argument
	return args.Get(0).([][]byte), args.Error(1)
}

func (m *MockCodec) Marshal(msg _proto.Message) ([]byte, error) {
	args := m.Called(msg)
	// Ensure the mock setup correctly returns a byte slice for the first argument
	return args.Get(0).([]byte), args.Error(1)
}

func (m *MockCodec) MarshalInterface(msg _proto.Message) ([]byte, error) {
	// Mock the MarshalInterface method according to your test needs.
	return nil, nil
}

func (m *MockCodec) MarshalInterfaceJSON(msg _proto.Message) ([]byte, error) {
	args := m.Called(msg)
	// Assuming the mock setup correctly returns a byte slice for the first argument
	return args.Get(0).([]byte), args.Error(1)
}

func (m *MockCodec) MarshalJSON(msg _proto.Message) ([]byte, error) {
	args := m.Called(msg)
	return args.Get(0).([]byte), args.Error(1)
}

func (m *MockCodec) MarshalLengthPrefixed(msg _proto.Message) ([]byte, error) {
	args := m.Called(msg)
	// Assuming the mock setup correctly returns a byte slice for the first argument
	return args.Get(0).([]byte), args.Error(1)
}

func (m *MockCodec) MustMarshal(msg _proto.Message) []byte {
	args := m.Called(msg)
	// Assuming the mock setup correctly returns a byte slice for the first argument
	return args.Get(0).([]byte)
}

func (m *MockCodec) MustMarshalJSON(msg _proto.Message) []byte {
	args := m.Called(msg)
	// Assuming the mock setup correctly returns a byte slice for the first argument
	return args.Get(0).([]byte)
}

func (m *MockCodec) MustMarshalLengthPrefixed(msg _proto.Message) []byte {
	args := m.Called(msg)
	// Assuming the mock setup correctly returns a byte slice for the first argument.
	// This method should panic if marshaling fails, mimicking the behavior of the actual method it mocks.
	return args.Get(0).([]byte)
}

func (m *MockCodec) MustUnmarshal(bz []byte, ptr _proto.Message) {
	// args := m.Called(bz, ptr)
	// Assuming you set up your mock to simulate the behavior of unmarshaling.
	// This method should panic if unmarshaling fails, mimicking the behavior of the actual method.
	// You might need to simulate setting the ptr value if your test depends on the unmarshaled result.
}

func (m *MockCodec) MustUnmarshalJSON(bz []byte, ptr _proto.Message) {
	// args := m.Called(bz, ptr)
	// Here you can simulate the behavior of unmarshaling JSON into the ptr.
	// Since MustUnmarshalJSON panics on failure, you should also simulate this behavior if necessary.
	// Typically, you would set the ptr to some value based on your mock setup.
}

func (m *MockCodec) MustUnmarshalLengthPrefixed(bz []byte, ptr _proto.Message) {
	// args := m.Called(bz, ptr)
	// Here you simulate the behavior of unmarshaling length-prefixed data into ptr.
	// Since MustUnmarshalLengthPrefixed panics on failure, simulate this behavior if necessary.
	// Typically, you would set ptr to some value based on your mock setup.
}

func (m *MockCodec) Unmarshal(bz []byte, ptr _proto.Message) error {
	args := m.Called(bz, ptr)
	// Optionally, you can simulate the behavior of unmarshaling the data into ptr.
	// You might return an error based on your mock setup to mimic failure scenarios.
	return args.Error(0)
}

func (m *MockCodec) UnmarshalInterface(bz []byte, ptr interface{}) error {
	args := m.Called(bz, ptr)
	// Optionally, simulate the behavior of unmarshaling the data into ptr.
	// You might return an error based on your mock setup to mimic failure scenarios.
	return args.Error(0)
}

func (m *MockCodec) UnmarshalInterfaceJSON(bz []byte, ptr interface{}) error {
	args := m.Called(bz, ptr)
	// Optionally simulate behavior of unmarshaling the JSON into ptr.
	// You might want to return an error based on your mock setup to mimic failure scenarios.
	return args.Error(0)
}

func (m *MockCodec) UnmarshalJSON(bz []byte, ptr _proto.Message) error {
	args := m.Called(bz, ptr)
	// Optionally simulate behavior of unmarshaling the JSON into ptr.
	// You might want to return an error based on your mock setup to mimic failure scenarios.
	return args.Error(0)
}

func (m *MockCodec) UnmarshalLengthPrefixed(bz []byte, ptr _proto.Message) error {
	args := m.Called(bz, ptr)
	// Optionally simulate behavior of unmarshaling the length-prefixed data into ptr.
	// You might want to return an error based on your mock setup to mimic failure scenarios.
	return args.Error(0)
}

func (m *MockCodec) UnpackAny(any *codectypes.Any, iface interface{}) error {
	args := m.Called(any, iface)
	// You can simulate behavior or return values here based on your test setup.
	// For example, to simulate unpacking the any into the provided iface:
	// if mockIface, ok := iface.(**YourExpectedType); ok {
	//     *mockIface = &YourExpectedTypeInstance
	//     return nil
	// }
	return args.Error(0)
}

func (m *MockCodec) mustEmbedCodec() {
	// This method is a no-op. It's only here to satisfy the Codec interface.
}

// MockKVStoreService is a mock for the KVStoreService interface
type MockKVStoreService struct {
	mock.Mock
}

func (m *MockKVStoreService) OpenKVStore(ctx context.Context) store.KVStore {
	// Ensure this method returns an object that implements store.KVStore
	return &MockKVStore{}
}

func setupMockKeeper(t *testing.T) (*keeper.Keeper, context.Context, *MockCodec, *MockKVStoreService) {
	w3llApp := app.Setup(t)
	mockedCodec := w3llApp.AppCodec()
	mockCodec := new(MockCodec)
	mockStoreService := new(MockKVStoreService)

	// ctx := context.Background()
	// Create a mock sdk.Context
	ctx := sdk.NewContext(nil, tmproto.Header{}, false, log.NewNopLogger())

	// Create a no-op logger. This logger does nothing but satisfies the interface requirement.
	noOpLogger := log.NewNopLogger()
	kpr := keeper.NewKeeper(mockedCodec, mockStoreService, noOpLogger)
	return &kpr, ctx, mockCodec, mockStoreService
}

const (
	// Example hexadecimal strings for public key and signature
	publicKeyHex = "3059...yourPublicKeyHexHere...AOGQ=="
	signatureHex = "0D02...yourSignatureHexHere...9C"
)

func TestKeeperCreateWill(t *testing.T) {
	kpr, ctx, mockCodec, mockStoreService := setupMockKeeper(t)

	// Set expectations on mocks
	// Example: mockCodec.On("MarshalBinaryBare", mock.Anything).Return([]byte{}, nil)
	// Example: mockStoreService.On("OpenKVStore", mock.Anything).Return(mockStore)

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

	// Verify that all expectations were met
	mockCodec.AssertExpectations(t)
	mockCodec.On("MustMarshal", mock.Anything).Return([]byte("expected output"))
	mockCodec.On("MustUnmarshal", mock.AnythingOfType("[]byte"), mock.AnythingOfType("*YourProtobufMessageType")).Run(func(args mock.Arguments) {
		// Optionally manipulate args[1] (*YourProtobufMessageType) to simulate unmarshaling if needed.
	}).Return()
	mockCodec.On("MustUnmarshalJSON", mock.AnythingOfType("[]byte"), mock.Anything).Run(func(args mock.Arguments) {
		// Optionally manipulate args[1] (interface{}) to simulate unmarshaling if needed.
	}).Return()
	mockCodec.On("MustUnmarshalLengthPrefixed", mock.AnythingOfType("[]byte"), mock.Anything).Run(func(args mock.Arguments) {
		// Optionally manipulate args[1] (interface{}) to simulate unmarshaling if needed.
	}).Return()
	mockCodec.On("Unmarshal", mock.AnythingOfType("[]byte"), mock.Anything).Return(nil)
	mockCodec.On("UnmarshalInterface", mock.AnythingOfType("[]byte"), mock.Anything).Return(nil)
	mockCodec.On("UnmarshalJSON", mock.AnythingOfType("[]byte"), mock.Anything).Return(nil)
	mockCodec.On("UnmarshalLengthPrefixed", mock.AnythingOfType("[]byte"), mock.Anything).Return(nil)
	mockStoreService.AssertExpectations(t)

	// Assertions to verify the contents of the will
	assert.Equal(t, msg.Creator, will.Creator, "will creator should match the request")
	assert.Equal(t, msg.Name, will.Name, "will name should match the request")
	assert.Equal(t, msg.Beneficiary, will.Beneficiary, "will beneficiary should match the request")
	assert.Equal(t, msg.Height, will.Height, "will height should match the request")

	// If you have specific expectations for the Components, verify those as well
	// This example assumes you want to check the length of the components slice
	assert.Len(t, will.Components, len(msg.Components), "number of will components should match the request")

}

/*
func TestKeeperCreateAndClaimWill(t *testing.T) {
	kpr, ctx, _, _ := setupMockKeeper(t)

	// Create Will
	msg := &types.MsgCreateWillRequest{
		Creator:     "creator-address",
		Name:        "Test Will",
		Beneficiary: "beneficiary-address",
		Height:      1,
		Components:  []*types.ExecutionComponent{},
	}

	will, err := kpr.CreateWill(sdk.UnwrapSDKContext(ctx), msg)
	require.NoError(t, err)
	assert.NotNil(t, will)

	// Prepare for claim
	publicKeyHex := "3059301306072a8648ce3d020106082a8648ce3d03010703420004..." // Example public key in hex
	signatureHex := "3045022100a34f..."                                         // Example signature in hex

	publicKeyBytes, err := hex.DecodeString(publicKeyHex)
	require.NoError(t, err)

	signatureBytes, err := hex.DecodeString(signatureHex)
	require.NoError(t, err)

	claimMsg := &types.MsgClaimRequest{
		WillId: will.ID,
		ClaimType: &types.MsgClaimRequest_SchnorrClaim{ // Assuming this is the correct type
			SchnorrClaim: &types.SchnorrClaim{
				PublicKey: publicKeyBytes,
				Signature: signatureBytes,
				Message:   []byte("Claim message"), // The message that was signed
			},
		},
	}

	// Process the claim
	// Note: Adjust the method call if your Keeper has a different method for processing claims
	err = kpr.Claim(sdk.UnwrapSDKContext(ctx), claimMsg)
	require.NoError(t, err)

	// Additional assertions to verify the claim was processed as expected
	// For example, verify the will's status, or that the beneficiary received the expected outcome
}
*/

func TestKeeperClaimWithSchnorrSignature(t *testing.T) {
	kpr, ctx, _, _ := setupMockKeeper(t)

	// Hardcoded values from your Schnorr signature example
	publicKeyHex := "2320a2da28561875cedbb0c25ae458e0a1d087834ae49b96a3f93cec79a8190c"
	signatureRHex := "7ab0edb9b0929b5bb4b47dfb927d071ecc5de75985662032bb52ef3c5ace640b"
	signatureSHex := "165c6df5ea8911a6c0195a3140be5119a5b882e91b34cbcdd31ef3f5b0035b06"
	message := "message-2b-signed"

	// Convert hexadecimal strings to bytes
	publicKeyBytes, err := hex.DecodeString(publicKeyHex)
	require.NoError(t, err)
	signatureRBytes, err := hex.DecodeString(signatureRHex)
	require.NoError(t, err)
	signatureSBytes, err := hex.DecodeString(signatureSHex)
	require.NoError(t, err)
	// messageBytes := []byte(message)

	// Assuming the signature is the concatenation of R and S components
	signatureBytes := append(signatureRBytes, signatureSBytes...)

	msg := &types.MsgCreateWillRequest{
		Creator:     "creator-address",
		Name:        "Test Will",
		Beneficiary: "beneficiary-address",
		Height:      1,
		Components: []*types.ExecutionComponent{
			{
				Name:   "SchnorrSignatureComponent",
				Id:     "unique-component-id",
				Status: "active",
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

	// Construct the claim request with the Schnorr claim
	claimMsg := &types.MsgClaimRequest{
		WillId:      willID,
		Claimer:     "jovi",
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
	err = kpr.Claim(sdk.UnwrapSDKContext(ctx), claimMsg)
	require.NoError(t, err, "processing Schnorr claim should not produce an error")

	// Additional assertions can be made here to verify the claim was processed as expected
}

package keeper

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strconv"
	"time"

	"github.com/bwesterb/go-ristretto"
	capabilitykeeper "github.com/cosmos/ibc-go/modules/capability/keeper"
	icatypes "github.com/cosmos/ibc-go/v8/modules/apps/27-interchain-accounts/types"
	clienttypes "github.com/cosmos/ibc-go/v8/modules/core/02-client/types" //nolint:staticcheck
	channeltypes "github.com/cosmos/ibc-go/v8/modules/core/04-channel/types"
	host "github.com/cosmos/ibc-go/v8/modules/core/24-host"
	"go.dedis.ch/kyber/v3/group/edwards25519"

	"cosmossdk.io/collections"
	corestoretypes "cosmossdk.io/core/store"
	"cosmossdk.io/errors"
	"cosmossdk.io/log"
	storetypes "cosmossdk.io/store/types"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/telemetry"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	authkeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"

	wasmkeeper "github.com/CosmWasm/wasmd/x/wasm/keeper"
	"github.com/CosmWasm/wasmd/x/will/schemes/schnorr"
	"github.com/CosmWasm/wasmd/x/will/types"
)

type IKeeper interface {
	CreateWill(
		ctx context.Context,
		msg *types.MsgCreateWillRequest,
	) (*types.Will, error)
	GetWillByID(ctx context.Context, id string) (*types.Will, error)
	ListWillsByAddress(ctx context.Context, address string) ([]*types.Will, error)
	Claim(ctx context.Context, msg *types.MsgClaimRequest) error
}

type IContractCall interface {
	CallContract(address sdk.AccAddress)
}

// var (
// 	_ ibctypes.ChannelKeeper = (*Keeper)(nil)
// 	_ ibctypes.PortKeeper    = (*Keeper)(nil)
// )

type (
	// ContractHandler struct{}

	Keeper struct {
		storeService corestoretypes.KVStoreService
		// storeService storetypes.KVStoreKey
		cdc                    codec.Codec
		storeKey               storetypes.StoreKey // Add this line
		channelKeeper          ChannelKeeper
		scopedKeeper           capabilitykeeper.ScopedKeeper
		portKeeper             PortKeeper
		wasmKeeper             wasmkeeper.Keeper
		bankKeeper             bankkeeper.Keeper
		permissionedWasmKeeper wasmkeeper.PermissionedKeeper
		// capabilityKeeper CapabilityKeeper
		capabilityKeeper capabilitykeeper.Keeper
		accountKeeper    authkeeper.AccountKeeper

		params    collections.Item[types.Params]
		authority string
	}

	// ScopedKeeper struct {
	// 	cdc      codec.BinaryCodec
	// 	storeKey storetypes.StoreKey
	// 	memKey   storetypes.StoreKey
	// 	capMap   map[uint64]*captypes.Capability
	// 	module   string
	// }
)

func NewKeeper(
	cdc codec.Codec,
	storeService corestoretypes.KVStoreService,
	// storeService storetypes.KVStoreKey,
	logger log.Logger,
	channelKeeper ChannelKeeper,
	portKeeper icatypes.PortKeeper,
	scopedKeeper capabilitykeeper.ScopedKeeper,
	scopedIBCKeeper capabilitykeeper.ScopedKeeper,

	// capabilityKeeper CapabilityKeeper,
	capabilityKeeper capabilitykeeper.Keeper,
	wk wasmkeeper.Keeper,
	bk bankkeeper.Keeper,
	pwk wasmkeeper.PermissionedKeeper,
	ak authkeeper.AccountKeeper,
) Keeper {
	// fmt.Println("NewKeeper:")
	// sb := collections.NewSchemaBuilder(storeService)

	// sk := ScopedKeeper {
	// 		cdc      codec.BinaryCodec
	// 		storeKey storetypes.StoreKey
	// 		memKey   storetypes.StoreKey
	// 		capMap   map[uint64]*captypes.Capability
	// 		module   string
	// }

	// fmt.Println("scopedKeeper:")
	// fmt.Println(scopedKeeper)
	// fmt.Println(scopedIBCKeeper)
	// fmt.Println(capabilityKeeper)
	keeper := &Keeper{
		storeService:           storeService,
		cdc:                    cdc,
		channelKeeper:          channelKeeper,
		portKeeper:             portKeeper,
		scopedKeeper:           scopedKeeper,
		wasmKeeper:             wk,
		bankKeeper:             bk,
		permissionedWasmKeeper: pwk,
		capabilityKeeper:       capabilityKeeper,
		accountKeeper:          ak,
	}

	return *keeper
}

// GetParams returns the total set of wasm parameters.
func (k Keeper) GetParams(ctx context.Context) types.Params {
	p, err := k.params.Get(ctx)
	if err != nil {
		panic(err)
	}
	return p
}

// SetParams sets all will parameters.
// func (k Keeper) SetParams(ctx context.Context, ps types.Params) error {
// 	return k.params.Set(ctx, ps)
// }

// SetParams sets the transfer module parameters.
func (k Keeper) SetParams(ctx sdk.Context, params types.Params) {
	fmt.Println("SetParams k.storeKey")
	fmt.Println(k.storeKey)
	// store := ctx.KVStore(k.storeKey)
	store := k.storeService.OpenKVStore(ctx)
	bz := k.cdc.MustMarshal(&params)
	store.Set([]byte(types.ParamsKey), bz)
}

// GetAuthority returns the x/will module's authority.
func (k Keeper) GetAuthority() string {
	return k.authority
}

func (k *Keeper) GetBankKeeper() bankkeeper.Keeper {
	return k.bankKeeper
}

func (k *Keeper) GetAccountKeeper() *authkeeper.AccountKeeper {
	return &k.accountKeeper
}

func (k *Keeper) GetChannelKeeper() ChannelKeeper {
	return k.channelKeeper
}

/*
@name
@desc
@param
*/
// TruncateHash creates a shorter hash by taking the first n bytes of the SHA256 hash.
func TruncateHash(input []byte, n int) ([]byte, error) {
	if n <= 0 {
		return nil, fmt.Errorf("n must be positive")
	}
	hash := sha256.Sum256(input)
	if n > len(hash) {
		return nil, fmt.Errorf("n is greater than the hash size")
	}
	return hash[:n], nil
}

/*
@name
@desc
@param
*/
func (k Keeper) GetWillByID(ctx context.Context, id string) (*types.Will, error) {
	sdkCtx := sdk.UnwrapSDKContext(ctx) // Make sure you have a way to convert or access sdk.Context
	store := k.storeService.OpenKVStore(sdkCtx)
	var will types.Will
	fmt.Println("GetWillByID: " + id)
	willBz, willErr := store.Get(types.GetWillKey(id))
	fmt.Println("========= get will by id ============")
	fmt.Println(willBz)
	if willErr != nil {
		return nil, fmt.Errorf("will with ID %s not found", id)
	}
	k.cdc.MustUnmarshal(willBz, &will)
	return &will, nil
}

// TODO: use Decentralized Identifier in Will ID
// @note: https://pkg.go.dev/go.bryk.io/pkg/did
// @note: https://w3c.github.io/did-core/
//
//	func createWillId(creator string, name string, beneficiary string, height int64) string {
//		willID := fmt.Sprintf("%s-%s-%s-%s", creator, name, beneficiary, strconv.Itoa(int(height)))
//		fmt.Println("New Will ID: ", willID)
//		return willID
//	}
func createWillId(creator string, name string, beneficiary string, height int64) string {
	baseString := fmt.Sprintf("%s|%s|%s|%d", creator, name, beneficiary, height)
	hash := sha256.Sum256([]byte(baseString))
	willID := fmt.Sprintf("did:will:%x", hash[:])
	fmt.Println("New Will ID: ", willID)
	return willID
}

/*
@name CreateWill
@desc
@param
*/
func (k *Keeper) CreateWill(ctx context.Context, msg *types.MsgCreateWillRequest) (*types.Will, error) {
	store := k.storeService.OpenKVStore(ctx)
	if sdk.UnwrapSDKContext(ctx).BlockHeight() > msg.Height {
		// var errheight error
		fmt.Println("Target Block height is less than the current block height")
		// return nil, errors.Wrap(errheight, "inside k.createWill, block height is greater than submitted will execution height")
		return nil, errors.Wrapf(sdkerrors.ErrInvalidRequest, "block height %d is greater than submitted will execution height %d", sdk.UnwrapSDKContext(ctx).BlockHeight(), msg.Height)

	}
	// Concatenate values to generate a unique hash
	concatValues := createWillId(msg.Creator, msg.Name, msg.Beneficiary, msg.Height)
	// idBytes := []byte(concatValues)

	// Generate a truncated hash of the concatenated values
	// truncatedHash, err := TruncateHash(idBytes, 16) // Truncate SHA256 to 16 bytes
	// if err != nil {
	// 	return nil, err
	// }

	// Convert the truncated hash bytes to a hex string for safe serialization
	// idString := hex.EncodeToString(idBytes)
	// fmt.Println(fmt.Printf("NEWLY CREATED WILL: %s", idString))

	// TODO: verify components, as this is already done in client/cli/tx.go
	// verifyComponents(msg.components)

	// TODO: ENSURE EACH COMPONENT HAS AN OUTPUT TYPE AND AN ACCESS TYPE
	// Construct the will object
	will := types.Will{
		// ID:          fmt.Sprintf("did:will:%x", idString),
		ID:          concatValues,
		Creator:     msg.Creator,
		Name:        msg.Name,
		Beneficiary: msg.Beneficiary,
		Height:      msg.Height,
		Status:      "live",
		Components:  msg.Components,
	}
	// Marshal the will object to bytes
	willBz := k.cdc.MustMarshal(&will)
	fmt.Println("inside k.createWill: " + concatValues)
	if willBz == nil {
		var errBz error
		return nil, errors.Wrap(errBz, "inside k.createWill, willBz is nil") // Make sure to handle the error appropriately
	}

	// Use the GetWillKey function to get a unique byte key for this will
	key := types.GetWillKey(concatValues)
	// key := types.GetWillKey("zmxjiudojne844jdsbndsbdyuikdbaazxqetrsdshudjsdhuekdsxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxnxnmcnmcndhdiohsiodsdhsdoshdsdjksdhjksdsdsdhjsdjsdhjksdjshjdhjshdjksjdhsjdhks")
	fmt.Println("KEY")
	fmt.Println(key)

	// Store the marshaled will in the module's store
	// storeErr := store.Set(key, willBz)
	storeErr := store.Set([]byte(key), willBz)
	if storeErr != nil {
		return nil, errors.Wrap(storeErr, "inside k.createWill storeErr, KV store set threw an error")
	}
	///////////// store at height
	// Assuming you want to store the will's ID under a key derived from its height for some indexing purpose
	// Handling storage for heightKey with WillIds message
	heightKey := types.GetWillKey(strconv.Itoa(int(will.Height)))
	var willIdsAtHeight types.WillIds
	existingWillsBz, _ := store.Get([]byte(heightKey)) // Simplified error handling
	if existingWillsBz != nil {
		k.cdc.MustUnmarshal(existingWillsBz, &willIdsAtHeight)
	}

	// TODO: make this a chain param to be changed via governance
	if len(willIdsAtHeight.Ids) > 10 {
		var errBlockHeightLength error
		return nil, errors.Wrapf(errBlockHeightLength, "error: cannot add will during create will, too many wills at block height: %s", willIdsAtHeight.Ids)
	}

	willIdsAtHeight.Ids = append(willIdsAtHeight.Ids, will.ID)
	updatedHeightBz := k.cdc.MustMarshal(&willIdsAtHeight)
	store.Set([]byte(heightKey), updatedHeightBz)

	// Handling storage for creator key, ensuring unique insertion
	creatorKey := types.GetWillKey(msg.Creator)
	var willIdsAtCreator types.WillIds
	existingWillsForCreatorBz, _ := store.Get([]byte(creatorKey)) // Simplified error handling
	if existingWillsForCreatorBz != nil {
		k.cdc.MustUnmarshal(existingWillsForCreatorBz, &willIdsAtCreator)
	}
	if !contains(willIdsAtCreator.Ids, will.ID) {
		willIdsAtCreator.Ids = append(willIdsAtCreator.Ids, will.ID)
	}
	updatedCreatorBz := k.cdc.MustMarshal(&willIdsAtCreator)
	store.Set([]byte(creatorKey), updatedCreatorBz)

	fmt.Println("KEEPER TEST DEBUG:")
	fmt.Println(will.ID)
	fmt.Println(willIdsAtHeight.Ids)

	return &will, nil
}

// contains checks if a string is present in a slice of strings.
func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

/*
@name
@desc
@param
*/

func (k Keeper) ListWillsByAddress(ctx context.Context, address string) ([]*types.Will, error) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	store := k.storeService.OpenKVStore(sdkCtx)

	// Use address to construct the key for fetching associated will IDs
	addressKey := types.GetWillKey(address)
	willIDsBz, err := store.Get([]byte(addressKey))
	if err != nil {
		return nil, errors.Wrap(err, "failed to fetch will IDs for address")
	}
	if willIDsBz == nil {
		// No wills associated with this address
		return []*types.Will{}, nil
	}

	// Deserialize the will IDs
	var willIds types.WillIds
	err = k.cdc.Unmarshal(willIDsBz, &willIds)
	if err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal will IDs")
	}

	// Fetch and collect the wills by their IDs
	var wills []*types.Will
	for _, willID := range willIds.Ids {
		will, err := k.GetWillByID(ctx, willID)
		if err != nil {
			// Log the error and continue to the next ID if a specific will cannot be fetched
			fmt.Printf("Error fetching will by ID %s: %v\n", willID, err)
			continue
		}
		wills = append(wills, will)
	}

	return wills, nil
}

/*
@name
@desc
@param
*/
func (k Keeper) updateWillStatusAndStore(ctx context.Context, will *types.Will, componentIndex int) error {
	fmt.Println("Updating will status and storing it")
	store := k.storeService.OpenKVStore(ctx)
	// concatValues := createWillId(will.Creator, will.Name, will.Beneficiary, will.Height)
	// willID := hex.EncodeToString([]byte(concatValues))
	willID := createWillId(will.Creator, will.Name, will.Beneficiary, will.Height)
	key := types.GetWillKey(willID)
	fmt.Println(fmt.Sprintf("Storing will with ID: %s", willID))

	willBz := k.cdc.MustMarshal(will)
	storeErr := store.Set(key, willBz)
	if storeErr != nil {
		return errors.Wrapf(storeErr, "error: could not save will ID with updated component status")
	}

	return nil
}

/*
@name Claim
@description the function to make a claim on a will component
@param ctx Context to pass context from the sdk
@param msg MsgClaimRequest the message structure holding params for the claim request
*/
func (k Keeper) Claim(ctx context.Context, msg *types.MsgClaimRequest) error {
	// Retrieve the will by ID to ensure it exists and to process the claim against it
	fmt.Println("CLAIM FROM KEEPER: ", msg.Claimer)
	fmt.Println(msg)
	will, err := k.GetWillByID(ctx, msg.WillId)
	fmt.Println("THE WILL TO CLAIM")
	fmt.Println(will)
	if err != nil {
		return err
	}

	// Assuming GetWillByID returns nil for non-existent wills
	if will == nil {
		// Handle the case where the will doesn't exist
		fmt.Println("CANNOT CLAIM WILL, WILL DOES NOT EXIST")
		return fmt.Errorf("will with ID %s not found", msg.WillId)
	}

	// If there are specific fields that should be checked to determine if the will is "blank"
	if will.ID == "" || will.Creator == "" {
		// TODO: THIS SHOULD NEVER FIRE BECAUSE WILL.ID AND WILL.CREATOR SHOULD BE CHECKED UPON CREATION
		fmt.Println("CANNOT CLAIM WILL, WILL id and creator are nil")
		// Assuming ID and Creator being empty means the will is "blank"
		return fmt.Errorf("will with ID %s is blank", msg.WillId)
	}

	// will must be expired
	if will.Status != "expired" {
		fmt.Println("CANNOT CLAIM WILL, AS IT IS NOT EXPIRED")
		return fmt.Errorf("will with ID %s is NOT EXPIRED", msg.WillId)
	}

	fmt.Printf("CLAIMING WILL %s", msg.WillId)
	fmt.Println(msg)

	// TODO: the will component is in the will object here
	// but we have the claim ID, so iterate over the will to fetch the actual component itself
	// that matches the claim id
	var componentIndex int = -1
	for i, component := range will.Components {
		if component.Id == msg.ComponentId {
			componentIndex = i
			break
		}
	}

	if componentIndex == -1 {
		fmt.Printf("component with ID %s not found in will ID %s\n", msg.ComponentId, msg.WillId)
		return fmt.Errorf("component with ID %s not found in will ID %s", msg.ComponentId, msg.WillId)
	}

	// At this point, you have the index of the component being claimed.
	// var component *types.ExecutionComponent = will.Components[componentIndex]
	// You can now check its status before proceeding with the claim.
	if will.Components[componentIndex].Status != "active" {
		fmt.Printf("component with ID %s is not active and cannot be claimed\n", msg.ComponentId)
		return fmt.Errorf("component with ID %s is not active and cannot be claimed", msg.ComponentId)
	}

	// check access handler to ensure users can only access it with the proper access
	if accessErr := k.AccessHandler(ctx, will.Components[componentIndex], *will, msg); accessErr != nil {
		fmt.Printf("component with ID %s, Access Errored: %s\n", msg.ComponentId, accessErr)
		return fmt.Errorf("component with ID %s, Access Errored: %s", msg.ComponentId, accessErr)
	}

	var claimErr error

	// Process the claim based on its type
	switch claim := msg.ClaimType.(type) {
	case *types.MsgClaimRequest_SchnorrClaim:

		// todo: invoke schno

		// Assuming the public key and signature are provided as byte slices in the claim
		fmt.Println(claim)
		// TODO: pass in the component, not the component id lol
		claimErr = k.processSchnorrClaim(ctx, claim, will, componentIndex)
		if claimErr != nil {
			return fmt.Errorf("component with ID %s processSchnorrClaim FAILED and cannot be claimed: %s", msg.ComponentId, claimErr)
		}

		fmt.Println("Schnorr signature verified and saved now successfully.")
		k.OutputHandler(sdk.UnwrapSDKContext(ctx), will.Components[componentIndex], *will)

	case *types.MsgClaimRequest_PedersenClaim:

		// Process PedersenClaim
		fmt.Printf("Processing Pedersen claim with commitment: %x, blinding factor: %x, and value: %x\n", claim.PedersenClaim.Commitment, claim.PedersenClaim.BlindingFactor, claim.PedersenClaim.Value)
		// TODO
		fmt.Println(claim)

		claimErr = k.processPedersenClaim(ctx, will, componentIndex, claim)
		if claimErr != nil {
			return fmt.Errorf("component with ID %s processPedersenClaim FAILED and cannot be claimed: %s", msg.ComponentId, claimErr)
		}

		fmt.Println("Pedersen commitments verified and saved now successfully.")
		k.OutputHandler(sdk.UnwrapSDKContext(ctx), will.Components[componentIndex], *will)

	case *types.MsgClaimRequest_GnarkClaim:
		// Process GnarkClaim
		fmt.Printf("Processing Gnark claim with proof: %x and public inputs: %x\n", claim.GnarkClaim.Proof, claim.GnarkClaim.PublicInputs)
		// TODO

		k.OutputHandler(sdk.UnwrapSDKContext(ctx), will.Components[componentIndex], *will)

	default:
		// Handle unknown claim type
		fmt.Println("unknown claim type provided")
		return fmt.Errorf("unknown claim type provided")
	}

	if claimErr != nil {
		return claimErr // Properly propagate the error
	}

	// Assuming the claim has been validated successfully, you can then update the will's status or components accordingly
	return nil
}

/*
@name
@desc
@param
*/
func (k Keeper) processSchnorrClaim(ctx context.Context, claim *types.MsgClaimRequest_SchnorrClaim, will *types.Will, componentIndex int) error {
	// publicKeyBytes := claim.SchnorrClaim.PublicKey // The public key bytes
	// NOTE: use the public key from the will
	// publicKeyBytes, _ := hex.DecodeString(string(claim.SchnorrClaim.PublicKey))
	// publicKeyBytes, _ := hex.DecodeString(string(claim.SchnorrClaim.PublicKey))
	publicKeyBytes, _ := hex.DecodeString(string(claim.SchnorrClaim.PublicKey))

	fmt.Printf("string claim.SchnorrClaim.PublicKey %s: \n", string(claim.SchnorrClaim.PublicKey))
	fmt.Printf("claim.SchnorrClaim.PublicKey %s: \n", claim.SchnorrClaim.PublicKey)
	fmt.Printf("publicKeyBytes %s: \n", publicKeyBytes)

	// signatureBytes := claim.SchnorrClaim.Signature // The signature bytes
	signatureBytes, _ := hex.DecodeString(string(claim.SchnorrClaim.Signature))
	fmt.Printf("string claim.SchnorrClaim.Signature %s: \n", string(claim.SchnorrClaim.Signature))
	fmt.Printf("claim.SchnorrClaim.Signature %s: \n", claim.SchnorrClaim.Signature)
	fmt.Printf("Signature %s: \n", signatureBytes)

	message := claim.SchnorrClaim.Message // The message as a byte slice
	// message, _ := hex.DecodeString(string(claim.SchnorrClaim.Message))

	curve := edwards25519.NewBlakeSHA256Ed25519()
	// Convert the public key bytes to a kyber.Point
	publicKeyPoint := curve.Point()
	if err := publicKeyPoint.UnmarshalBinary(publicKeyBytes); err != nil {
		fmt.Printf("failed to unmarshal public key: %v\n", err)
		return fmt.Errorf("failed to unmarshal public key: %v", err)
	}

	// Assuming the signature consists of R and S components concatenated
	// and that each component is of equal length
	sigLen := len(signatureBytes) / 2
	rBytes := signatureBytes[:sigLen]
	sBytes := signatureBytes[sigLen:]

	// Convert R and S bytes to kyber.Point and kyber.Scalar respectively
	r := curve.Point()
	if err := r.UnmarshalBinary(rBytes); err != nil {
		fmt.Printf("failed to unmarshal R component: %v", err)
		return fmt.Errorf("failed to unmarshal R component: %v", err)
	}
	s := curve.Scalar().SetBytes(sBytes)

	// Hash the message to a scalar using your Schnorr Hash function
	messageScalar := schnorr.Hash(string(message) + string(hex.EncodeToString(publicKeyBytes))) // Convert the message to a string if your Hash function expects a string

	// Construct the Signature struct
	schnorrSignature := schnorr.Signature{R: r, S: s}

	// Verify the Schnorr signature
	if !schnorr.Verify(messageScalar, schnorrSignature, publicKeyPoint) {
		return fmt.Errorf("schnorr signature verification failed")
		// panic(99)
	}
	// panic(99)

	// TODO: IF MESSAGE IS ENCRYPTED:?
	// verify the encrypted message matches one stored in will

	fmt.Println("Schnorr signature verified successfully.")
	will.Components[componentIndex].Status = "claimed"
	return k.updateWillStatusAndStore(ctx, will, componentIndex)
}

/*
@name
@desc
@param
*/
func (k Keeper) processPedersenClaim(ctx context.Context, will *types.Will, componentIndex int, claim *types.MsgClaimRequest_PedersenClaim) error {
	fmt.Println("Starting processPedersenClaim")

	// Extract the Pedersen commitment from the component
	storedCommitment := will.Components[componentIndex].GetClaim().GetPedersen()

	if storedCommitment == nil {
		return fmt.Errorf("Error: Pedersen commitment not found in the component")
	}

	fmt.Println("1: ", storedCommitment.Commitment)
	// Deserialize the stored commitment and target commitment
	storedCommitmentPoint, err := k.DeserializeCommitment(storedCommitment.Commitment)
	if err != nil {
		return fmt.Errorf("failed to deserialize stored commitment: %v", err)
	}

	fmt.Println("2")
	claimCommitmentPoint, err := k.DeserializeCommitment(claim.PedersenClaim.Commitment)
	if err != nil {
		return fmt.Errorf("failed to deserialize claimed commitment: %v", err)
	}

	fmt.Println("3")
	targetCommitmentPoint, err := k.DeserializeCommitment(storedCommitment.TargetCommitment)
	if err != nil {
		return fmt.Errorf("failed to deserialize target commitment: %v", err)
	}

	// Add commitments
	resultCommitment := k.AddCommitments(storedCommitmentPoint, claimCommitmentPoint)
	fmt.Println(storedCommitmentPoint)
	fmt.Println(claimCommitmentPoint)
	fmt.Println(targetCommitmentPoint)
	fmt.Println(resultCommitment)
	// Check if the result matches the target
	if !resultCommitment.Equals(&targetCommitmentPoint) {
		return fmt.Errorf("commitment verification failed")
	}

	fmt.Println("Commitment verified successfully.")
	will.Components[componentIndex].Status = "claimed"
	return k.updateWillStatusAndStore(ctx, will, componentIndex)
}

/*
@name
@desc
@param
*/
// func EncryptWithPublicKey(cosmosPub secp256k1.PubKey, message []byte) ([]byte, error) {
// 	// Convert Cosmos SDK PubKey to ECDSA Public Key
// 	pub := &ecdsa.PublicKey{
// 		Curve: elliptic.P256(), // make sure to use the correct curve
// 		X:     new(big.Int).SetBytes(cosmosPub.XBytes()),
// 		Y:     new(big.Int).SetBytes(cosmosPub.YBytes()),
// 	}

// 	// Import public key to ECIES
// 	pubECIES := ecies.ImportECDSAPublic(pub)

// 	// Encrypt using ECIES
// 	return ecies.Encrypt(rand.Reader, pubECIES, message, nil, nil)
// }

/*
@name
@desc
@param
*/
func (k Keeper) AddCommitments(a, b ristretto.Point) ristretto.Point {
	var result ristretto.Point
	result.Add(&a, &b) // Add points
	return result
}

// Deserialize a commitment from bytes to a ristretto.Point
func (k Keeper) DeserializeCommitment(data []byte) (ristretto.Point, error) {
	var point ristretto.Point
	err := point.UnmarshalBinary(data)
	if err != nil {
		return ristretto.Point{}, err // return an empty point on error
	}
	return point, nil
}

// ///////////////////////////////////// expirations
// expiration
func (k Keeper) SetWillExpiryIndex(ctx sdk.Context, expiryHeight int64, willID string) {
	store := ctx.KVStore(k.storeKey)
	expiryKey := []byte(fmt.Sprintf("expiry:%d:%s", expiryHeight, willID))
	store.Set(expiryKey, []byte(willID))
}

/*
@name
@desc
@param
*/
func (k Keeper) GetWillsByExpiry(ctx sdk.Context, expiryHeight int64) ([]*types.Will, error) {
	store := ctx.KVStore(k.storeKey)
	iterator := storetypes.KVStorePrefixIterator(store, []byte(fmt.Sprintf("expiry:%d:", expiryHeight)))
	defer iterator.Close()

	var wills []*types.Will
	for ; iterator.Valid(); iterator.Next() {
		willID := string(iterator.Value())
		will, err := k.GetWillByID(ctx, willID)
		if err != nil {
			return nil, err
		}
		wills = append(wills, will)
	}
	return wills, nil
}

/*
@name
@desc
@param
*/
func (k Keeper) BeginBlocker(ctx sdk.Context) error {
	defer telemetry.ModuleMeasureSince(types.ModuleName, time.Now(), telemetry.MetricKeyBeginBlocker)

	// Get the current block height
	blockHeight := ctx.BlockHeight()
	fmt.Printf("Processing wills at block height: %d\n", blockHeight)

	// Access the store
	store := k.storeService.OpenKVStore(ctx)

	// Construct the height key to fetch will IDs associated with the current block height
	heightKey := types.GetWillKey(strconv.Itoa(int(blockHeight)))
	willIDsBz, err := store.Get([]byte(heightKey))

	// If there is an error fetching the will IDs or if there are no wills for this block height, return early
	if err != nil || willIDsBz == nil {
		fmt.Println("No wills to process for this block height or unable to fetch will IDs.")
		return nil
	}

	// Deserialize the list of will IDs
	var willIds types.WillIds
	err = k.cdc.Unmarshal(willIDsBz, &willIds)
	if err != nil {
		fmt.Printf("Error unmarshaling will IDs: %v\n", err)
		return nil
	}

	// Iterate over each will ID
	for _, willID := range willIds.Ids {
		// Fetch the will object using its ID
		will, err := k.GetWillByID(ctx, willID)
		if err != nil {
			fmt.Printf("Error fetching will with ID %s: %v\n", willID, err)
			continue // Proceed to the next will if there's an issue fetching this one
		}

		// Perform the desired operations on the will object here
		// This is where you would implement the logic specific to your application's requirements
		fmt.Printf("Successfully fetched will with ID %s for further processing.\n", will.ID)

		// if the will is not live, because this transition should only happen is the will is going from inactive->active (maybe terminology can be better)
		if will.Status != "live" {
			fmt.Printf("Error executing will components with WILL ID %s, will is NOT EXPIRED: %v\n", willID, err)
			continue
		}

		for component_index, component := range will.Components {
			fmt.Printf("Iterating over compnents for will ID %s for further processing.\n", will.ID)
			fmt.Println(component_index)
			fmt.Println(component)
			switch c := component.ComponentType.(type) {

			case *types.ExecutionComponent_Transfer:
				fmt.Printf("Transfer component found, to: %s, amount: %s\n", c.Transfer.To, c.Transfer.Amount.String())

				// TODO: actually execute the token send
				transferErr := k.ExecuteTransfer(ctx, component, *will)
				if transferErr != nil {
					continue
				}

				// update status to executed
				component.Status = "executed"

				// TODO: should we do outputs on execution components, or only claims?
				// HandleOutput()
				// k.OutputHandler(ctx, component, *will)

			case *types.ExecutionComponent_Claim:
				fmt.Printf("Claim component found, evidence")
				// set all claimable components to active - can now have claims submitted
				component.Status = "active"
				// fmt.Printf("Claim component found, evidence: %s\n", c.Claim.Evidence)

				// TODO:
				// DONT EXECUTE HandleOutput()
			case *types.ExecutionComponent_Contract:

				_, err := k.ExecuteContract(ctx, c, will.Creator)
				if err != nil {
					// Handle error, maybe log it or take appropriate action.
					continue
				}

				// Update the status based on the execution result.
				component.Status = "executed"
				// Handle other component outputs if necessary.
				// k.OutputHandler(ctx, component, *will)

			case *types.ExecutionComponent_IbcMsg:
				// send an IBC message
				// TODO: DEV TESTING FOR SENDIBCMESSAGE
				// HandleOutput()

				k.SendIBCMessage(sdk.UnwrapSDKContext(ctx), component, *will)
				// change status depending on result
				component.Status = "executed"

				// k.OutputHandler(ctx, component, *will)

			case *types.ExecutionComponent_IbcSend:
				// change status depending on result
				component.Status = "executed"

			default:
				fmt.Println("Unknown component type found")
			}
		}

		fmt.Printf("Will ID: %s, Name: %s, Beneficiary: %s, Height: %d\n", will.ID, will.Name, will.Beneficiary, will.Height)

		// update will
		will.Status = "expired"
		// concatValues := createWillId(will.Creator, will.Name, will.Beneficiary, will.Height)
		// willID := hex.EncodeToString([]byte(concatValues))
		// willID := hex.EncodeToString(idString)
		willID := createWillId(will.Creator, will.Name, will.Beneficiary, will.Height)
		key := types.GetWillKey(willID)
		fmt.Println(fmt.Printf("BEGIN BLOCKER WILL EXECUTED: %s", willID))

		willBz := k.cdc.MustMarshal(will)
		storeErr := store.Set(key, willBz)

		if storeErr != nil {
			return errors.Wrapf(storeErr, "inside k.beginBlocker storeErr, KV store set threw an error after updating will: %s", will.ID)
		}

	}

	// DEBUG
	// os.Exit(10)

	return nil
}

// func (k *Keeper) EndBlocker(ctx context.Context) error {
func (k *Keeper) EndBlocker(ctx sdk.Context) error {
	defer telemetry.ModuleMeasureSince(types.ModuleName, time.Now(), telemetry.MetricKeyEndBlocker)
	fmt.Println("INSIDE END BLOCKER FOR WILL MODULE")
	return nil
}

/*
@name OutputHandler
@desc OutputHandler processes the output based on the component's output type and executes corresponding actions.
@param
*/
func (k Keeper) OutputHandler(ctx sdk.Context, component *types.ExecutionComponent, will types.Will) error {
	// Assuming OutputType is correctly configured to be used as a type switch
	fmt.Println("Output Handler:")
	fmt.Println(component)
	// todo if output is nil
	switch output := component.OutputType.OutputType.(type) {
	case *types.ComponentOutput_OutputTransfer:
		toAddr, err := sdk.AccAddressFromBech32(output.OutputTransfer.Address)
		if err != nil {
			return err
		}
		coins := sdk.NewCoins(*output.OutputTransfer.Amount)
		if err := k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, toAddr, coins); err != nil {
			return fmt.Errorf("failed to send coins: %v", err)
		}
		ctx.EventManager().EmitEvent(
			sdk.NewEvent("transfer",
				sdk.NewAttribute("from_module", types.ModuleName),
				sdk.NewAttribute("to", output.OutputTransfer.Address),
				sdk.NewAttribute("amount", output.OutputTransfer.Amount.String()),
			),
		)
		return nil

	case *types.ComponentOutput_OutputContractCall:
		// Correctly extract the Contract type component before executing
		contractComponent, ok := component.ComponentType.(*types.ExecutionComponent_Contract)
		if !ok {
			return fmt.Errorf("component is not a Contract component")
		}
		// contract caller could be beneficiary?
		if _, err := k.ExecuteContract(ctx, contractComponent, will.Creator); err != nil {
			return fmt.Errorf("failed to call contract: %v", err)
		}
		return nil

	case *types.ComponentOutput_OutputIbcSend:
		// Adjust to match the correct parameters and method definition
		// data := []byte(output.OutputIbcSend.Denom) // Simplistic assumption; adjust as needed!
		return k.SendIBCMessage(ctx, component, will)

	case *types.ComponentOutput_OutputEmit:
		ctx.EventManager().EmitEvent(
			sdk.NewEvent("emit_message",
				sdk.NewAttribute("message", output.OutputEmit.Message),
			),
		)
		return nil

	default:
		return fmt.Errorf("unsupported output type: %T", output)
	}
}

// AccessHandler checks if the signer has the right to submit a claim based on the component's access settings.
func (k Keeper) AccessHandler(ctx context.Context, component *types.ExecutionComponent, will types.Will, msg *types.MsgClaimRequest) error {
	// Type assertion to make sure we're working with a claim component
	claimComponent, ok := component.ComponentType.(*types.ExecutionComponent_Claim)
	if !ok {
		return fmt.Errorf("access control is only applicable to claim components")
	}

	access := claimComponent.Claim.Access

	switch acc := access.AccessType.(type) {
	case *types.ClaimAccessControl_Public:
		// If it's public, anyone can claim, so do nothing here
		fmt.Println("AccessHandler: Component access is public")

		return nil
	case *types.ClaimAccessControl_Private:
		// If private, check if the signer is in the allowed list
		for _, addr := range acc.Private.Addresses {
			if msg.Claimer == addr {
				fmt.Println("AccessHandler: Claimer is Authorized")
				return nil // Signer is authorized
			}
		}

		fmt.Println("AccessHandler: Claimer is NOT AUTHORIZED")

		return fmt.Errorf("signer %s is not authorized to claim this component", msg.Claimer)
	default:
		return fmt.Errorf("unsupported access type")
	}
}

// TODO: early claiming (im not a fan of this)
// at will creation, for each component, configure if it
// can be claimable early. If yes, whenever a beneficiary makes a claim
// we will check if the will component is early claimable.
// if so, and the verification is successful
// store will at a new key in the store [block_number]:early_claim

// //////////////// EXECUTE TRANSFER
// ExecuteTransferComponent handles the execution of a transfer component within a will.

// ExecuteTransfer executes a transfer component within a will.
func (k *Keeper) ExecuteTransfer(ctx sdk.Context, component *types.ExecutionComponent, will types.Will) error {
	// Check if the component is a TransferComponent
	transferComponent, ok := component.ComponentType.(*types.ExecutionComponent_Transfer)
	if !ok {
		return fmt.Errorf("component is not a TransferComponent")
	}

	// Prepare the coins for transfer
	coins := sdk.NewCoins(*transferComponent.Transfer.Amount)

	// Parse addresses
	// fromAddr, err := sdk.AccAddressFromBech32(transferComponent.Transfer.From)

	// fmt.Println(transferComponent.Transfer.From)
	fmt.Println(will.Creator)
	fmt.Println(transferComponent.Transfer.To)

	fromAddr, err := sdk.AccAddressFromBech32(will.Creator)
	if err != nil {
		return fmt.Errorf("parsing from address failed: %w", err)
	}

	toAddr, err := sdk.AccAddressFromBech32(transferComponent.Transfer.To)
	if err != nil {
		return fmt.Errorf("parsing to address failed: %w", err)
	}

	// fromAddr := sdk.AccAddress{}
	// toAddr := sdk.AccAddress{}

	// Use the bank module's SendCoinsFromAccountToAccount to transfer the coins
	fmt.Println("ExecuteTransfer")
	fmt.Println(fromAddr)
	fmt.Println(toAddr)

	fmt.Println(coins)
	if err := k.bankKeeper.SendCoins(ctx, fromAddr, toAddr, coins); err != nil {
		return fmt.Errorf("failed to send coins: %w", err)
	}

	// Optionally log the transfer event
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(component.Name,
			sdk.NewAttribute("will_id", will.ID),
			sdk.NewAttribute("from", will.Creator),
			sdk.NewAttribute("to", component.GetTransfer().To),
			sdk.NewAttribute("amount", coins.String()),
		),
	)

	return nil
}

// //////////////////////////// EXECUTE CONTRACT
// function to invoke contract during will execution, or claim
func (k Keeper) ExecuteContract(ctx sdk.Context, c *types.ExecutionComponent_Contract, caller string) ([]byte, error) {
	// Prepare the message you want to send to the contract. You might need to serialize it if it's not already in []byte form.
	msg := c.Contract.Data // Assuming this is already in []byte form.

	// Convert sdk.Context to context.Context. Be cautious with context conversions and make sure you're handling it correctly across your entire application.
	ctxContext := sdk.UnwrapSDKContext(ctx)

	// Prepare coins if your contract call requires sending tokens along. If not, just pass nil or an empty sdk.Coins{}.
	coins := sdk.NewCoins() // Assuming no coins are needed for this example.

	// Call the execute function. You need to replace "contractAddress" with the actual address of the contract and "caller" with the appropriate caller address.
	contractAddr, err := sdk.AccAddressFromBech32(c.Contract.Address)
	if err != nil {
		// handle error
		return nil, fmt.Errorf("failed to send coins: %w", err)
	}

	// this should be the will creator address, or parameterized
	// callerAddr := sdk.AccAddress{} // Determine how you get or set the caller address.
	callerAddr, callerErr := sdk.AccAddressFromBech32(caller)
	if callerErr != nil {
		// handle error
		return nil, fmt.Errorf("failed to send coins: %w", err)
	}

	// k.wasmKeeper.
	return k.permissionedWasmKeeper.Execute(ctxContext, contractAddr, callerAddr, msg, coins)
}

func (k Keeper) ExecutePrivateTransfer(ctx sdk.Context, component *types.ExecutionComponent) error {
	return nil
}

//////////////////////////////////////////////// IBC

// SendIBCMessage sends an IBC message from the specified port and channel with the given data
func (k *Keeper) SendIBCMessage(ctx sdk.Context, component *types.ExecutionComponent, will types.Will) error {
	fmt.Println("SendIBCMessage 1")
	var channelID, portID string
	var data []byte

	channelID = component.GetIbcMsg().Channel
	portID = component.GetIbcMsg().PortId
	data = component.GetIbcMsg().Data // Assuming Data is a field in IbcContractCall
	fmt.Println("send ibc messsage DEBUG CHANNEL KEEPER")
	fmt.Println(k.GetChannelKeeper())
	fmt.Println(k.channelKeeper)
	// panic(99)
	sequence, found := k.GetChannelKeeper().GetNextSequenceSend(ctx, portID, channelID)
	fmt.Println("will.keeper SendIBCMessage, sequence: ", sequence, " : ", found)
	if !found {
		fmt.Println("sequence not found for channel")
		return fmt.Errorf("sequence not found for channel")
	}

	timeoutHeight := clienttypes.NewHeight(0, 10000)
	fmt.Println("timeoutHeight")
	fmt.Println(timeoutHeight)

	timeoutTimestamp := uint64(ctx.BlockTime().UnixNano()) + 100000000000
	fmt.Println("timeoutTimestamp")
	fmt.Println(timeoutTimestamp)

	packet := channeltypes.NewPacket(data, sequence, "will", channelID, portID, channelID, timeoutHeight, timeoutTimestamp)
	fmt.Println("will.keeper SendIBCMessage, packet: ", packet)

	var capabilityName string = host.ChannelCapabilityPath(portID, channelID)
	fmt.Println("capabilityName: ", capabilityName)
	channelCap, ok := k.scopedKeeper.GetCapability(ctx, capabilityName)
	channelCap1, ok := k.scopedKeeper.GetCapability(ctx, capabilityName)
	fmt.Println("will.keeper SendIBCMessage, channelCap: ", channelCap, " : ", capabilityName)
	fmt.Println("will.keeper SendIBCMessage, channelCap1: ", channelCap1, " : ", capabilityName)

	/////////

	if !ok {
		return fmt.Errorf("channel capability not found")
	}

	////////////////
	fmt.Println("k.GetChannelKeeper().GetAllChannels()")
	fmt.Println(k.GetChannelKeeper().GetAllChannels(ctx))
	_, err := k.GetChannelKeeper().SendPacket(ctx, channelCap, portID, channelID, timeoutHeight, timeoutTimestamp, packet.GetData())
	return err
}

// hasCapability checks if the transfer module owns the port capability for the desired port
func (k *Keeper) hasCapability(ctx sdk.Context, portID string) bool {
	var portPath string = host.PortPath(portID)
	fmt.Println("portpath: %s", portPath)
	_, ok := k.scopedKeeper.GetCapability(ctx, portPath)
	return ok
}

// // dev for visibility
// func (k *Keeper) HasCapability(ctx sdk.Context, portID string) bool {
// 	var portPath string = host.PortPath(portID)
// 	fmt.Println("portpath 2: %s", portPath)

// 	_, ok := k.scopedKeeper.GetCapability(ctx, host.PortPath(portID))
// 	return ok
// }

// BindPort defines a wrapper function for the ort Keeper's function in
// order to expose it to module's InitGenesis function
func (k *Keeper) BindPort(ctx sdk.Context, portID string) error {
	capability := k.portKeeper.BindPort(ctx, portID)
	// ports/will
	var capabilityName string = host.PortPath(portID)
	fmt.Println("binding port: capability, ", capabilityName)
	fmt.Println(capability)
	// return k.ClaimCapability(ctx, capability, host.PortPath(portID))
	k.scopedKeeper.ClaimCapability(ctx, capability, capabilityName)
	// k.portKeeper.BindPort(ctx, portID)
	// k.capabilityKeeper.Cl
	k.ClaimCapability(ctx, capability, host.PortPath(portID))
	return nil
}

// GetPort returns the portID for the transfer module. Used in ExportGenesis
func (k *Keeper) GetPort(ctx sdk.Context) string {
	store := ctx.KVStore(k.storeKey)
	return string(store.Get(types.PortKey))
}

// SetPort sets the portID for the transfer module. Used in InitGenesis
func (k Keeper) SetPort(ctx sdk.Context, portID string) {
	// fmt.Println("SETTING PORT HERE: %s", portID)
	// fmt.Println("SETTING PORT HERE storekey: %s", k.storeKey)
	// store := ctx.KVStore(k.storeKey)
	store := k.storeService.OpenKVStore(ctx)
	store.Set(types.PortKey, []byte(portID))
}

// anteprotoreflect.ProtoMessage
func (k Keeper) VerifyWillAddress(ctx sdk.Context, msg sdk.Msg) (bool, error) {
	// func (k Keeper) VerifyWillAddress(ctx sdk.Context, msg anteprotoreflect.ProtoMessage) (bool, error) {
	fmt.Println("INSIDE WILL VERIFY ADDRESS")
	return true, nil
}

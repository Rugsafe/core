package keeper

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strconv"
	"time"

	capabilitykeeper "github.com/cosmos/ibc-go/modules/capability/keeper"
	clienttypes "github.com/cosmos/ibc-go/v8/modules/core/02-client/types" //nolint:staticcheck
	channeltypes "github.com/cosmos/ibc-go/v8/modules/core/04-channel/types"
	host "github.com/cosmos/ibc-go/v8/modules/core/24-host"
	"go.dedis.ch/kyber/v3/group/edwards25519"

	"cosmossdk.io/collections"
	corestoretypes "cosmossdk.io/core/store"
	"cosmossdk.io/errors"
	"cosmossdk.io/log"
	storetypes "cosmossdk.io/store/types"
	icatypes "github.com/cosmos/ibc-go/v8/modules/apps/27-interchain-accounts/types"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/telemetry"
	sdk "github.com/cosmos/cosmos-sdk/types"

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

// var (
// 	_ ibctypes.ChannelKeeper = (*Keeper)(nil)
// 	_ ibctypes.PortKeeper    = (*Keeper)(nil)
// )

type (
	Keeper struct {
		storeService corestoretypes.KVStoreService
		// storeService storetypes.KVStoreKey
		cdc           codec.Codec
		storeKey      storetypes.StoreKey // Add this line
		ChannelKeeper ChannelKeeper
		scopedKeeper  capabilitykeeper.ScopedKeeper
		portKeeper    PortKeeper

		// capabilityKeeper CapabilityKeeper
		capabilityKeeper capabilitykeeper.Keeper

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
) Keeper {
	fmt.Println("NewKeeper:")
	// sb := collections.NewSchemaBuilder(storeService)

	// sk := ScopedKeeper {
	// 		cdc      codec.BinaryCodec
	// 		storeKey storetypes.StoreKey
	// 		memKey   storetypes.StoreKey
	// 		capMap   map[uint64]*captypes.Capability
	// 		module   string
	// }

	fmt.Println("scopedKeeper:")
	fmt.Println(scopedKeeper)
	fmt.Println(scopedIBCKeeper)
	fmt.Println(capabilityKeeper)
	keeper := &Keeper{
		storeService:     storeService,
		cdc:              cdc,
		ChannelKeeper:    channelKeeper,
		portKeeper:       portKeeper,
		scopedKeeper:     scopedKeeper,
		capabilityKeeper: capabilityKeeper,
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

func createWillId(creator string, name string, beneficiary string, height int64) string {
	return fmt.Sprintf("%s-%s-%s-%s", creator, name, beneficiary, strconv.Itoa(int(height)))
}

func (k *Keeper) CreateWill(ctx context.Context, msg *types.MsgCreateWillRequest) (*types.Will, error) {
	store := k.storeService.OpenKVStore(ctx)

	// Concatenate values to generate a unique hash
	concatValues := createWillId(msg.Creator, msg.Name, msg.Beneficiary, msg.Height)
	idBytes := []byte(concatValues)

	// Generate a truncated hash of the concatenated values
	// truncatedHash, err := TruncateHash(idBytes, 16) // Truncate SHA256 to 16 bytes
	// if err != nil {
	// 	return nil, err
	// }

	// Convert the truncated hash bytes to a hex string for safe serialization
	idString := hex.EncodeToString(idBytes)
	fmt.Println(fmt.Printf("NEWLY CREATED WILL: %s", idString))

	// TODO: verify components, as this is already done in client/cli/tx.go
	// verifyComponents(msg.components)

	// Construct the will object
	will := types.Will{
		ID:          idString,
		Creator:     msg.Creator,
		Name:        msg.Name,
		Beneficiary: msg.Beneficiary,
		Height:      msg.Height,
		Status:      "live",
		Components:  msg.Components,
	}

	// Marshal the will object to bytes
	willBz := k.cdc.MustMarshal(&will)
	fmt.Println("inside k.createWill: " + idString)
	if willBz == nil {
		var errBz error
		return nil, errors.Wrap(errBz, "inside k.createWill, willBz is nil") // Make sure to handle the error appropriately
	}

	// Use the GetWillKey function to get a unique byte key for this will
	key := types.GetWillKey(idString)
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

// CLAIM
func (k Keeper) Claim(ctx context.Context, msg *types.MsgClaimRequest) error {
	// Retrieve the will by ID to ensure it exists and to process the claim against it
	fmt.Println("CLAIM FROM KEEPER: ")
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

	// Process the claim based on its type
	switch claim := msg.ClaimType.(type) {
	case *types.MsgClaimRequest_SchnorrClaim:

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

		// You can now check its status before proceeding with the claim.
		if will.Components[componentIndex].Status != "active" {
			fmt.Printf("component with ID %s is not active and cannot be claimed\n", msg.ComponentId)
			return fmt.Errorf("component with ID %s is not active and cannot be claimed", msg.ComponentId)
		}

		// Assuming the public key and signature are provided as byte slices in the claim
		fmt.Println(claim)

		// publicKeyBytes := claim.SchnorrClaim.PublicKey // The public key bytes
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

		fmt.Println("Schnorr signature verified successfully.")
		will.Components[componentIndex].Status = "claimed"

		// store will update with
		store := k.storeService.OpenKVStore(ctx)
		concatValues := createWillId(will.Creator, will.Name, will.Beneficiary, will.Height)
		willID := hex.EncodeToString([]byte(concatValues))
		key := types.GetWillKey(willID)
		fmt.Println(fmt.Printf("BEGIN BLOCKER WILL EXECUTED: %s", willID))

		willBz := k.cdc.MustMarshal(will)
		storeErr := store.Set(key, willBz)
		if storeErr != nil {
			return errors.Wrapf(storeErr, "schnorr claim error: could not save will ID with updated component status")
		}

		fmt.Println("Schnorr signature verified and saved now successfully.")
	case *types.MsgClaimRequest_PedersenClaim:
		// Process PedersenClaim
		fmt.Printf("Processing Pedersen claim with commitment: %x, blinding factor: %x, and value: %x\n", claim.PedersenClaim.Commitment, claim.PedersenClaim.BlindingFactor, claim.PedersenClaim.Value)
		// TODO

	case *types.MsgClaimRequest_GnarkClaim:
		// Process GnarkClaim
		fmt.Printf("Processing Gnark claim with proof: %x and public inputs: %x\n", claim.GnarkClaim.Proof, claim.GnarkClaim.PublicInputs)
		// TODO

	default:
		// Handle unknown claim type
		fmt.Println("unknown claim type provided")
		return fmt.Errorf("unknown claim type provided")
	}

	// Assuming the claim has been validated successfully, you can then update the will's status or components accordingly
	return nil
}

// expiration
func (k Keeper) SetWillExpiryIndex(ctx sdk.Context, expiryHeight int64, willID string) {
	store := ctx.KVStore(k.storeKey)
	expiryKey := []byte(fmt.Sprintf("expiry:%d:%s", expiryHeight, willID))
	store.Set(expiryKey, []byte(willID))
}

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

// BEGIN BLOCKER
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

		for component_index, component := range will.Components {
			fmt.Printf("Iterating over compnents for will ID %s for further processing.\n", will.ID)
			fmt.Println(component_index)
			fmt.Println(component)
			switch c := component.ComponentType.(type) {
			case *types.ExecutionComponent_Transfer:
				fmt.Printf("Transfer component found, to: %s, amount: %s\n", c.Transfer.To, c.Transfer.Amount.String())

				// TODO: actually execute the trades

				// update status to executed
				component.Status = "executed"

				// TODO: DEV TESTING FOR SENDIBCMESSAGE
				channelID := "testChannelID"
				portID := "testPortID"
				data := []byte("testData")
				k.SendIBCMessage(sdk.UnwrapSDKContext(ctx), channelID, portID, data)

			case *types.ExecutionComponent_Claim:
				fmt.Printf("Claim component found, evidence")
				// set all claimable components to active - can now have claims submitted
				component.Status = "active"
				// fmt.Printf("Claim component found, evidence: %s\n", c.Claim.Evidence)
			// case *types.ExecutionComponent_ContractCall:

			default:
				fmt.Println("Unknown component type found")
			}
		}

		fmt.Printf("Will ID: %s, Name: %s, Beneficiary: %s, Height: %d\n", will.ID, will.Name, will.Beneficiary, will.Height)

		// update will
		will.Status = "expired"
		concatValues := createWillId(will.Creator, will.Name, will.Beneficiary, will.Height)
		willID := hex.EncodeToString([]byte(concatValues))
		// willID := hex.EncodeToString(idString)
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

//////////////////////////////////////////////// IBC

// SendIBCMessage sends an IBC message from the specified port and channel with the given data
func (k *Keeper) SendIBCMessage(ctx sdk.Context, channelID, portID string, data []byte) error {
	// Retrieve the next sequence send for the channel
	fmt.Println("SendIBCMessage 1")
	sequence, found := k.ChannelKeeper.GetNextSequenceSend(ctx, portID, channelID)
	if !found {
		fmt.Println("SendIBCMessage 2")
		return errors.New("sequence not found for channel", 1, "k.channelKeeper.GetNextSequenceSend ran out")
	}

	// Define packet timeout (adjust as needed)
	fmt.Println("SendIBCMessage 3")
	timeoutHeight := clienttypes.NewHeight(0, 10000)                      // Use appropriate timeout height
	timeoutTimestamp := uint64(ctx.BlockTime().UnixNano()) + 100000000000 // 100 seconds; adjust as needed

	// Construct the packet
	packet := channeltypes.NewPacket(data, sequence, portID, channelID, "destPort", "destChannel", timeoutHeight, timeoutTimestamp)

	fmt.Println("SendIBCMessage 4")
	// Retrieve the capability for the port and channel
	channelCap, ok := k.scopedKeeper.GetCapability(ctx, host.ChannelCapabilityPath(portID, channelID))
	if !ok {
		fmt.Println("SendIBCMessage 5")
		return errors.New("channel capability not found: ", 1, "k.scopedKeeper.GetCapability ran out")
	}

	// Send the packet
	_, err := k.ChannelKeeper.SendPacket(ctx, channelCap, portID, channelID, timeoutHeight, timeoutTimestamp, packet.GetData())
	fmt.Println("SendIBCMessage 6")
	return err
}

// hasCapability checks if the transfer module owns the port capability for the desired port
func (k *Keeper) hasCapability(ctx sdk.Context, portID string) bool {
	var portPath string = host.PortPath(portID)
	fmt.Println("portpath: %s", portPath)
	_, ok := k.scopedKeeper.GetCapability(ctx, portPath)
	return ok
}

// dev for visibility
func (k *Keeper) HasCapability(ctx sdk.Context, portID string) bool {
	var portPath string = host.PortPath(portID)
	fmt.Println("portpath 2: %s", portPath)

	_, ok := k.scopedKeeper.GetCapability(ctx, host.PortPath(portID))
	return ok
}

// BindPort defines a wrapper function for the ort Keeper's function in
// order to expose it to module's InitGenesis function
func (k *Keeper) BindPort(ctx sdk.Context, portID string) error {
	capability := k.portKeeper.BindPort(ctx, portID)
	fmt.Println("capability")
	fmt.Println(capability)
	// return k.ClaimCapability(ctx, capability, host.PortPath(portID))
	return k.scopedKeeper.ClaimCapability(ctx, capability, host.PortPath(portID))
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

// SetPort sets the portID for the transfer module. Used in InitGenesis
// func (k *Keeper) SetPort(ctx sdk.Context, portID string) {
// 	// fmt.Println("SETTING PORT HERE: %s", portID)
// 	// fmt.Println("SETTING PORT HERE storekey: %s", k.storeKey)
// 	// store := ctx.KVStore(k.storeKey)

// 	store.Set(types.PortKey, []byte(portID))
// 	store := k.storeService.OpenKVStore(ctx)
// 	store.Set(expiryKey, []byte(willID))

// }

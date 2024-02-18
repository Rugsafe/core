package keeper

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"
	"strconv"
	"time"

	corestoretypes "cosmossdk.io/core/store"
	"cosmossdk.io/errors"
	"cosmossdk.io/log"
	storetypes "cosmossdk.io/store/types"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/telemetry"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/CosmWasm/wasmd/x/will/types"
)

type IKeeper interface {
	CreateWill(
		ctx context.Context,
		msg *types.MsgCreateWillRequest,
	) (*types.Will, error)
	GetWillByID(ctx context.Context, id string) (*types.Will, error)
	ListWillsByAddress(ctx context.Context, address string) ([]*types.Will, error)
}

type Keeper struct {
	storeService corestoretypes.KVStoreService
	// storeService storetypes.KVStoreKey
	cdc      codec.Codec
	storeKey storetypes.StoreKey // Add this line
}

func NewKeeper(
	cdc codec.Codec,
	storeService corestoretypes.KVStoreService,
	// storeService storetypes.KVStoreKey,
	logger log.Logger,
) Keeper {
	// sb := collections.NewSchemaBuilder(storeService)
	keeper := &Keeper{
		storeService: storeService,
		cdc:          cdc,
	}

	return *keeper
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
	return fmt.Sprintf("%s-%s-%s-%s", creator, name, beneficiary, height)
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

func (k Keeper) Claim(ctx sdk.Context, msg *types.MsgClaimRequest) {
	// claim specific component
}

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
			case *types.ExecutionComponent_Claim:
				fmt.Printf("Claim component found, evidence")
				// fmt.Printf("Claim component found, evidence: %s\n", c.Claim.Evidence)
				// case *types.ExecutionComponent_ContractCall:

			default:
				fmt.Println("Unknown component type found")
			}
		}

		fmt.Printf("Will ID: %s, Name: %s, Beneficiary: %s, Height: %d\n", will.ID, will.Name, will.Beneficiary, will.Height)

		// update will
		will.Status = "expired"
		// Store the marshaled will in the module's store
		// storeErr := store.Set(key, willBz)
		// concatValues := createWillId(will.Creator, will.Name, will.Beneficiary, will.Height)
		// idBytes := []byte(concatValues)
		// idString := hex.EncodeToString(idBytes)
		idString := hex.EncodeToString([]byte(will.ID))
		key := types.GetWillKey(idString)
		fmt.Println(fmt.Printf("BEGIN BLOCKER WILL EXECUTED: %s", idString))

		willBz := k.cdc.MustMarshal(will)
		storeErr := store.Set([]byte(key), willBz)

		if storeErr != nil {
			return errors.Wrapf(storeErr, "inside k.beginBlocker storeErr, KV store set threw an error after updating will: %s", will.ID)
		}

	}

	os.Exit(10)

	return nil
}

// func (k *Keeper) EndBlocker(ctx context.Context) error {
func (k *Keeper) EndBlocker(ctx sdk.Context) error {
	defer telemetry.ModuleMeasureSince(types.ModuleName, time.Now(), telemetry.MetricKeyEndBlocker)
	fmt.Println("INSIDE END BLOCKER FOR WILL MODULE")
	return nil
}

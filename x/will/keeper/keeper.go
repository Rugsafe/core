package keeper

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strconv"
	"strings"
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

func (k *Keeper) CreateWill(ctx context.Context, msg *types.MsgCreateWillRequest) (*types.Will, error) {
	store := k.storeService.OpenKVStore(ctx)

	// Concatenate values to generate a unique hash
	concatValues := fmt.Sprintf("%s-%s-%s", msg.Creator, msg.Name, msg.Beneficiary)
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
		Name:        msg.Name,
		Beneficiary: msg.Beneficiary,
		Height:      msg.Height,
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
	heightKey := types.GetWillKey(strconv.Itoa(int(will.Height)))
	existingWillsBz, existingWillsErr := store.Get(heightKey)

	if existingWillsErr != nil {
		return nil, errors.Wrap(existingWillsErr, "inside k.createWill existingWillsErr, KV store set threw an error")
	}

	var willsAtHeight []string
	if existingWillsBz != nil {
		// Deserialize existing wills if any
		existingWills := string(existingWillsBz)
		willsAtHeight = strings.Split(existingWills, ",")
	}

	// Add the new will ID to the set of wills at this height
	willsAtHeight = append(willsAtHeight, will.ID)

	// Serialize the updated set of wills and store it
	updatedWillsBz := []byte(strings.Join(willsAtHeight, ","))
	store.Set(heightKey, updatedWillsBz)

	/////////////// store for creator
	creatorKey := types.GetWillKey(msg.Creator)
	existingWillsForCreatorBz, existingWillsForCreatorErr := store.Get(creatorKey)

	if existingWillsForCreatorErr != nil {
		return nil, errors.Wrap(existingWillsForCreatorErr, "inside k.createWill existingWillsForCreatorErr, KV store set threw an error")
	}

	var willsAtCreator []string
	if existingWillsForCreatorBz != nil {
		existingWillsForCreator := string(existingWillsBz)
		willsAtCreator = strings.Split(existingWillsForCreator, ",")
	}

	// Add the new will ID to the set of wills at this height
	willsAtHeight = append(willsAtCreator, will.ID)

	return &will, nil
}

// func (k Keeper) ListWillsByAddress(ctx context.Context, address string) ([]*types.Will, error) {
// 	sdkCtx := sdk.UnwrapSDKContext(ctx)
// 	store := k.storeService.OpenKVStore(sdkCtx)

// 	fmt.Println("listWillsByAddress: Starting to list wills for address", address)

// 	// Use address to generate prefix for will keys.
// 	prefix := types.GetWillKey(address)
// 	// Iterator with prefix and nil as end parameter to get all keys with the prefix.
// 	iterator, _ := store.Iterator(prefix, nil)
// 	defer iterator.Close()

// 	var wills []*types.Will
// 	// var wills *types.Wills
// 	for ; iterator.Valid(); iterator.Next() {
// 		fmt.Println("listWillsByAddress: Found will in store")
// 		var will types.Will
// 		fmt.Println(iterator.Value())
// 		err := k.cdc.Unmarshal(iterator.Value(), &will)
// 		if err != nil {
// 			fmt.Println("listWillsByAddress: Error unmarshaling will", err)
// 			return nil, errors.Wrap(err, "failed to unmarshal will")
// 		}
// 		wills = append(wills, &will)
// 	}

//		fmt.Println("listWillsByAddress: Completed listing wills")
//		return wills, nil
//	}
func (k Keeper) ListWillsByAddress(ctx context.Context, address string) ([]*types.Will, error) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	store := k.storeService.OpenKVStore(sdkCtx)

	// Assuming you're storing a 'Wills' object under a specific key related to the address
	key := types.GetWillKey(address) // Ensure you have a method to generate a unique key for storing `Wills` by address

	bz, err := store.Get(key)
	if err != nil {
		fmt.Println("listWillsByAddress: Error fetching wills for address", err)
		return nil, errors.Wrap(err, "failed to fetch wills for address")
	}
	if bz == nil {
		// No wills found for this address
		return []*types.Will{}, nil
	}

	var wills types.Wills
	err = k.cdc.Unmarshal(bz, &wills)
	if err != nil {
		fmt.Println("listWillsByAddress: Error unmarshaling wills", err)
		return nil, errors.Wrap(err, "failed to unmarshal wills")
	}

	fmt.Println("listWillsByAddress: Completed listing wills")
	return wills.Wills, nil // Assuming 'Wills' type has a field named 'Wills' which is a slice of *Will
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

func (k Keeper) BeginBlocker(ctx sdk.Context) error {
	defer telemetry.ModuleMeasureSince(types.ModuleName, time.Now(), telemetry.MetricKeyBeginBlocker)

	// Get the current block height
	blockHeight := ctx.BlockHeight()
	fmt.Printf("Processing wills at block height: %d\n", blockHeight)

	// Access the store
	store := k.storeService.OpenKVStore(ctx)

	// Construct the height key to fetch will IDs associated with the current block height
	heightKey := types.GetWillKey(strconv.Itoa(int(blockHeight)))
	willIDsBz, err := store.Get(heightKey)

	// If there is an error fetching the will IDs or if there are no wills for this block height, return early
	if err != nil || willIDsBz == nil {
		fmt.Println("No wills to process for this block height or unable to fetch will IDs.")
		return nil // You might want to handle the error more gracefully depending on your application's needs
	}

	// Deserialize the list of will IDs
	willIDs := strings.Split(string(willIDsBz), ",")

	// Iterate over each will ID
	for _, willID := range willIDs {
		// Fetch the will object using its ID
		willBz, willFetchErr := store.Get(types.GetWillKey(willID))
		if willFetchErr != nil {
			fmt.Printf("Error fetching will with ID %s: %v\n", willID, willFetchErr)
			continue // Proceed to the next will if there's an issue fetching this one
		}

		var will types.Will
		unmarshalErr := k.cdc.Unmarshal(willBz, &will)
		if unmarshalErr != nil {
			fmt.Printf("Error unmarshaling will with ID %s: %v\n", willID, unmarshalErr)
			continue // Continue to the next will if unmarshaling fails
		}

		// Perform the desired operations on the will object here
		// For example, checking conditions, updating state, etc.
		fmt.Printf("Successfully fetched and unmarshaled will with ID %s for further processing.\n", will.ID)

		// Example operation: Print will details
		// Adjust this section based on what you actually need to do with each will
		fmt.Printf("Will ID: %s, Name: %s, Beneficiary: %s, Height: %d\n", will.ID, will.Name, will.Beneficiary, will.Height)
	}

	return nil
}

// func (k *Keeper) EndBlocker(ctx context.Context) error {
func (k *Keeper) EndBlocker(ctx sdk.Context) error {
	defer telemetry.ModuleMeasureSince(types.ModuleName, time.Now(), telemetry.MetricKeyEndBlocker)
	fmt.Println("INSIDE END BLOCKER FOR WILL MODULE")
	return nil
}

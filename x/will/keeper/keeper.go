package keeper

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"

	corestoretypes "cosmossdk.io/core/store"
	"cosmossdk.io/errors"
	"cosmossdk.io/log"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/CosmWasm/wasmd/x/will/types"
)

type Keeper struct {
	storeService corestoretypes.KVStoreService
	// storeService storetypes.KVStoreKey
	cdc codec.Codec
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

func (k Keeper) getWillByID(ctx context.Context, id string) (*types.Will, error) {
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

func (k Keeper) createWill(ctx context.Context, msg *types.MsgCreateWillRequest) (*types.Will, error) {
	store := k.storeService.OpenKVStore(ctx)

	// Concatenate values to generate a unique hash
	concatValues := fmt.Sprintf("%s-%s-%s", msg.Creator, msg.Name, msg.Beneficiary)
	idBytes := []byte(concatValues)

	// Generate a truncated hash of the concatenated values
	truncatedHash, err := TruncateHash(idBytes, 16) // Truncate SHA256 to 16 bytes
	if err != nil {
		return nil, err
	}

	// Convert the truncated hash bytes to a hex string for safe serialization
	idString := hex.EncodeToString(truncatedHash)
	fmt.Println(fmt.Printf("NEWLY CREATED WILL: %s", idString))

	// Construct the will object
	will := types.Will{
		ID:          idString,
		Name:        msg.Name,
		Beneficiary: msg.Beneficiary,
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
	storeErr := store.Set(key, willBz)

	if storeErr != nil {
		return nil, errors.Wrap(storeErr, "inside k.createWill, KV store set threw an error")
	}

	return &will, nil
}

// func (k Keeper) listWillsByAddress(ctx context.Context, address string) ([]*types.Will, error) {
// 	sdkCtx := sdk.UnwrapSDKContext(ctx)
// 	store := k.storeService.OpenKVStore(sdkCtx)

// 	var wills []*types.Will

// 	// Prefix for scanning the store
// 	prefix := types.GetWillKey(address)
// 	iterator := store.Iterator(prefix, sdk.PrefixEndBytes(prefix))
// 	defer iterator.Close()

// 	for ; iterator.Valid(); iterator.Next() {
// 		var will types.Will
// 		err := k.cdc.Unmarshal(iterator.Value(), &will)
// 		if err != nil {
// 			return nil, errors.Wrap(err, "failed to unmarshal will")
// 		}
// 		wills = append(wills, &will)
// 	}

// 	return wills, nil
// }

// func (k Keeper) listWillsByAddress(ctx context.Context, address string) ([]*types.Will, *query.PageResponse, error) {
// 	sdkCtx := sdk.UnwrapSDKContext(ctx)
// 	store := k.storeService.OpenKVStore(sdkCtx)

// 	// Assuming address is somehow part of the key, and GetWillAddressPrefix generates a prefix to scan all wills for that address.
// 	prefix := types.GetWillKey(address) // You need to implement this based on your key design.
// 	iterator, err := store.Iterator(prefix, sdk.PrefixEndBytes(prefix))
// 	if err != nil {
// 		return nil, nil, err // Correct error handling.
// 	}
// 	defer iterator.Close()

// 	var wills []*types.Will
// 	for ; iterator.Valid(); iterator.Next() {
// 		var will types.Will
// 		err := k.cdc.Unmarshal(iterator.Value(), &will)
// 		if err != nil {
// 			return nil, nil, errors.Wrap(err, "failed to unmarshal will")
// 		}
// 		wills = append(wills, &will)
// 	}

// 	// Assuming pagination is not implemented or needed, return nil for *query.PageResponse.
// 	// Implement pagination if necessary, adjusting PageResponse accordingly.
// 	return wills, nil, nil
// }

// func (k Keeper) listWillsByAddress(ctx context.Context, address string, pagination *query.PageRequest) ([]*types.Will, *query.PageResponse, error) {
// 	sdkCtx := sdk.UnwrapSDKContext(ctx)
// 	store := k.storeService.OpenKVStore(sdkCtx)

// 	// Prefix for scanning the store
// 	// prefix := types.GetWillKey(address) // Ensure this function exists or is implemented to convert address to prefix
// 	var wills []*types.Will

// 	pageRes, err := query.Paginate(store, pagination, func(key []byte, value []byte) error {
// 		var will types.Will
// 		if err := k.cdc.Unmarshal(value, &will); err != nil {
// 			return err
// 		}
// 		wills = append(wills, &will)
// 		return nil
// 	})

// 	if err != nil {
// 		return nil, nil, err
// 	}

// 	return wills, pageRes, nil
// }

// func (k Keeper) listWillsByAddress(ctx context.Context, address string) ([]*types.Will, error) {
// 	sdkCtx := sdk.UnwrapSDKContext(ctx)
// 	store := k.storeService.OpenKVStore(sdkCtx)

// 	// Assuming address is somehow part of the key to scan all wills for that address.
// 	prefix := types.GetWillKey(address)
// 	iterator, _ := store.Iterator(prefix, sdk.PrefixEndBytes(prefix))
// 	defer iterator.Close()

// 	var wills []*types.Will
// 	for ; iterator.Valid(); iterator.Next() {
// 		var will types.Will
// 		err := k.cdc.Unmarshal(iterator.Value(), &will)
// 		if err != nil {
// 			return nil, errors.Wrap(err, "failed to unmarshal will")
// 		}
// 		wills = append(wills, &will)
// 	}

// 	return wills, nil
// }

func (k Keeper) listWillsByAddress(ctx context.Context, address string) ([]*types.Will, error) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	store := k.storeService.OpenKVStore(sdkCtx)

	fmt.Println("listWillsByAddress: Starting to list wills for address", address)

	// Use address to generate prefix for will keys.
	prefix := types.GetWillKey(address)
	// Iterator with prefix and nil as end parameter to get all keys with the prefix.
	iterator, _ := store.Iterator(prefix, nil)
	defer iterator.Close()

	var wills []*types.Will
	for ; iterator.Valid(); iterator.Next() {
		fmt.Println("listWillsByAddress: Found will in store")
		var will types.Will
		err := k.cdc.Unmarshal(iterator.Value(), &will)
		if err != nil {
			fmt.Println("listWillsByAddress: Error unmarshaling will", err)
			return nil, errors.Wrap(err, "failed to unmarshal will")
		}
		wills = append(wills, &will)
	}

	fmt.Println("listWillsByAddress: Completed listing wills")
	return wills, nil
}

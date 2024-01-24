package keeper

import (
	"context"

	corestoretypes "cosmossdk.io/core/store"
	"cosmossdk.io/log"

	"github.com/cosmos/cosmos-sdk/codec"

	// "github.com/gogo/protobuf/codec"
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

func (k Keeper) GetWillByID(ctx context.Context, id uint64) *types.Will {
	store := k.storeService.OpenKVStore(ctx)
	var will types.Will
	willBz, _ := store.Get(types.GetWillKey(id))
	k.cdc.MustUnmarshal(willBz, &will)
	return &will
}

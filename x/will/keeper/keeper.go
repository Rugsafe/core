package keeper

import (
	corestoretypes "cosmossdk.io/core/store"
	"cosmossdk.io/log"

	// storetypes "cosmossdk.io/store/types"

	// "github.com/gogo/protobuf/codec"
	"github.com/cosmos/cosmos-sdk/codec"
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

	// keeper.wasmVMQueryHandler = DefaultQueryPlugins(bankKeeper, stakingKeeper, distrKeeper, channelKeeper, keeper)
	// preOpts, postOpts := splitOpts(opts)
	// for _, o := range preOpts {
	// 	o.apply(keeper)
	// }
	// // only set the wasmvm if no one set this in the options
	// // NewVM does a lot, so better not to create it and silently drop it.
	// if keeper.wasmVM == nil {
	// 	var err error
	// 	keeper.wasmVM, err = wasmvm.NewVM(filepath.Join(homeDir, "wasm"), availableCapabilities, contractMemoryLimit, wasmConfig.ContractDebugMode, wasmConfig.MemoryCacheSize)
	// 	if err != nil {
	// 		panic(err)
	// 	}
	// }

	// for _, o := range postOpts {
	// 	o.apply(keeper)
	// }
	// // not updatable, yet
	// keeper.wasmVMResponseHandler = NewDefaultWasmVMContractResponseHandler(NewMessageDispatcher(keeper.messenger, keeper))
	return *keeper
}

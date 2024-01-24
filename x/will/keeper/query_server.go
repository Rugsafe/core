package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/CosmWasm/wasmd/x/will/types"
)

var _ types.QueryServer = &queryServer{}

type queryServer struct {
	// cdc           codec.Codec
	// storeService  corestoretypes.KVStoreService
	keeper *Keeper
	// queryGasLimit storetypes.Gas
}

// NewGrpcQuerier constructor
//
//	func NewGrpcQuerier(cdc codec.Codec, storeService corestoretypes.KVStoreService, keeper types.ViewKeeper, queryGasLimit storetypes.Gas) *GrpcQuerier {
//		return &queryServer{cdc: cdc, storeService: storeService, keeper: keeper, queryGasLimit: queryGasLimit}
//	}
func NewGrpcQuerier(k *Keeper) types.QueryServer {
	return &queryServer{keeper: k}
}

// GetWill gets a will from the store
func (q queryServer) GetWill(c context.Context, req *types.QueryGetWillRequest) (*types.QueryGetWillResponse, error) {
	rsp, _ := queryGetWill(sdk.UnwrapSDKContext(c), 1, q.keeper)
	return rsp, nil
}

func queryGetWill(ctx sdk.Context, id uint64, keeper *Keeper) (*types.QueryGetWillResponse, error) {
	info := keeper.GetWillByID(ctx, id)
	// if info == nil {
	// 	return nil, types.ErrNoSuchContractFn(addr.String()).
	// 		Wrapf("address %s", addr.String())
	// }
	return &types.QueryGetWillResponse{
		WillId: id,
		Will:   info,
	}, nil
}

package keeper

import (
	"context"
	"fmt"

	"cosmossdk.io/errors"

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
func NewGrpcQuerier(k *Keeper) types.QueryServer {
	return &queryServer{keeper: k}
}

// GetWill gets a will from the store
func (q queryServer) GetWill(c context.Context, req *types.QueryGetWillRequest) (*types.QueryGetWillResponse, error) {
	rsp, _ := queryGetWill(sdk.UnwrapSDKContext(c), req.WillId, q.keeper)
	fmt.Println("rsp")
	fmt.Println(rsp)
	return rsp, nil
}

func queryGetWill(ctx sdk.Context, id string, keeper *Keeper) (*types.QueryGetWillResponse, error) {
	fmt.Println("QUERY SERVER.GO, Getting will by ID")
	will, err := keeper.GetWillByID(ctx, id)
	fmt.Println("=======QUERY GET WILL=======")
	fmt.Println(will)
	if err != nil {
		return nil, errors.Wrapf(err, "queryGetWill: error when getting will by id: %s", id)
	}
	return &types.QueryGetWillResponse{
		Will: will,
	}, nil
}

func queryListWills(ctx context.Context, keeper *Keeper, req *types.QueryListWillsRequest) (*types.QueryListWillsResponse, error) {
	fmt.Println("QUERY SERVER.GO, Getting wills by address", req.Address)
	wills, err := keeper.ListWillsByAddress(ctx, req.Address)
	if err != nil {
		fmt.Println("queryListWills: Error listing wills for address", req.Address, err)
		return nil, errors.Wrapf(err, "queryListWills: error when listing wills for address: %s", req.Address)
	}

	// Convert []*types.Will (Go slice of pointers) to []types.Will (protobuf slice)
	willsProto := make([]types.Will, len(wills))
	for i, will := range wills {
		willsProto[i] = *will // Dereference the pointer to match the protobuf expectation
	}

	fmt.Println("queryListWills: Successfully listed wills for address", req.Address)
	return &types.QueryListWillsResponse{
		Wills: willsProto, // Adjust to use the dereferenced slice
	}, nil
}

func (q queryServer) ListWills(ctx context.Context, req *types.QueryListWillsRequest) (*types.QueryListWillsResponse, error) {
	wills, err := q.keeper.ListWillsByAddress(ctx, req.Address)
	if err != nil {
		return nil, err
	}

	// Prepare the slice of Will for the protobuf response
	willsProto := make([]types.Will, len(wills))
	for i, will := range wills {
		willsProto[i] = *will // Dereference the pointer to match the protobuf expectation
	}

	return &types.QueryListWillsResponse{
		Wills: willsProto, // This now matches the protobuf definition
	}, nil
}

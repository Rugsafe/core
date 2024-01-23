package keeper

import (
	"context"

	"github.com/CosmWasm/wasmd/x/will/types"
)

var _ types.MsgServer = msgServer{}

type msgServer struct {
	keeper *Keeper
}

func NewMsgServerImpl(k *Keeper) types.MsgServer {
	return &msgServer{keeper: k}
}

func (m msgServer) CreateWill(
	ctx context.Context,
	msg *types.MsgCreateWill,
) (*types.MsgCreateWillResponse, error) {
	return &types.MsgCreateWillResponse{}, nil
}

// UpdateParams updates the module parameters
func (m msgServer) UpdateParams(ctx context.Context, req *types.MsgUpdateParams) (*types.MsgUpdateParamsResponse, error) {
	return &types.MsgUpdateParamsResponse{}, nil
}

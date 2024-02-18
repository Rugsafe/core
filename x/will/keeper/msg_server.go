package keeper

import (
	"context"
	"fmt"

	"cosmossdk.io/errors"

	"github.com/CosmWasm/wasmd/x/will/types"
)

var _ types.MsgServer = msgServer{}

type msgServer struct {
	keeper IKeeper
}

func NewMsgServerImpl(k IKeeper) types.MsgServer {
	return &msgServer{keeper: k}
}

func (m msgServer) CreateWill(
	ctx context.Context,
	msg *types.MsgCreateWillRequest,
) (*types.MsgCreateWillResponse, error) {
	fmt.Println("Inside msg_server, CreateWill")
	will, err := m.keeper.CreateWill(ctx, msg)
	if err != nil {
		return nil, errors.Wrap(err, "upon creating will")
	}
	return &types.MsgCreateWillResponse{
		Id:          will.ID,
		Creator:     msg.GetCreator(),
		Name:        will.Name,
		Beneficiary: will.Beneficiary,
		Height:      will.Height,
	}, nil
}

func (m msgServer) CheckIn(
	ctx context.Context,
	msg *types.MsgCheckInRequest,
) (*types.MsgCheckInResponse, error) {
	return &types.MsgCheckInResponse{}, nil
}

func (m msgServer) Claim(ctx context.Context, msg *types.MsgClaimRequest) (*types.MsgClaimResponse, error) {
	// m.keeper.Claim(ctx, msg)
	return &types.MsgClaimResponse{
		Success: true,
		Message: "Claim processed successfully",
	}, nil
}

// UpdateParams updates the module parameters
func (m msgServer) UpdateParams(ctx context.Context, req *types.MsgUpdateParams) (*types.MsgUpdateParamsResponse, error) {
	return &types.MsgUpdateParamsResponse{}, nil
}

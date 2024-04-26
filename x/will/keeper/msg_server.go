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
	// signer := sdk.AccAddress(ctx.Signers()[0].Bytes()).String()

	will, err := m.keeper.CreateWill(ctx, msg)
	fmt.Println("MSG SERVER CREATE WILL")
	fmt.Println(will)
	fmt.Println(err)

	if err != nil {
		// return nil, errors.Wrap(err, "error upon creating will")
		return nil, err
	} else {
		return &types.MsgCreateWillResponse{
			Id:          will.ID,
			Creator:     msg.GetCreator(),
			Name:        will.Name,
			Beneficiary: will.Beneficiary,
			Height:      will.Height,
		}, nil
	}

}

func (m msgServer) Claim(ctx context.Context, msg *types.MsgClaimRequest) (*types.MsgClaimResponse, error) {
	fmt.Println("INSIDE CLAIM FUNCTION")
	err := m.keeper.Claim(ctx, msg)
	if err != nil {
		return nil, errors.Wrap(err, "upon claiming will component")
	}
	return &types.MsgClaimResponse{
		Success: true,
		Message: "Claim processed successfully",
	}, nil
}

func (m msgServer) CheckIn(
	ctx context.Context,
	msg *types.MsgCheckInRequest,
) (*types.MsgCheckInResponse, error) {
	return &types.MsgCheckInResponse{}, nil
}

// UpdateParams updates the module parameters
func (m msgServer) UpdateParams(ctx context.Context, req *types.MsgUpdateParams) (*types.MsgUpdateParamsResponse, error) {
	return &types.MsgUpdateParamsResponse{}, nil
}

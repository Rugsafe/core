package keeper

import (
	"context"
	"fmt"

	"github.com/cockroachdb/errors"
	// willkeeper "github.com/CosmWasm/wasmd/x/will/keeper"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type Will interface {
	AddressCheck(ctx context.Context, address string) (bool, error)
}

type WillDecorator struct {
	willKeeper Keeper
}

func NewWillDecorator(w Keeper) WillDecorator {
	return WillDecorator{
		willKeeper: w,
	}
}

func (wd WillDecorator) AnteHandle(ctx sdk.Context, tx sdk.Tx, simulate bool, next sdk.AnteHandler) (sdk.Context, error) {
	// loop through all the messages and check if the message type is allowed
	fmt.Println("ANTE HANDLER TX")
	fmt.Println(tx)
	msgv2, _ := tx.GetMsgsV2()
	fmt.Println(msgv2)
	for _, msg := range tx.GetMsgs() {
		// address, _ := sdk.AccAddressFromBech32(msg)
		fmt.Println("WILL ANTE HANDLER")
		fmt.Println(msg)
		fmt.Println(ctx)
		ctx.TxBytes()
		isAllowed, err := wd.willKeeper.VerifyWillAddress(ctx, msg)
		if err != nil {
			return ctx, err
		}

		if !isAllowed {
			return ctx, errors.New("tx type not allowed")
		}
	}

	return next(ctx, tx, simulate)
}

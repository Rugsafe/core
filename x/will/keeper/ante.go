package keeper

import (
	"context"
	"fmt"

	// willkeeper "github.com/CosmWasm/wasmd/x/will/keeper"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	"github.com/cosmos/cosmos-sdk/crypto/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type Will interface {
	AddressCheck(ctx context.Context, address string) (bool, error)
}

type WillDecorator struct {
	willKeeper Keeper
	txConfig   client.TxConfig
}

func NewWillDecorator(w Keeper, txc client.TxConfig) WillDecorator {
	return WillDecorator{
		willKeeper: w,
		txConfig:   txc,
	}
}

func convertBytesToAddresses(signersBytes [][]byte) []string {
	var addresses []string
	for _, bytes := range signersBytes {
		var pubKey types.PubKey
		fmt.Println("len(bytes): ", len(bytes))
		switch len(bytes) {
		// case 32: // Assuming ed25519 key
		// pubKey = &ed25519.PubKey{Key: bytes}
		case 33: // Assuming compressed secp256k1 key
			pubKey = &secp256k1.PubKey{Key: bytes}

		default:
			fmt.Println("Unknown or unsupported key type")
			continue
		}

		fmt.Println("pubKey.Address()")
		fmt.Println(pubKey.Address())
		address := sdk.AccAddress(pubKey.Address()).String()
		addresses = append(addresses, address)
	}
	return addresses
}

func (wd WillDecorator) AnteHandle(ctx sdk.Context, tx sdk.Tx, simulate bool, next sdk.AnteHandler) (sdk.Context, error) {
	// loop through all the messages and check if the message type is allowed
	fmt.Println("ANTE HANDLER TX")
	// fmt.Println(tx.Signature)
	msgv2, _ := tx.GetMsgsV2()
	fmt.Println(msgv2)
	fmt.Println(ctx.TxBytes())
	// signers := tx.GetSigners()
	fmt.Println("=====================")
	// for _, msg := range tx.GetMsgs() {
	for _, msg := range msgv2 {
		// address, _ := sdk.AccAddressFromBech32(msg)
		fmt.Println("WILL ANTE HANDLER")
		fmt.Println(msg)
		fmt.Println(ctx)
		ctx.TxBytes()
		tx := wd.txConfig.TxDecoder()
		fmt.Println("TX")
		fmt.Println(tx)
		fmt.Println(wd.txConfig.TxJSONDecoder())
		signer, _ := wd.txConfig.SigningContext().GetSigners(msg)
		fmt.Println("signer")
		fmt.Println(signer[0])
		fmt.Println(convertBytesToAddresses(signer))
		fmt.Println("signer string")
		fmt.Println(wd.txConfig.SigningContext().AddressCodec().BytesToString(signer[0]))

		// isAllowed, err := wd.willKeeper.VerifyWillAddress(ctx, msg)
		// if err != nil {
		// 	return ctx, err
		// }

		// if !isAllowed {
		// 	return ctx, errors.New("tx type not allowed")
		// }
	}

	return next(ctx, tx, simulate)
}

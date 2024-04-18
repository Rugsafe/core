package keeper

import (
	"context"
	"fmt"

	// willkeeper "github.com/CosmWasm/wasmd/x/will/keeper"

	"github.com/CosmWasm/wasmd/x/will/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	sdk "github.com/cosmos/cosmos-sdk/types"

	// "google.golang.org/protobuf/reflect/protoreflect"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
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
		var pubKey cryptotypes.PubKey
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
	fmt.Println("ANTE HANDLER TX")
	msgs := tx.GetMsgs()
	msgsV2, _ := tx.GetMsgsV2()
	fmt.Println(msgs)
	fmt.Println(msgsV2)
	fmt.Println("===========")
	// if err != nil {
	// 	return ctx, err
	// }

	for i, msg := range msgs {
		switch msg := msg.(type) {
		case *types.MsgCreateWillRequest:
			fmt.Println("Processing MsgCreateWillRequest")
			creator := msg.Creator
			fmt.Println("Creator:", creator)

			// Continue with your existing logic
			signers, err := wd.txConfig.SigningContext().GetSigners(msgsV2[i])
			if err != nil {
				return ctx, err
			}
			signerAddressStr, err := wd.txConfig.SigningContext().AddressCodec().BytesToString(signers[0])
			if err != nil {
				return ctx, err
			}
			fmt.Println("Signer Address:", signerAddressStr)

			if creator != signerAddressStr {
				return ctx, sdkerrors.ErrAppConfig.Wrapf("signer address does not match the creator address")
			}
		default:
			fmt.Println("Received message is not of type MsgCreateWillRequest, skipping specific checks.")
		}
	}

	return next(ctx, tx, simulate)
}

package keeper

import (
	"context"
	"fmt"

	// willkeeper "github.com/CosmWasm/wasmd/x/will/keeper"

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
	msgs, err := tx.GetMsgsV2()
	if err != nil {
		return ctx, err
	}

	for _, msg := range msgs {
		fmt.Println("WILL ANTE HANDLER")
		message := msg.ProtoReflect()

		// Find the field descriptor for the "Creator" field
		fd := message.Descriptor().Fields().ByName("creator")
		if fd == nil {
			fmt.Println("Creator field not found")
			continue
		}

		// Get the value of the "Creator" field
		creator := message.Get(fd).String()
		fmt.Println("Creator:", creator)

		// Assuming you want to do something with the Creator value, such as checking against signers
		signers, err := wd.txConfig.SigningContext().GetSigners(msg)
		if err != nil {
			return ctx, err
		}
		signerAddressStr, err := wd.txConfig.SigningContext().AddressCodec().BytesToString(signers[0])
		if err != nil {
			return ctx, err
		}
		fmt.Println("Signer Address:", signerAddressStr)

		if creator != signerAddressStr {
			fmt.Println("Signer address does not match the creator address")
			return ctx, sdkerrors.ErrAppConfig.Wrapf("Signer address does not match the creator address")
		}
	}

	return next(ctx, tx, simulate)
}

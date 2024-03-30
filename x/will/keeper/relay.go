package keeper

import (
	"fmt"
	"time"

	// capabilitytypes "github.com/cosmos/ibc-go/modules/capability/types"
	capabilitytypes "github.com/cosmos/ibc-go/modules/capability/types"
	ibcexported "github.com/cosmos/ibc-go/v8/modules/core/exported"

	"github.com/cosmos/cosmos-sdk/telemetry"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/CosmWasm/wasmd/x/will/types"
)

func (k Keeper) OnAckPacket(
	ctx sdk.Context,
	// contractAddr sdk.AccAddress,
	msg types.IBCPacketAckMsg,
) error {
	defer telemetry.MeasureSince(time.Now(), "will", "contract", "ibc-ack-packet")
	fmt.Println("INSIDE OnAckPacket")
	// panic(1)
	return nil
}

func (k Keeper) OnCloseChannel(
	ctx sdk.Context,
	// contractAddr sdk.AccAddress,
	msg types.IBCChannelCloseMsg,
) error {
	defer telemetry.MeasureSince(time.Now(), "will", "contract", "ibc-close-channel")
	fmt.Println("INSIDE OnCloseChannel")
	// panic(1)
	return nil
}

func (k Keeper) OnConnectChannel(
	ctx sdk.Context,
	// contractAddr sdk.AccAddress,
	msg types.IBCChannelConnectMsg,
) error {
	defer telemetry.MeasureSince(time.Now(), "will", "contract", "ibc-connect-channel")
	fmt.Println("INSIDE OnConnectChannel")
	// panic(1)
	return nil
}

func (k Keeper) OnOpenChannel(
	ctx sdk.Context,
	// contractAddr sdk.AccAddress,
	msg types.IBCChannelOpenMsg,
) (string, error) {
	defer telemetry.MeasureSince(time.Now(), "wasm", "contract", "ibc-open-channel")
	fmt.Println("INSIDE OnOpenChannel")

	// panic(1)
	return "version", nil
}

func (k Keeper) AuthenticateCapability(ctx sdk.Context, cap *capabilitytypes.Capability, name string) bool {
	fmt.Println("INSIDE AuthenticateCapability")
	// panic(1)
	return k.scopedKeeper.AuthenticateCapability(ctx, cap, name)
}

// ClaimCapability allows the transfer module to claim a capability
// that IBC module passes to it
func (k Keeper) ClaimCapability(ctx sdk.Context, cap *capabilitytypes.Capability, name string) error {
	fmt.Println("INSIDE ClaimCapability")
	fmt.Println(k.capabilityKeeper)
	fmt.Println(cap)
	// panic(1)
	return k.scopedKeeper.ClaimCapability(ctx, cap, name)
}

// OnRecvPacket calls the contract to process the incoming IBC packet. The contract fully owns the data processing and
// returns the acknowledgement data for the chain level. This allows custom applications and protocols on top
// of IBC. Although it is recommended to use the standard acknowledgement envelope defined in
// https://github.com/cosmos/ics/tree/master/spec/ics-004-channel-and-packet-semantics#acknowledgement-envelope
//
// For more information see: https://github.com/cosmos/ics/tree/master/spec/ics-004-channel-and-packet-semantics#packet-flow--handling
func (k Keeper) OnRecvPacket(
	ctx sdk.Context,
	// contractAddr sdk.AccAddress,
	msg types.IBCPacketReceiveMsg,
) (ibcexported.Acknowledgement, error) {
	defer telemetry.MeasureSince(time.Now(), "wasm", "contract", "ibc-recv-packet")
	fmt.Println("INSIDE OnRecvPacket")
	// panic(1)
	return nil, nil
}

// OnTimeoutPacket calls the contract to let it know the packet was never received on the destination chain within
// the timeout boundaries.
// The contract should handle this on the application level and undo the original operation
func (k Keeper) OnTimeoutPacket(
	ctx sdk.Context,
	// contractAddr sdk.AccAddress,
	msg types.IBCPacketTimeoutMsg,
) error {
	defer telemetry.MeasureSince(time.Now(), "wasm", "contract", "ibc-timeout-packet")
	fmt.Println("INSIDE OnTimeoutPacket")
	// panic(1)
	return nil
}

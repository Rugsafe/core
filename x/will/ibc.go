package will

import (
	// willKeeper "github.com/CosmWasm/wasmd/x/will/keeper"
	// willtypes "github.com/CosmWasm/wasmd/x/will/types"
	"math"

	// wasmvmtypes "github.com/CosmWasm/wasmvm/types"
	capabilitytypes "github.com/cosmos/ibc-go/modules/capability/types"
	channeltypes "github.com/cosmos/ibc-go/v8/modules/core/04-channel/types"
	porttypes "github.com/cosmos/ibc-go/v8/modules/core/05-port/types"
	host "github.com/cosmos/ibc-go/v8/modules/core/24-host"
	ibcexported "github.com/cosmos/ibc-go/v8/modules/core/exported"

	errorsmod "cosmossdk.io/errors"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/CosmWasm/wasmd/x/wasm/types"

	willkeeper "github.com/CosmWasm/wasmd/x/will/keeper"
	willtypes "github.com/CosmWasm/wasmd/x/will/types"
)

var _ porttypes.IBCModule = IBCModule{}

// internal interface that is implemented by ibc middleware
type appVersionGetter interface {
	// GetAppVersion returns the application level version with all middleware data stripped out
	GetAppVersion(ctx sdk.Context, portID, channelID string) (string, bool)
}

// type IBCHandler struct {
// 	keeper           types.IBCContractKeeper
// 	channelKeeper    types.ChannelKeeper
// 	appVersionGetter appVersionGetter
// }

// func NewIBCHandler(k types.IBCContractKeeper, ck types.ChannelKeeper, vg appVersionGetter) IBCHandler {
// 	return IBCHandler{keeper: k, channelKeeper: ck, appVersionGetter: vg}
// }

// IBCModule implements the ICS26 interface for interchain accounts host chains
type IBCModule struct {
	keeper           willkeeper.Keeper //keeper.Keeper
	channelKeeper    types.ChannelKeeper
	appVersionGetter appVersionGetter
}

// NewIBCModule creates a new IBCModule given the associated keeper
func NewIBCModule(k willkeeper.Keeper, ck types.ChannelKeeper, vg appVersionGetter) IBCModule {
	return IBCModule{
		keeper:           k,
		channelKeeper:    ck,
		appVersionGetter: vg,
	}
}

// OnAcknowledgementPacket implements the IBCModule interface
func (i IBCModule) OnAcknowledgementPacket(
	ctx sdk.Context,
	packet channeltypes.Packet,
	acknowledgement []byte,
	relayer sdk.AccAddress,
) error {
	// contractAddr, err := keeper.ContractFromPortID(packet.SourcePort)
	// if err != nil {
	// 	return errorsmod.Wrapf(err, "contract port id")
	// }

	// err = i.keeper.OnAckPacket(ctx, contractAddr, willtypes.IBCPacketAckMsg{
	// 	Acknowledgement: willtypes.IBCAcknowledgement{Data: acknowledgement},
	// 	OriginalPacket:  newIBCPacket(packet),
	// 	Relayer:         relayer.String(),
	// })
	// if err != nil {
	// 	return errorsmod.Wrap(err, "on ack")
	// }
	return nil
}

func newIBCPacket(packet channeltypes.Packet) willtypes.IBCPacket {
	timeout := willtypes.IBCTimeout{
		Timestamp: packet.TimeoutTimestamp,
	}
	if !packet.TimeoutHeight.IsZero() {
		timeout.Block = &willtypes.IBCTimeoutBlock{
			Height:   packet.TimeoutHeight.RevisionHeight,
			Revision: packet.TimeoutHeight.RevisionNumber,
		}
	}

	return willtypes.IBCPacket{
		Data:     packet.Data,
		Src:      willtypes.IBCEndpoint{ChannelID: packet.SourceChannel, PortID: packet.SourcePort},
		Dest:     willtypes.IBCEndpoint{ChannelID: packet.DestinationChannel, PortID: packet.DestinationPort},
		Sequence: packet.Sequence,
		Timeout:  timeout,
	}
}

// OnChanCloseConfirm implements the IBCModule interface
func (i IBCModule) OnChanCloseConfirm(ctx sdk.Context, portID, channelID string) error {
	// // counterparty has closed the channel
	// contractAddr, err := keeper.ContractFromPortID(portID)
	// if err != nil {
	// 	return errorsmod.Wrapf(err, "contract port id")
	// }
	// channelInfo, ok := i.channelKeeper.GetChannel(ctx, portID, channelID)
	// if !ok {
	// 	return errorsmod.Wrapf(channeltypes.ErrChannelNotFound, "port ID (%s) channel ID (%s)", portID, channelID)
	// }
	// appVersion, ok := i.appVersionGetter.GetAppVersion(ctx, portID, channelID)
	// if !ok {
	// 	return errorsmod.Wrapf(channeltypes.ErrInvalidChannelVersion, "port ID (%s) channel ID (%s)", portID, channelID)
	// }

	// msg := willtypes.IBCChannelCloseMsg{
	// 	CloseConfirm: &willtypes.IBCCloseConfirm{Channel: toWasmVMChannel(portID, channelID, channelInfo, appVersion)},
	// }
	// err = i.keeper.OnCloseChannel(ctx, contractAddr, msg)
	// if err != nil {
	// 	return err
	// }
	// // emit events?

	// return err
	return nil
}

// OnChanCloseInit implements the IBCModule interface
func (i IBCModule) OnChanCloseInit(ctx sdk.Context, portID, channelID string) error {
	// contractAddr, err := keeper.ContractFromPortID(portID)
	// if err != nil {
	// 	return errorsmod.Wrapf(err, "contract port id")
	// }
	// channelInfo, ok := i.channelKeeper.GetChannel(ctx, portID, channelID)
	// if !ok {
	// 	return errorsmod.Wrapf(channeltypes.ErrChannelNotFound, "port ID (%s) channel ID (%s)", portID, channelID)
	// }
	// appVersion, ok := i.appVersionGetter.GetAppVersion(ctx, portID, channelID)
	// if !ok {
	// 	return errorsmod.Wrapf(channeltypes.ErrInvalidChannelVersion, "port ID (%s) channel ID (%s)", portID, channelID)
	// }

	// msg := willtypes.IBCChannelCloseMsg{
	// 	CloseInit: &willtypes.IBCCloseInit{Channel: toWasmVMChannel(portID, channelID, channelInfo, appVersion)},
	// }
	// err = i.keeper.OnCloseChannel(ctx, contractAddr, msg)
	// if err != nil {
	// 	return err
	// }
	// // emit events?

	// return err
	return nil
}

// OnChanOpenAck implements the IBCModule interface
func (i IBCModule) OnChanOpenAck(
	ctx sdk.Context,
	portID, channelID string,
	counterpartyChannelID string,
	counterpartyVersion string,
) error {
	contractAddr, err := sdk.AccAddressFromBech32("abc") //keeper.ContractFromPortID(portID)
	if err != nil {
		return errorsmod.Wrapf(err, "address errored out")
	}
	channelInfo, ok := i.channelKeeper.GetChannel(ctx, portID, channelID)
	if !ok {
		return errorsmod.Wrapf(channeltypes.ErrChannelNotFound, "port ID (%s) channel ID (%s)", portID, channelID)
	}
	channelInfo.Counterparty.ChannelId = counterpartyChannelID

	appVersion, ok := i.appVersionGetter.GetAppVersion(ctx, portID, channelID)
	if !ok {
		return errorsmod.Wrapf(channeltypes.ErrInvalidChannelVersion, "port ID (%s) channel ID (%s)", portID, channelID)
	}

	msg := willtypes.IBCChannelConnectMsg{
		OpenAck: &willtypes.IBCOpenAck{
			Channel:             toWasmVMChannel(portID, channelID, channelInfo, appVersion),
			CounterpartyVersion: counterpartyVersion,
		},
	}
	return i.keeper.OnConnectChannel(ctx, contractAddr, msg)
}

// OnChanOpenConfirm implements the IBCModule interface
func (i IBCModule) OnChanOpenConfirm(ctx sdk.Context, portID, channelID string) error {
	contractAddr, err := sdk.AccAddressFromBech32(portID)

	if err != nil {
		return errorsmod.Wrapf(err, "contract port id")
	}
	channelInfo, ok := i.channelKeeper.GetChannel(ctx, portID, channelID)
	if !ok {
		return errorsmod.Wrapf(channeltypes.ErrChannelNotFound, "port ID (%s) channel ID (%s)", portID, channelID)
	}
	appVersion, ok := i.appVersionGetter.GetAppVersion(ctx, portID, channelID)
	if !ok {
		return errorsmod.Wrapf(channeltypes.ErrInvalidChannelVersion, "port ID (%s) channel ID (%s)", portID, channelID)
	}
	msg := willtypes.IBCChannelConnectMsg{
		OpenConfirm: &willtypes.IBCOpenConfirm{
			Channel: toWasmVMChannel(portID, channelID, channelInfo, appVersion),
		},
	}
	return i.keeper.OnConnectChannel(ctx, contractAddr, msg)
}

// OnChanOpenInit implements the IBCModule interface
func (i IBCModule) OnChanOpenInit(
	ctx sdk.Context,
	order channeltypes.Order,
	connectionHops []string,
	portID string,
	channelID string,
	chanCap *capabilitytypes.Capability,
	counterParty channeltypes.Counterparty,
	version string,
) (string, error) {
	// ensure port, version, capability
	if err := ValidateChannelParams(channelID); err != nil {
		return "", err
	}
	contractAddr, err := sdk.AccAddressFromBech32(portID)
	if err != nil {
		return "", errorsmod.Wrapf(err, "contract port id")
	}

	msg := willtypes.IBCChannelOpenMsg{
		OpenInit: &willtypes.IBCOpenInit{
			Channel: willtypes.IBCChannel{
				Endpoint:             willtypes.IBCEndpoint{PortID: portID, ChannelID: channelID},
				CounterpartyEndpoint: willtypes.IBCEndpoint{PortID: counterParty.PortId, ChannelID: counterParty.ChannelId},
				Order:                order.String(),
				// DESIGN V3: this may be "" ??
				Version:      version,
				ConnectionID: connectionHops[0], // At the moment this list must be of length 1. In the future multi-hop channels may be supported.
			},
		},
	}

	// Allow contracts to return a version (or default to proposed version if unset)
	acceptedVersion, err := i.keeper.OnOpenChannel(ctx, contractAddr, msg)
	if err != nil {
		return "", err
	}
	if acceptedVersion == "" { // accept incoming version when nothing returned by contract
		if version == "" {
			return "", types.ErrEmpty.Wrap("version")
		}
		acceptedVersion = version
	}

	// Claim channel capability passed back by IBC module
	if err := i.keeper.ClaimCapability(ctx, chanCap, host.ChannelCapabilityPath(portID, channelID)); err != nil {
		return "", errorsmod.Wrap(err, "claim capability")
	}
	return acceptedVersion, nil
}

// OnChanOpenTry implements the IBCModule interface
func (i IBCModule) OnChanOpenTry(
	ctx sdk.Context,
	order channeltypes.Order,
	connectionHops []string,
	portID, channelID string,
	chanCap *capabilitytypes.Capability,
	counterParty channeltypes.Counterparty,
	counterpartyVersion string,
) (string, error) {
	// ensure port, version, capability
	if err := ValidateChannelParams(channelID); err != nil {
		return "", err
	}

	contractAddr, err := sdk.AccAddressFromBech32(portID)
	if err != nil {
		return "", errorsmod.Wrapf(err, "contract port id")
	}

	msg := willtypes.IBCChannelOpenMsg{
		OpenTry: &willtypes.IBCOpenTry{
			Channel: willtypes.IBCChannel{
				Endpoint:             willtypes.IBCEndpoint{PortID: portID, ChannelID: channelID},
				CounterpartyEndpoint: willtypes.IBCEndpoint{PortID: counterParty.PortId, ChannelID: counterParty.ChannelId},
				Order:                order.String(),
				Version:              counterpartyVersion,
				ConnectionID:         connectionHops[0], // At the moment this list must be of length 1. In the future multi-hop channels may be supported.
			},
			CounterpartyVersion: counterpartyVersion,
		},
	}

	// Allow contracts to return a version (or default to counterpartyVersion if unset)
	version, err := i.keeper.OnOpenChannel(ctx, contractAddr, msg)
	if err != nil {
		return "", err
	}
	if version == "" {
		version = counterpartyVersion
	}

	// Module may have already claimed capability in OnChanOpenInit in the case of crossing hellos
	// (ie chainA and chainB both call ChanOpenInit before one of them calls ChanOpenTry)
	// If module can already authenticate the capability then module already owns it, so we don't need to claim
	// Otherwise, module does not have channel capability, and we must claim it from IBC
	if !i.keeper.AuthenticateCapability(ctx, chanCap, host.ChannelCapabilityPath(portID, channelID)) {
		// Only claim channel capability passed back by IBC module if we do not already own it
		if err := i.keeper.ClaimCapability(ctx, chanCap, host.ChannelCapabilityPath(portID, channelID)); err != nil {
			return "", errorsmod.Wrap(err, "claim capability")
		}
	}

	return version, nil
}

// OnRecvPacket implements the IBCModule interface
func (i IBCModule) OnRecvPacket(
	ctx sdk.Context,
	packet channeltypes.Packet,
	relayer sdk.AccAddress,
) ibcexported.Acknowledgement {
	contractAddr, err := sdk.AccAddressFromBech32(packet.DestinationPort)
	if err != nil {
		// this must not happen as ports were registered before
		panic(errorsmod.Wrapf(err, "contract port id"))
	}

	em := sdk.NewEventManager()
	msg := willtypes.IBCPacketReceiveMsg{Packet: newIBCPacket(packet), Relayer: relayer.String()}
	ack, err := i.keeper.OnRecvPacket(ctx.WithEventManager(em), contractAddr, msg)
	if err != nil {
		ack = channeltypes.NewErrorAcknowledgement(err)
		// the state gets reverted, so we drop all captured events
	} else if ack == nil || ack.Success() {
		// emit all contract and submessage events on success
		// nil ack is a success case, see: https://github.com/cosmos/ibc-go/blob/v7.0.0/modules/core/keeper/msg_server.go#L453
		ctx.EventManager().EmitEvents(em.Events())
	}
	types.EmitAcknowledgementEvent(ctx, contractAddr, ack, err)
	return ack
}

// OnTimeoutPacket implements the IBCModule interface
func (i IBCModule) OnTimeoutPacket(ctx sdk.Context, packet channeltypes.Packet, relayer sdk.AccAddress) error {
	contractAddr, err := sdk.AccAddressFromBech32(packet.SourcePort)
	if err != nil {
		return errorsmod.Wrapf(err, "contract port id")
	}
	msg := willtypes.IBCPacketTimeoutMsg{Packet: newIBCPacket(packet), Relayer: relayer.String()}
	err = i.keeper.OnTimeoutPacket(ctx, contractAddr, msg)
	if err != nil {
		return errorsmod.Wrap(err, "on timeout")
	}
	return nil
}

// helpers
func toWasmVMChannel(portID, channelID string, channelInfo channeltypes.Channel, appVersion string) willtypes.IBCChannel {
	return willtypes.IBCChannel{
		Endpoint:             willtypes.IBCEndpoint{PortID: portID, ChannelID: channelID},
		CounterpartyEndpoint: willtypes.IBCEndpoint{PortID: channelInfo.Counterparty.PortId, ChannelID: channelInfo.Counterparty.ChannelId},
		Order:                channelInfo.Ordering.String(),
		Version:              appVersion,
		ConnectionID:         channelInfo.ConnectionHops[0], // At the moment this list must be of length 1. In the future multi-hop channels may be supported.
	}
}

func ValidateChannelParams(channelID string) error {
	// NOTE: for escrow address security only 2^32 channels are allowed to be created
	// Issue: https://github.com/cosmos/cosmos-sdk/issues/7737
	channelSequence, err := channeltypes.ParseChannelSequence(channelID)
	if err != nil {
		return err
	}
	if channelSequence > math.MaxUint32 {
		return errorsmod.Wrapf(types.ErrMaxIBCChannels, "channel sequence %d is greater than max allowed transfer channels %d", channelSequence, math.MaxUint32)
	}
	return nil
}
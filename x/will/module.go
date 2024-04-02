package will

import (
	"context"
	"encoding/json"
	"fmt"

	// "github.com/CosmWasm/wasmd/x/wasm/types"
	abci "github.com/cometbft/cometbft/abci/types"
	// "github.com/gogo/protobuf/codec"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/spf13/cobra"

	"cosmossdk.io/errors"
	"cosmossdk.io/log"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"

	// "github.com/CosmWasm/wasmd/x/wasm/keeper"
	"github.com/CosmWasm/wasmd/x/wasm/simulation"
	"github.com/CosmWasm/wasmd/x/will/client/cli"
	"github.com/CosmWasm/wasmd/x/will/keeper"
	"github.com/CosmWasm/wasmd/x/will/types"
)

var (
	_ module.AppModuleBasic      = AppModuleBasic{}
	_ module.AppModuleSimulation = AppModule{}
)

type AppModuleBasic struct{}

type AppModule struct {
	AppModuleBasic
	cdc    codec.Codec
	keeper *keeper.Keeper
}

func NewAppModule(
	cdc codec.Codec,
	keeper *keeper.Keeper,
	logger log.Logger,
) AppModule {
	Main()
	return AppModule{
		AppModuleBasic: AppModuleBasic{},
		cdc:            cdc,
		keeper:         keeper,
	}
}

// Name returns the wasm module's name.
func (AppModuleBasic) Name() string {
	return types.ModuleName
}

func (am AppModule) IsAppModule() {
}

// IsOnePerModuleType implements the depinject.OnePerModuleType interface.
func (am AppModule) IsOnePerModuleType() { // marker
}

func (b AppModuleBasic) RegisterGRPCGatewayRoutes(clientCtx client.Context, serveMux *runtime.ServeMux) {
	// err := types.RegisterQueryHandlerClient(context.Background(), serveMux, types.NewQueryClient(clientCtx))
	// if err != nil {
	// 	panic(err)
	// }
}

// RegisterInterfaces implements InterfaceModule
func (b AppModuleBasic) RegisterInterfaces(registry cdctypes.InterfaceRegistry) {
	types.RegisterInterfaces(registry)
}

func (b AppModuleBasic) RegisterLegacyAminoCodec(amino *codec.LegacyAmino) {
	// types.RegisterLegacyAminoCodec(amino)
}

func (am AppModule) RegisterServices(cfg module.Configurator) {
	types.RegisterMsgServer(cfg.MsgServer(), keeper.NewMsgServerImpl(am.keeper))
	types.RegisterQueryServer(cfg.QueryServer(), keeper.NewGrpcQuerier(am.keeper))

	// for migrations!
	// m := keeper.NewMigrator(*am.keeper, am.legacySubspace)
	// err := cfg.RegisterMigration(types.ModuleName, 1, m.Migrate1to2)
	// if err != nil {
	// 	panic(err)
	// }
}

// GenerateGenesisState creates a randomized GenState of the bank module.
func (AppModule) GenerateGenesisState(simState *module.SimulationState) {
	simulation.RandomizedGenState(simState)
}

// ProposalMsgs returns msgs used for governance proposals for simulations.
// func (am AppModule) ProposalMsgs(simState module.SimulationState) []simtypes.WeightedProposalMsg {
// 	return simulation.ProposalMsgs(am.bankKeeper, am.keeper)
// }

// RegisterStoreDecoder registers a decoder for supply module's types
func (am AppModule) RegisterStoreDecoder(_ simtypes.StoreDecoderRegistry) {
}

// // WeightedOperations returns the all the gov module operations with their respective weights.
func (am AppModule) WeightedOperations(simState module.SimulationState) []simtypes.WeightedOperation {
	// return simulation.WeightedOperations(simState.AppParams, am.accountKeeper, am.bankKeeper, am.keeper)
	return nil
}

//////////////////////////

// EndBlock returns the end blocker for the delay module.
// func (am AppModule) EndBlock(_ client.Context) []abci.ValidatorUpdate {
// 	return []abci.ValidatorUpdate{}
// }

// func (am AppModule) EndBlock(ctx context.Context) ([]abci.ValidatorUpdate, error) {
// 	return am.keeper.EndBlocker(ctx)
// }

// BeginBlock executes delayed items.
func (am AppModuleBasic) BeginBlock(sdk sdk.Context) {
	fmt.Println("NOW IM ACTUALLY IN THE WILL MODULES BEGIN BLOCKER 1")
}

func (am AppModule) BeginBlock(ctx context.Context) error {
	fmt.Println("NOW IM ACTUALLY IN THE WILL MODULES BEGIN BLOCKER 2")
	fmt.Println(ctx)
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	denom, _ := sdk.GetBaseDenom()
	denom_log, _ := fmt.Printf("Base Denom %s", denom)
	fmt.Println(denom_log)
	// sdk.NewSearchBlocksResult()
	endBlockerError := am.keeper.BeginBlocker(sdkCtx)
	if endBlockerError != nil {
		error_msg, _ := fmt.Printf("ERROR RUNNING will.am.keeper.BeginBlock: %s", endBlockerError)
		fmt.Println(error_msg)
	} else {
	}
	return nil
}

func (am AppModule) EndBlock(ctx context.Context) error {
	fmt.Println("NOW IM ACTUALLY IN THE WILL MODULES END BLOCKER")
	fmt.Println(ctx)
	return nil
}

/////////////////////

// func (AppModuleBasic) RegisterInterfaces(registry codectypes.InterfaceRegistry) {
//     // Register all your message types and interfaces here
//     types.RegisterMsgServer(registry, &_Msg_serviceDesc)
// }

// will structue the will module for cli
func (b AppModuleBasic) GetTxCmd() *cobra.Command {
	return cli.GetTxCmd()
}

// GetQueryCmd returns no root query command for the wasm module.
func (b AppModuleBasic) GetQueryCmd() *cobra.Command {
	return cli.GetQueryCmd()
}

///////////////////////////////

// // func (am AppModule) InitGenesis(ctx sdk.Context, cdc codec.JSONCodec, data json.RawMessage) []abci.ValidatorUpdate {
// func (am AppModule) InitGenesis(ctx sdk.Context, data json.RawMessage) []abci.ValidatorUpdate {
// 	fmt.Println("ABOUT TO INVOKE INITGENESIS")
// 	var genesisState types.GenesisState
// 	// cdc.MustUnmarshalJSON(data, &genesisState)
// 	// am.keeper.InitGenesis(ctx, genesisState)
// 	_, _ = keeper.InitGenesis(ctx, am.keeper, genesisState)
// 	panic(1)
// 	// if err != nil {
// 	// 	panic(err)
// 	// }
// 	// return validators
// 	return nil
// }

// func (b AppModuleBasic) DefaultGenesis(cdc codec.JSONCodec) json.RawMessage {
// 	fmt.Println("INSIDE WILL DEFAULT GENESIS")
// 	return cdc.MustMarshalJSON(&types.GenesisState{
// 		Params: types.DefaultParams(),
// 		PortId: b.Name(),
// 	})
// }

// func (b AppModule) DefaultGenesis(cdc codec.JSONCodec) json.RawMessage {
// 	fmt.Println("INSIDE WILL DEFAULT GENESIS")
// 	return cdc.MustMarshalJSON(&types.GenesisState{
// 		Params: types.DefaultParams(),
// 		PortId: b.Name(),
// 	})
// }

// func (am AppModule) ExportGenesis(ctx sdk.Context, cdc codec.JSONCodec) json.RawMessage {
// 	gs := keeper.ExportGenesis(ctx, am.keeper)
// 	return cdc.MustMarshalJSON(gs)
// }

// // ValidateGenesis performs genesis state validation for the delay module.
// func (am AppModule) ValidateGenesis(cdc codec.JSONCodec, _ client.TxEncodingConfig, bz json.RawMessage) error {
// 	var genState types.GenesisState
// 	if err := cdc.UnmarshalJSON(bz, &genState); err != nil {
// 		return errors.Wrapf(err, "failed to unmarshal %s genesis state", types.ModuleName)
// 	}
// 	return nil //genState.Validate()
// }

//////////////////////////

// func (AppModuleBasic) DefaultGenesis(cdc codec.JSONCodec) json.RawMessage {
//     return cdc.MustMarshalJSON(types.DefaultGenesis())
// }

// DefaultGenesis returns default genesis state as raw bytes for the delay module.
func (amb AppModuleBasic) DefaultGenesis(cdc codec.JSONCodec) json.RawMessage {
	fmt.Println("INSIDE WILL DEFAULT GENESIS")
	return cdc.MustMarshalJSON(&types.GenesisState{
		Params: types.DefaultParams(),
		PortId: amb.Name(),
	})
}

// ValidateGenesis performs genesis state validation for the delay module.
func (AppModuleBasic) ValidateGenesis(cdc codec.JSONCodec, _ client.TxEncodingConfig, bz json.RawMessage) error {
	var genState types.GenesisState
	if err := cdc.UnmarshalJSON(bz, &genState); err != nil {
		return errors.Wrapf(err, "failed to unmarshal %s genesis state", types.ModuleName)
	}
	return nil
}

// InitGenesis performs genesis initialization for the delay module.
func (am AppModule) InitGenesis(ctx sdk.Context, cdc codec.JSONCodec, gs json.RawMessage) []abci.ValidatorUpdate {
	var genState types.GenesisState
	cdc.MustUnmarshalJSON(gs, &genState)

	// keeper.InitGenesis(ctx, cdc, gs)
	keeper.InitGenesis(ctx, am.keeper, genState)

	return []abci.ValidatorUpdate{}
}

// ExportGenesis returns the exported genesis state as raw bytes for the delay module.
func (am AppModule) ExportGenesis(ctx sdk.Context, cdc codec.JSONCodec) json.RawMessage {
	gs := keeper.ExportGenesis(ctx, am.keeper)
	return cdc.MustMarshalJSON(gs)
}

func Main() {
	fmt.Println("INSIDE WILL MODULE.GO")
}

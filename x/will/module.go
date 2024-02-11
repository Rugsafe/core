package will

import (
	"context"
	"fmt"

	// "github.com/gogo/protobuf/codec"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/spf13/cobra"

	"cosmossdk.io/log"

	// "github.com/CosmWasm/wasmd/x/wasm/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"

	// "github.com/CosmWasm/wasmd/x/wasm/keeper"
	"github.com/CosmWasm/wasmd/x/will/client/cli"
	"github.com/CosmWasm/wasmd/x/will/keeper"
	"github.com/CosmWasm/wasmd/x/will/types"
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
	endBlockerError := am.keeper.EndBlocker(sdkCtx)
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

func Main() {
	fmt.Println("INSIDE WILL MODULE.GO")
}

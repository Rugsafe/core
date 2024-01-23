package will

import (
	"fmt"

	// "github.com/gogo/protobuf/codec"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/spf13/cobra"

	"cosmossdk.io/log"

	// "github.com/CosmWasm/wasmd/x/wasm/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"

	// "github.com/CosmWasm/wasmd/x/wasm/keeper"
	"github.com/CosmWasm/wasmd/x/will/client/cli"
	"github.com/CosmWasm/wasmd/x/will/keeper"
	"github.com/CosmWasm/wasmd/x/will/types"

	"github.com/cosmos/cosmos-sdk/types/module"
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
	logger.Info("yo im here......")
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
	// types.RegisterQueryServer(cfg.QueryServer(), keeper.Querier(am.keeper))

	// for migrations!
	// m := keeper.NewMigrator(*am.keeper, am.legacySubspace)
	// err := cfg.RegisterMigration(types.ModuleName, 1, m.Migrate1to2)
	// if err != nil {
	// 	panic(err)
	// }

}

// func (AppModuleBasic) RegisterInterfaces(registry codectypes.InterfaceRegistry) {
//     // Register all your message types and interfaces here
//     types.RegisterMsgServer(registry, &_Msg_serviceDesc)
// }

// will structue the will module for cli
func (b AppModuleBasic) GetTxCmd() *cobra.Command {
	return cli.GetTxCmd()
}

func Main() {
	fmt.Println("INSIDE WILL MODULE.GO")
}

package will

import (
	"fmt"

	"cosmossdk.io/log"
	// "github.com/CosmWasm/wasmd/x/wasm/keeper"
	"github.com/CosmWasm/wasmd/x/will/keeper"

	// "github.com/CosmWasm/wasmd/x/wasm/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"

	// "github.com/gogo/protobuf/codec"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"

	"github.com/CosmWasm/wasmd/x/will/types"
)

type AppModuleBasic struct {
}

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
	// types.RegisterInterfaces(registry)
}

func (b AppModuleBasic) RegisterLegacyAminoCodec(amino *codec.LegacyAmino) {
	// types.RegisterLegacyAminoCodec(amino)
}

func Main() {
	fmt.Println("lol")

}

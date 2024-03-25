package will

// import (
// 	"encoding/json"
// 	"fmt"

// 	"github.com/CosmWasm/wasmd/x/will/types"
// 	abci "github.com/cometbft/cometbft/abci/types"
// 	"github.com/cosmos/cosmos-sdk/codec"
// 	sdk "github.com/cosmos/cosmos-sdk/types"
// )

// func (am AppModule) InitGenesis(ctx sdk.Context, cdc codec.JSONCodec, data json.RawMessage) []abci.ValidatorUpdate {
// 	var genesisState types.GenesisState
// 	cdc.MustUnmarshalJSON(data, &genesisState)
// 	fmt.Println("ABOUT TO INVOKE INITGENESIS")
// 	panic(1)
// 	am.keeper.InitGenesis(ctx, genesisState)
// 	// if err != nil {
// 	// 	panic(err)
// 	// }
// 	// return validators
// 	return nil
// }

// func (am AppModule) ExportGenesis(ctx sdk.Context, cdc codec.JSONCodec) json.RawMessage {
// 	gs := keeper.ExportGenesis(ctx, am.keeper)
// 	return cdc.MustMarshalJSON(gs)
// }

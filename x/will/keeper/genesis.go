package keeper

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"

	// "github.com/cosmos/ibc-go/v8/modules/apps/transfer/types"
	"github.com/CosmWasm/wasmd/x/will/types"
	abci "github.com/cometbft/cometbft/abci/types"
)

// InitGenesis initializes the ibc-transfer state and binds to PortID.
// func (k *Keeper) InitGenesis(ctx sdk.Context, state types.GenesisState) ([]abci.ValidatorUpdate, error) {
func InitGenesis(ctx sdk.Context, k *Keeper, state types.GenesisState) ([]abci.ValidatorUpdate, error) {
	fmt.Println("ABOUT TO BIND PORT IN INITGENESIS")

	k.SetPort(ctx, state.PortId)
	// Only try to bind to port if it is not already bound, since we may already own
	// port capability from capability InitGenesis
	if !k.hasCapability(ctx, state.PortId) {
		// transfer module binds to the transfer port on InitChain
		// and claims the returned capability
		err := k.BindPort(ctx, state.PortId)
		if err != nil {
			panic(fmt.Errorf("could not claim port capability: %v", err))
		} else {
			fmt.Println("INITGENESIS PORT BIND WILL SUCCESSFULL")
		}
	}
	// panic(2)
	//set params
	k.SetParams(ctx, state.Params)
	fmt.Println("AFTER SET PARAMS")
	return nil, nil
}

// ExportGenesis returns a GenesisState for a given context and keeper.
func ExportGenesis(ctx sdk.Context, keeper *Keeper) *types.GenesisState {
	var genState types.GenesisState
	return &genState
}

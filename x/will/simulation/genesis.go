package simulation

import (
	"github.com/cosmos/cosmos-sdk/types/module"

	"github.com/CosmWasm/wasmd/x/wasm/types"
)

// RandomizeGenState generates a random GenesisState for wasm
func RandomizedGenState(simstate *module.SimulationState) {
	params := types.DefaultParams()
	willGenesis := types.GenesisState{
		Params: params,
	}

	_, err := simstate.Cdc.MarshalJSON(&willGenesis)
	if err != nil {
		panic(err)
	}

	simstate.GenState[types.ModuleName] = simstate.Cdc.MustMarshalJSON(&willGenesis)
}

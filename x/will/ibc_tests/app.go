// Copyright Tharsis Labs Ltd.(Evmos)
// SPDX-License-Identifier:ENCL-1.0(https://github.com/evmos/evmos/blob/main/LICENSE)

package ibctesting

import (
	"encoding/json"
	"os"
	"testing"
	"time"

	"cosmossdk.io/math"
	abci "github.com/cometbft/cometbft/abci/types"
	"github.com/cometbft/cometbft/libs/log"
	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"
	tmtypes "github.com/cometbft/cometbft/types"

	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	cryptocodec "github.com/cosmos/cosmos-sdk/crypto/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/stretchr/testify/require"

	// ibcgotesting "github.com/cosmos/ibc-go/v7/testing"
	dbm "github.com/cosmos/cosmos-db"
	ibcgotesting "github.com/cosmos/ibc-go/v8/testing"

	// evmosapp "github.com/evmos/evmos/v16/app"

	// "github.com/evmos/evmos/v16/types"
	// "github.com/evmos/evmos/v16/utils"

	// ibctesting "github.com/cosmos/ibc-go/testing"

	simapp "github.com/cosmos/ibc-go/v8/testing/simapp"
	simappparams "github.com/cosmos/ibc-go/v8/testing/simapp/params"
	//	simappparams "github.com/cosmos/ibc-go/v8/testing/simapp/params"
)

var _ ibcgotesting.TestingApp = &simapp.SimApp{}

type AppOptions map[string]interface{}

func (ao AppOptions) Get(key string) interface{} {
	return ao[key]
}

// DefaultTestingAppInit is a test helper function used to initialize an App
// on the ibc testing pkg
// need this design to make it compatible with the SetupTestinApp func on ibctesting pkg

func SetupTestingApp() (ibcgotesting.TestingApp, map[string]json.RawMessage) {
	db := dbm.NewMemDB()
	encCfg := simappparams.MakeTestEncodingConfig()

	app := simapp.NewSimApp(log.NewTMLogger(os.Stdout), db, nil, true, map[int64]bool{}, "/tmp", 0, encCfg, simapp.EmptyAppOptions{})

	return app, simapp.NewDefaultGenesisState(encCfg.Marshaler)
}

var DefaultTestingAppInit func(chainID string) func() (ibcgotesting.TestingApp, map[string]json.RawMessage) = SetupTestingApp

var DefaultConsensusParams = &tmproto.ConsensusParams{
	Block: &tmproto.BlockParams{
		MaxBytes: 200000,
		MaxGas:   -1, // no limit
	},
	Evidence: &tmproto.EvidenceParams{
		MaxAgeNumBlocks: 302400,
		MaxAgeDuration:  504 * time.Hour, // 3 weeks is the max duration
		MaxBytes:        10000,
	},
	Validator: &tmproto.ValidatorParams{
		PubKeyTypes: []string{
			tmtypes.ABCIPubKeyTypeEd25519,
		},
	},
}

// SetupWithGenesisValSet initializes a new SimApp with a validator set and genesis accounts
// that also act as delegators. For simplicity, each validator is bonded with a delegation
// of one consensus engine unit (10^6) in the default token of the simapp from first genesis
// account. A Nop logger is set in SimApp.
func SetupWithGenesisValSet(t *testing.T, valSet *tmtypes.ValidatorSet, genAccs []authtypes.GenesisAccount, chainID string, balances ...banktypes.Balance) ibcgotesting.TestingApp {
	app, genesisState := DefaultTestingAppInit(chainID)()
	// set genesis accounts
	authGenesis := authtypes.NewGenesisState(authtypes.DefaultParams(), genAccs)
	genesisState[authtypes.ModuleName] = app.AppCodec().MustMarshalJSON(authGenesis)

	validators := make([]stakingtypes.Validator, 0, len(valSet.Validators))
	delegations := make([]stakingtypes.Delegation, 0, len(valSet.Validators))

	// bondAmt := sdk.TokensFromConsensusPower(1, types.PowerReduction)
	// powerred := 1000
	powerred := math.NewInt(1000)

	bondAmt := sdk.TokensFromConsensusPower(1, powerred)

	for _, val := range valSet.Validators {
		pk, err := cryptocodec.FromTmPubKeyInterface(val.PubKey)
		require.NoError(t, err)
		pkAny, err := codectypes.NewAnyWithValue(pk)
		require.NoError(t, err)
		validator := stakingtypes.Validator{
			OperatorAddress:   sdk.ValAddress(val.Address).String(),
			ConsensusPubkey:   pkAny,
			Jailed:            false,
			Status:            stakingtypes.Bonded,
			Tokens:            bondAmt,
			DelegatorShares:   math.LegacyOneDec(),
			Description:       stakingtypes.Description{},
			UnbondingHeight:   int64(0),
			UnbondingTime:     time.Unix(0, 0).UTC(),
			Commission:        stakingtypes.NewCommission(math.LegacyZeroDec(), math.LegacyZeroDec(), math.LegacyZeroDec()),
			MinSelfDelegation: math.ZeroInt(),
		}
		validators = append(validators, validator)
		// delegations = append(delegations, stakingtypes.NewDelegation(genAccs[0].GetAddress(), val.Address.Bytes(), math.LegacyOneDec()))
		delegations = append(delegations, stakingtypes.NewDelegation("ADDRESS=w3ll1p0k8gygawzpggzwftv7cv47zvgg8zaun5xucxz", "ADDRESS=w3ll1p0k8gygawzpggzwftv7cv47zvgg8zaun5xucxz", math.LegacyOneDec()))
	}

	// set validators and delegations
	stakingParams := stakingtypes.DefaultParams()
	// set bond demon to be aevmos
	stakingParams.BondDenom = "stake" //utils.BaseDenom
	stakingGenesis := stakingtypes.NewGenesisState(stakingParams, validators, delegations)
	genesisState[stakingtypes.ModuleName] = app.AppCodec().MustMarshalJSON(stakingGenesis)

	totalSupply := sdk.NewCoins()
	for _, b := range balances {
		// add genesis acc tokens and delegated tokens to total supply
		// totalSupply = totalSupply.Add(b.Coins.Add(sdk.NewCoin(utils.BaseDenom, bondAmt))...)
		totalSupply = totalSupply.Add(b.Coins.Add(sdk.NewCoin("stake", bondAmt))...)
	}

	// add bonded amount to bonded pool module account
	balances = append(balances, banktypes.Balance{
		Address: authtypes.NewModuleAddress(stakingtypes.BondedPoolName).String(),
		// Coins:   sdk.Coins{sdk.NewCoin(utils.BaseDenom, bondAmt)},
		Coins: sdk.Coins{sdk.NewCoin("stake", bondAmt)},
	})

	// update total supply
	bankGenesis := banktypes.NewGenesisState(banktypes.DefaultGenesisState().Params, balances, totalSupply, []banktypes.Metadata{}, []banktypes.SendEnabled{})
	genesisState[banktypes.ModuleName] = app.AppCodec().MustMarshalJSON(bankGenesis)

	stateBytes, err := json.MarshalIndent(genesisState, "", " ")
	require.NoError(t, err)

	// init chain will set the validator set and initialize the genesis accounts
	app.InitChain(
		&abci.RequestInitChain{
			ChainId:         chainID,
			Validators:      []abci.ValidatorUpdate{},
			ConsensusParams: DefaultConsensusParams,
			AppStateBytes:   stateBytes,
		},
	)

	// commit genesis changes
	app.Commit()

	app.BeginBlocker(abci.RequestBeginBlock{Header: tmproto.Header{
		ChainID:            chainID,
		Height:             app.LastBlockHeight() + 1,
		AppHash:            app.LastCommitID().Hash,
		ValidatorsHash:     valSet.Hash(),
		NextValidatorsHash: valSet.Hash(),
	}})

	return app
}

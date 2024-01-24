package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	ModuleName = "will"

	// SubModuleName defines the interchain accounts controller module name
	SubModuleName = "will"

	// StoreKey is the store key string for the interchain accounts controller module
	StoreKey = SubModuleName

	// ParamsKey is the store key for the interchain accounts controller parameters
	ParamsKey = "params"
)

var WillPrefix = []byte{0x01}

func GetWillKey(willID uint64) []byte {
	willID8z := sdk.Uint64ToBigEndian(willID)
	return append(WillPrefix, willID8z...)
}

/*
var ints []int32 = []int32{1, 2}
fmt.Println(ints[0])
var _map map[int]int = map[int]int{1: 10}
print(_map[1])
*/

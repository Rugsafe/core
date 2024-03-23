package types

import "strings"

const (
	ModuleName = "will"

	// SubModuleName defines the interchain accounts controller module name
	SubModuleName = "will"

	// StoreKey is the store key string for the interchain accounts controller module
	StoreKey = SubModuleName

	// ParamsKey is the store key for the interchain accounts controller parameters
	ParamsKey = "params"
)

var (
	WillPrefix = []byte{0x01}
	// PortKey defines the key to store the port ID in store
	PortKey = []byte{0x02}
)

func GetWillKey(willID string) []byte {
	// willID8z := sdk.Uint64ToBigEndian(willID)
	stringKey := []byte(strings.ToLower(willID))
	// return append(WillPrefix, willID8z...)
	// return append(WillPrefix, willID...)
	return append(WillPrefix, stringKey...)
}

/*
var ints []int32 = []int32{1, 2}
fmt.Println(ints[0])
var _map map[int]int = map[int]int{1: 10}
print(_map[1])
*/

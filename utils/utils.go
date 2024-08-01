package utils

import (
	"fmt"
	"math/big"
)

func FloatToBigInt(val float64, decimals int) *big.Int {
	bigval := new(big.Float).SetFloat64(val)
	bigval.SetString(fmt.Sprintf("%fE%d", val, decimals))

	result := new(big.Int)
	bigval.Int(result)

	return result
}

func MergeMaps[K comparable, V any](maps ...map[K]V) map[K]V {
	result := make(map[K]V)
	for _, m := range maps {
		for k, v := range m {
			result[k] = v
		}
	}
	return result
}

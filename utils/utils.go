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

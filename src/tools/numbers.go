package tools

import (
	"crypto/rand"
	"math/big"
)

// RandomInt generate a random int between 0 and max
func RandomInt(max int) int {
	BigIntMax := big.NewInt(int64(max))
	i, _ := rand.Int(rand.Reader, BigIntMax)
	iInt64 := i.Int64()
	iInt := int(iInt64)

	return iInt
}

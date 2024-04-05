/*
 * Permutate an array using entropy provided by <seed>. Original typescript
 * implementation from Jordi Baylina @ www.github.com/jbaylina/random_permute/blob/main/test/test.js
 *
 * Replaced (arr []int) with [0,...,n-1]
 */
package client

import (
	"math/big"
)

func Permutate(seed *big.Int, n int) []int {
	arr := make([]int, n)
	for i := 0; i < n; i++ {
		arr[i] = i
	}

	maxVal := big.NewInt(1)
	maxVal.Lsh(maxVal, 250)

	seed.And(seed, maxVal.Sub(maxVal, big.NewInt(1)))

	for i := n; i > 0; i-- {
		mod := new(big.Int).Set(seed)
		mod.Mod(mod, big.NewInt(int64(i)))
		r := int(mod.Int64())

		arr[i-1], arr[r] = arr[r], arr[i-1]

		seed.Sub(seed, mod)
		seed.Div(seed, big.NewInt(int64(i)))
	}

	return arr
}

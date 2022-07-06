package random

import (
	"crypto/rand"
	"math/big"
)

// RangeInt Return: A random range of numbers as a slice within a range
//
// Usage:
// arr := random.RangeInt(2, 100, 3)
//
// Output: [40, 5, 81]
// Output: [12, 52, 31]
func RangeInt(min int64, max int64, n int) []int64 {
	arr := make([]int64, n)
	var r int
	for r = 0; r <= n-1; r++ {
		n, _ := rand.Int(rand.Reader, big.NewInt(max))
		arr[r] = n.Int64() + min
	}
	return arr
}

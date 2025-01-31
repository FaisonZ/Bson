package nums

import (
	"fmt"
	"math"
	"slices"
)

func IsInt(n float64) bool {
	return n == math.Trunc(n)
}

func MinIntSize(n int64) (int, error) {
	foundSize := 0
	sizes := []struct {
		high int64
		low  int64
		size int
	}{
		{high: math.MaxInt8, low: math.MinInt8, size: 8},
		{high: math.MaxInt16, low: math.MinInt16, size: 16},
		{high: math.MaxInt32, low: math.MinInt32, size: 32},
		{high: math.MaxInt64, low: math.MinInt64, size: 64},
	}

	for _, s := range sizes {
		if n >= s.low && n <= s.high {
			foundSize = s.size
			break
		}
	}

	if foundSize == 0 {
		return 0, fmt.Errorf("Valid minimum int size not found")
	}

	return foundSize, nil
}

func IntFitsInSize(n float64, s int) (bool, error) {
	INT_SIZES := []int{8, 16, 32, 64}

	if !IsInt(n) {
		return false, fmt.Errorf("Must pass in a Float without a decimal")
	}

	if !slices.Contains(INT_SIZES, s) {
		return false, fmt.Errorf("Invalid int size: %d", s)
	}

	high := math.Pow(2, float64(s-1)) - 1
	low := -math.Pow(2, float64(s-1))

	return n >= low && n <= high, nil
}

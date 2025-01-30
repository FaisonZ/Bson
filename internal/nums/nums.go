package nums

import (
	"fmt"
	"math"
	"slices"
)

func IsInt(n float64) bool {
	return n == math.Trunc(n)
}

func MinIntSize(n float64) (int, error) {
	INT_SIZES := []int{8, 16, 32, 64}

	if !IsInt(n) {
		return 8, fmt.Errorf("Must pass in a Float without a decimal")
	}

	foundSize := 0

	for _, s := range INT_SIZES {
		high := math.Pow(2, float64(s-1)) - 1
		low := -math.Pow(2, float64(s-1))

		if n >= low && n <= high {
			foundSize = s
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

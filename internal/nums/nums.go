package nums

import "math"

func IsInt(n float64) bool {
	return n == math.Trunc(n)
}

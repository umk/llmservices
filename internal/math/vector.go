package math

import (
	"math"

	"github.com/pehringer/simd"
)

func VectorsEq[T comparable](a, b []T) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

// VectorNorm calculates the Euclidean norm (magnitude) of a vector.
// It uses SIMD operations for performance optimization when available.
// The tmp slice is used for intermediate calculations and should be at least
// the same length as the vector.
func VectorNorm(vector []float32, tmp []float32) float64 {
	if tmp == nil {
		tmp = make([]float32, len(vector))
	} else if len(tmp) < len(vector) {
		panic("tmp slice size is less than vector length")
	}

	simd.MulFloat32(vector, vector, tmp)

	var sum float32
	for _, v := range tmp {
		sum += v
	}

	return math.Sqrt(float64(sum))
}

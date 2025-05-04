package math

import (
	"github.com/pehringer/simd"
)

func CosineSimilarity(a, b []float32, normA, normB float64, tmp []float32) float64 {
	simd.MulFloat32(a, b, tmp)

	var sum float32
	for _, v := range tmp {
		sum += v
	}

	return float64(sum) / (normA * normB)
}

func ApproxEq(a, b, tolerance float64) bool {
	d := a - b
	if d < 0 {
		d = -d
	}
	return d < tolerance
}

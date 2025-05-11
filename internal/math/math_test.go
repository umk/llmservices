package math

import (
	"math"
	"testing"
)

func TestCosineSimilarity(t *testing.T) {
	tests := []struct {
		name string
		a    []float32
		b    []float32
		want float64
	}{
		{
			name: "Parallel vectors",
			a:    []float32{1, 2, 3},
			b:    []float32{2, 4, 6},
			want: 1.0,
		},
		{
			name: "Orthogonal vectors",
			a:    []float32{1, 0, 0},
			b:    []float32{0, 1, 0},
			want: 0.0,
		},
		{
			name: "Opposite vectors",
			a:    []float32{1, 2, 3},
			b:    []float32{-1, -2, -3},
			want: -1.0,
		},
		{
			name: "Arbitrary vectors",
			a:    []float32{1, 2, 3},
			b:    []float32{4, 5, 6},
			// dot product / (|a| * |b|)
			want: 32.0 / (math.Sqrt(14) * math.Sqrt(77)),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmp := make([]float32, len(tt.a))

			// Calculate norms
			normA := VectorNorm(tt.a, tmp)
			normB := VectorNorm(tt.b, tmp)

			got := CosineSimilarity(tt.a, tt.b, normA, normB, tmp)

			if math.Abs(got-tt.want) > 1e-6 {
				t.Errorf("cosineSimilarity() = %v, want %v", got, tt.want)
			}
		})
	}
}

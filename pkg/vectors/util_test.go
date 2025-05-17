package vectors

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestSimilarity tests the Similarity function with various vector pairs
func TestSimilarity(t *testing.T) {
	tests := []struct {
		name string
		v1   Vector
		v2   Vector
		want float32
	}{
		{
			name: "Identical vectors",
			v1:   Vector{1, 2, 3},
			v2:   Vector{1, 2, 3},
			want: 1.0,
		},
		{
			name: "Parallel vectors",
			v1:   Vector{1, 2, 3},
			v2:   Vector{2, 4, 6},
			want: 1.0,
		},
		{
			name: "Orthogonal vectors",
			v1:   Vector{1, 0, 0},
			v2:   Vector{0, 1, 0},
			want: 0.0,
		},
		{
			name: "Opposite vectors",
			v1:   Vector{1, 2, 3},
			v2:   Vector{-1, -2, -3},
			want: -1.0,
		},
		{
			name: "Arbitrary vectors",
			v1:   Vector{1, 2, 3},
			v2:   Vector{4, 5, 6},
			want: 0.9746318, // Approximate value calculated externally
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Similarity(tt.v1, tt.v2)
			assert.InDelta(t, float64(tt.want), float64(got), 1e-6)
		})
	}
}

// TestSimilarityPanic tests that the Similarity function panics
// when vectors of different lengths are provided
func TestSimilarityPanic(t *testing.T) {
	v1 := Vector{1, 2, 3}
	v2 := Vector{1, 2}

	// This should panic
	assert.Panics(t, func() {
		Similarity(v1, v2)
	}, "Expected panic for vectors with different lengths")
}

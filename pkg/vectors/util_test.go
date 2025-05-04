package vectors

import (
	"math"
	"testing"
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
			if math.Abs(float64(got-tt.want)) > 1e-6 {
				t.Errorf("Similarity() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestSimilarityPanic tests that the Similarity function panics
// when vectors of different lengths are provided
func TestSimilarityPanic(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Expected panic for vectors with different lengths, but no panic occurred")
		}
	}()

	v1 := Vector{1, 2, 3}
	v2 := Vector{1, 2}

	// This should panic
	Similarity(v1, v2)
}

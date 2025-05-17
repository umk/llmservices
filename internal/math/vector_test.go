package math

import (
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestVectorNorm(t *testing.T) {
	tests := []struct {
		name   string
		vector []float32
		want   float64
	}{
		{
			name:   "Unit vector [1, 0, 0]",
			vector: []float32{1, 0, 0},
			want:   1.0,
		},
		{
			name:   "Vector [1, 1, 1]",
			vector: []float32{1, 1, 1},
			want:   math.Sqrt(3),
		},
		{
			name:   "Zero vector",
			vector: []float32{0, 0, 0},
			want:   0.0,
		},
		{
			name:   "Mixed values",
			vector: []float32{3, -4, 5},
			want:   math.Sqrt(50), // 3*3 + (-4)*(-4) + 5*5 = 9 + 16 + 25 = 50
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmp := make([]float32, len(tt.vector))
			got := float64(VectorNorm(tt.vector, tmp))
			assert.InDelta(t, tt.want, got, 1e-6, "Vector norm calculation failed")
		})
	}
}

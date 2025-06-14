package vectors

import (
	"github.com/umk/llmservices/internal/slices"
)

var vectorsPool = slices.NewSlicePool[float32](1 << 12)

func Init(vectorSize int) {
	vectorsPool = slices.NewSlicePool[float32](vectorSize)
}

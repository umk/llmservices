package vectors

import (
	"github.com/umk/llmservices/internal/slices"
)

var vectorsPool = slices.NewSlicePool[float32](20_000)

func Init(vectorSize int) {
	vectorsPool = slices.NewSlicePool[float32](vectorSize)
}

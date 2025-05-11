package vectors

import "github.com/umk/llmservices/internal/slices"

// minDistanceHeap implements a heap interface that compares vector
// records by the order of decreasing their similarity to another vector.
type minDistanceHeap = slices.LimitHeap[*minDistanceHeapItem]

type minDistanceHeapItem struct {
	record     *chunkRecord
	similarity float64
}

func (i *minDistanceHeapItem) Less(another *minDistanceHeapItem) bool {
	return i.similarity < another.similarity
}

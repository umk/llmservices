package vectors

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/umk/llmservices/internal/slices"
)

func TestMinDistanceHeap_Order(t *testing.T) {
	h := slices.MakeLimitHeap[*minDistanceHeapItem](3)

	h.Push(&minDistanceHeapItem{record: &chunkRecord{ID: 1}, similarity: 0.2})
	h.Push(&minDistanceHeapItem{record: &chunkRecord{ID: 2}, similarity: 0.8})
	h.Push(&minDistanceHeapItem{record: &chunkRecord{ID: 3}, similarity: 0.5})

	assert.Len(t, h, 3, "heap should have size 3")

	// The min should be the lowest similarity
	min := h[0]
	for _, item := range h {
		if item.similarity < min.similarity {
			min = item
		}
	}
	assert.Equal(t, 0.2, min.similarity, "expected lowest similarity at root")
}

func TestMinDistanceHeap_Capacity(t *testing.T) {
	h := make(minDistanceHeap, 0, 2)

	h.Push(&minDistanceHeapItem{record: &chunkRecord{ID: 1}, similarity: 0.1})
	h.Push(&minDistanceHeapItem{record: &chunkRecord{ID: 2}, similarity: 0.2})
	h.Push(&minDistanceHeapItem{record: &chunkRecord{ID: 3}, similarity: 0.3}) // Should replace the lowest

	assert.Len(t, h, 2, "heap should have size 2")

	found := map[ID]bool{}
	for _, item := range h {
		found[item.record.ID] = true
	}
	assert.True(t, found[2], "heap should contain ID 2")
	assert.True(t, found[3], "heap should contain ID 3")
}

func TestMinDistanceHeapItem_Less(t *testing.T) {
	a := &minDistanceHeapItem{similarity: 0.7}
	b := &minDistanceHeapItem{similarity: 0.5}
	assert.False(t, a.Less(b), "expected a.Less(b) to be false when a.similarity > b.similarity")
	assert.True(t, b.Less(a), "expected b.Less(a) to be true when b.similarity < a.similarity")
}

package vectors

import (
	"testing"

	"github.com/umk/llmservices/internal/slices"
)

func TestMinDistanceHeap_Order(t *testing.T) {
	h := slices.MakeLimitHeap[*minDistanceHeapItem](3)

	h.Push(&minDistanceHeapItem{record: &chunkRecord{id: 1}, similarity: 0.2})
	h.Push(&minDistanceHeapItem{record: &chunkRecord{id: 2}, similarity: 0.8})
	h.Push(&minDistanceHeapItem{record: &chunkRecord{id: 3}, similarity: 0.5})

	if len(h) != 3 {
		t.Fatalf("expected heap size 3, got %d", len(h))
	}

	// The min should be the lowest similarity
	min := h[0]
	for _, item := range h {
		if item.similarity < min.similarity {
			min = item
		}
	}
	if min.similarity != 0.2 {
		t.Errorf("expected highest similarity at root, got %v", min.similarity)
	}
}

func TestMinDistanceHeap_Capacity(t *testing.T) {
	h := make(minDistanceHeap, 0, 2)

	h.Push(&minDistanceHeapItem{record: &chunkRecord{id: 1}, similarity: 0.1})
	h.Push(&minDistanceHeapItem{record: &chunkRecord{id: 2}, similarity: 0.2})
	h.Push(&minDistanceHeapItem{record: &chunkRecord{id: 3}, similarity: 0.3}) // Should replace the lowest

	if len(h) != 2 {
		t.Fatalf("expected heap size 2, got %d", len(h))
	}

	found := map[ID]bool{}
	for _, item := range h {
		found[item.record.id] = true
	}
	if !found[2] || !found[3] {
		t.Errorf("expected heap to contain IDs 2 and 3, got %+v", found)
	}
}

func TestMinDistanceHeapItem_Less(t *testing.T) {
	a := &minDistanceHeapItem{similarity: 0.7}
	b := &minDistanceHeapItem{similarity: 0.5}
	if a.Less(b) {
		t.Errorf("expected a.Less(b) to be true when a.similarity < b.similarity")
	}
	if !b.Less(a) {
		t.Errorf("expected b.Less(a) to be false when b.similarity > a.similarity")
	}
}

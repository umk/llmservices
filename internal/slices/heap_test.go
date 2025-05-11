// heap_test.go
package slices

import (
	"container/heap"
	"testing"
)

type intItem int

func (a intItem) Less(b intItem) bool {
	return a < b
}

type intMaxItem int

func (a intMaxItem) Less(b intMaxItem) bool {
	return a > b
}

func TestSliceHeap_Basic(t *testing.T) {
	h := &SliceHeap[intItem]{}
	heap.Push(h, intItem(3))
	heap.Push(h, intItem(1))
	heap.Push(h, intItem(2))

	if h.Len() != 3 {
		t.Fatalf("expected len 3, got %d", h.Len())
	}

	min := heap.Pop(h).(intItem)
	if min != 1 {
		t.Errorf("expected min 1, got %d", min)
	}
}

func TestLimitHeap_Basic(t *testing.T) {
	h := MakeLimitHeap[intMaxItem](2)
	h.Push(5)
	h.Push(2)
	h.Push(1) // should replace the root

	if len(h) != 2 {
		t.Fatalf("expected len 2, got %d", len(h))
	}

	// The heap should contain 1 and 2 (since 5 was replaced)
	found := map[int]bool{}
	for _, v := range h {
		found[int(v)] = true
	}
	if !found[1] || !found[2] {
		t.Errorf("expected heap to contain 1 and 2, got %+v", h)
	}
}

func TestLimitHeap_Capacity(t *testing.T) {
	h := MakeLimitHeap[intMaxItem](1)
	h.Push(10)
	h.Push(5)
	if len(h) != 1 {
		t.Fatalf("expected len 1, got %d", len(h))
	}
}

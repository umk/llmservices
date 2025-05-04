package slices

import "container/heap"

type SliceHeap[T SliceHeapItem[T]] []T

type SliceHeapItem[T any] interface {
	Less(another T) bool
}

func (h SliceHeap[T]) Len() int           { return len(h) }
func (h SliceHeap[T]) Less(i, j int) bool { return h[i].Less(h[j]) }
func (h SliceHeap[T]) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }

func (h *SliceHeap[T]) Push(x any) {
	*h = append(*h, x.(T))
}

func (h *SliceHeap[T]) Pop() any {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}

// LimitHeap is a heap that maintains a fixed size. The items in the heap
// must implement the Less method the way that it returns true for the
// greater item (i.e. invert the comparison).
type LimitHeap[T SliceHeapItem[T]] SliceHeap[T]

func MakeLimitHeap[T SliceHeapItem[T]](n int) LimitHeap[T] {
	h := make(LimitHeap[T], 0, n)
	heap.Init((*SliceHeap[T])(&h))
	return h
}

func (h *LimitHeap[T]) Push(item T) {
	n := len(*h)
	if n < cap(*h) {
		heap.Push((*SliceHeap[T])(h), item)
	} else if n > 0 && (*h)[0].Less(item) {
		(*h)[0] = item
		heap.Fix((*SliceHeap[T])(h), 0)
	}
}

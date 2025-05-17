// heap_test.go
package slices

import (
	"container/heap"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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
	h := SliceHeap[intItem]{}
	heap.Push(&h, intItem(3))
	heap.Push(&h, intItem(1))
	heap.Push(&h, intItem(2))

	assert.Len(t, h, 3, "heap should have 3 elements")

	min := heap.Pop(&h).(intItem)
	assert.Equal(t, intItem(1), min, "minimum element should be 1")
}

func TestLimitHeap_Basic(t *testing.T) {
	h := MakeLimitHeap[intMaxItem](2)
	h.Push(5)
	h.Push(2)
	h.Push(1) // should replace the root

	require.Len(t, h, 2, "heap should have 2 elements")

	// The heap should contain 1 and 2 (since 5 was replaced)
	assert.Contains(t, h, intMaxItem(1), "heap should contain 1")
	assert.Contains(t, h, intMaxItem(2), "heap should contain 2")
}

func TestLimitHeap_Capacity(t *testing.T) {
	h := MakeLimitHeap[intMaxItem](1)
	h.Push(10)
	h.Push(5)
	assert.Len(t, h, 1, "heap should maintain capacity limit of 1")
}

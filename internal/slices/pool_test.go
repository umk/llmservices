package slices

import "testing"

func TestSlicePool_GetAndPut(t *testing.T) {
	pool := NewSlicePool[int](5)

	// Get a slice of size 3
	slice := pool.Get(3)
	if len(*slice) != 3 {
		t.Errorf("expected length 3, got %d", len(*slice))
	}
	if cap(*slice) != 5 {
		t.Errorf("expected cap 5, got %d", cap(*slice))
	}

	// Put the slice back
	pool.Put(slice)

	// Get again, should reuse the pooled slice
	slice2 := pool.Get(2)
	if len(*slice2) != 2 {
		t.Errorf("expected length 2, got %d", len(*slice2))
	}
	if cap(*slice2) != 5 {
		t.Errorf("expected cap 5, got %d", cap(*slice2))
	}

	// Get a slice larger than pool size
	bigSlice := pool.Get(10)
	if len(*bigSlice) != 10 {
		t.Errorf("expected length 10, got %d", len(*bigSlice))
	}
	if cap(*bigSlice) != 10 {
		t.Errorf("expected cap 10, got %d", cap(*bigSlice))
	}
}

func TestSlicePool_PutNonPooledSize(t *testing.T) {
	pool := NewSlicePool[int](4)
	s := make([]int, 0, 6)
	pool.Put(&s) // Should not panic or add to pool
}

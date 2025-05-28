package vectors

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	mathutil "github.com/umk/llmservices/internal/math"
)

func TestNewChunk(t *testing.T) {
	baseId := ID(100)
	chunkSize := 10
	chunk := newChunk(baseId, chunkSize)

	assert.Equal(t, baseId, chunk.BaseId, "newChunk baseId doesn't match")
	assert.Empty(t, chunk.Records, "newChunk records should be empty")
	assert.Equal(t, chunkSize, cap(chunk.Records), "newChunk cap(records) doesn't match")
}

func TestVectorsChunk_Add(t *testing.T) {
	baseId := ID(0)
	chunkSize := 3
	chunk := newChunk(baseId, chunkSize)

	vec1 := Vector{1.0, 0.0, 0.0}
	vec2 := Vector{0.0, 1.0, 0.0}
	vec3 := Vector{0.0, 0.0, 1.0}
	vec4 := Vector{1.0, 1.0, 1.0} // This one should fail to add

	// Add first vector
	id1 := chunk.add(vec1)
	assert.Equal(t, ID(0), id1, "add(vec1) returned incorrect id")
	require.Len(t, chunk.Records, 1, "len(records) after first add incorrect")

	rec1 := chunk.Records[0]
	expectedNorm1 := mathutil.VectorNorm(vec1, nil)
	assert.Equal(t, ID(0), rec1.Id, "first record id incorrect")
	assert.True(t, mathutil.VectorsEq(rec1.Vector, vec1), "first record vector incorrect")
	assert.InDelta(t, expectedNorm1, rec1.Norm, 1e-6, "first record norm incorrect")

	// Add second vector
	id2 := chunk.add(vec2)
	assert.Equal(t, ID(1), id2, "add(vec2) returned incorrect id")
	require.Len(t, chunk.Records, 2, "len(records) after second add incorrect")

	rec2 := chunk.Records[1]
	expectedNorm2 := mathutil.VectorNorm(vec2, nil)
	assert.Equal(t, ID(1), rec2.Id, "second record id incorrect")
	assert.True(t, mathutil.VectorsEq(rec2.Vector, vec2), "second record vector incorrect")
	assert.InDelta(t, expectedNorm2, rec2.Norm, 1e-6, "second record norm incorrect")

	// Add third vector
	id3 := chunk.add(vec3)
	assert.Equal(t, ID(2), id3, "add(vec3) returned incorrect id")
	require.Len(t, chunk.Records, 3, "len(records) after third add incorrect")

	rec3 := chunk.Records[2]
	expectedNorm3 := mathutil.VectorNorm(vec3, nil)
	assert.Equal(t, ID(2), rec3.Id, "third record id incorrect")
	assert.True(t, mathutil.VectorsEq(rec3.Vector, vec3), "third record vector incorrect")
	assert.InDelta(t, expectedNorm3, rec3.Norm, 1e-6, "third record norm incorrect")

	// Try adding when full
	id4 := chunk.add(vec4)
	assert.Equal(t, ID(-1), id4, "add(vec4) when full returned incorrect id")
	assert.Len(t, chunk.Records, 3, "len(records) after failed add incorrect")
}

func TestVectorsChunk_AddWithBaseId(t *testing.T) {
	baseId := ID(100)
	chunkSize := 2
	chunk := newChunk(baseId, chunkSize)

	vec1 := Vector{1.0, 2.0}
	vec2 := Vector{3.0, 4.0}

	id1 := chunk.add(vec1)
	assert.Equal(t, ID(100), id1, "add(vec1) returned incorrect id")

	id2 := chunk.add(vec2)
	assert.Equal(t, ID(101), id2, "add(vec2) returned incorrect id")

	require.Len(t, chunk.Records, 2, "len(records) incorrect")
	assert.Equal(t, ID(100), chunk.Records[0].Id, "first record id incorrect")
	assert.Equal(t, ID(101), chunk.Records[1].Id, "second record id incorrect")
}

func TestVectorsChunk_Delete(t *testing.T) {
	baseId := ID(10)
	chunkSize := 5
	chunk := newChunk(baseId, chunkSize)

	// Add some vectors
	ids := make([]ID, chunkSize)
	for i := range chunkSize {
		vec := Vector{float32(i)}
		ids[i] = chunk.add(vec)
		require.Equal(t, baseId+ID(i), ids[i], "Failed to add vector %d", i)
	}

	// Test deleting existing items
	testCases := []struct {
		name       string
		idToDelete ID
		expectOk   bool
		expectNil  bool // whether the record should be nil after deletion
	}{
		{name: "Delete middle", idToDelete: 12, expectOk: true, expectNil: true},
		{name: "Delete first", idToDelete: 10, expectOk: true, expectNil: true},
		{name: "Delete last", idToDelete: 14, expectOk: true, expectNil: true},
		{name: "Delete already deleted", idToDelete: 12, expectOk: false, expectNil: true},
		{name: "Delete non-existent (too low)", idToDelete: 9, expectOk: false, expectNil: false},
		{name: "Delete non-existent (too high)", idToDelete: 15, expectOk: false, expectNil: false},
		{name: "Delete another existing", idToDelete: 11, expectOk: true, expectNil: true},
		{name: "Delete the last remaining", idToDelete: 13, expectOk: true, expectNil: true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ok := chunk.delete(tc.idToDelete)
			assert.Equal(t, tc.expectOk, ok, "delete(%v) returned unexpected result", tc.idToDelete)

			// Verify record state (only if ID is within the initial range of this chunk)
			if tc.idToDelete >= baseId && tc.idToDelete < baseId+ID(chunkSize) {
				index := int(tc.idToDelete - baseId)
				if index < 0 || index >= len(chunk.Records) {
					if tc.expectNil {
						assert.Fail(t, "Index out of bounds for id, but expected record to be nil",
							"index=%d, id=%v", index, tc.idToDelete)
					}
					return // Cannot check chunk.records[index]
				}
				recordIsNil := chunk.Records[index].Vector == nil
				assert.Equal(t, tc.expectNil, recordIsNil,
					"record at index %d (id %v) has unexpected nil state", index, tc.idToDelete)
			}
		})
	}

	// Final check: ensure only nil records remain for deleted IDs 10, 11, 12, 13, 14
	expectedNils := map[ID]bool{10: true, 11: true, 12: true, 13: true, 14: true}
	for i, rec := range chunk.Records {
		id := baseId + ID(i)
		isNil := rec.Vector == nil
		shouldBeNil := expectedNils[id]
		assert.Equal(t, shouldBeNil, isNil, "Final state check: record for id %v has unexpected nil state", id)
	}
}

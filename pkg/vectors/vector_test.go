package vectors

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewVectors(t *testing.T) {
	chunkSize := 128
	v := NewVectors(chunkSize)

	require.NotNil(t, v, "NewVectors returned nil")
	assert.Equal(t, chunkSize, v.header.ChunkSize, "incorrect chunkSize")
	assert.Len(t, v.chunks, 1, "expected 1 chunk initially")
	require.NotNil(t, v.currentChunk, "currentChunk is nil")
	assert.Equal(t, ID(0), v.currentChunk.BaseId, "incorrect baseId")
	assert.Equal(t, chunkSize, cap(v.currentChunk.Records), "incorrect records capacity")
}

func TestVectors_Add(t *testing.T) {
	chunkSize := 2 // Small chunk size for testing chunk creation
	v := NewVectors(chunkSize)

	// Add first vector
	vec1 := Vector{1.0, 2.0, 3.0}
	id1 := v.Add(vec1)
	assert.Equal(t, ID(0), id1, "expected first ID to be 0")

	// Add second vector (still fits in first chunk)
	vec2 := Vector{4.0, 5.0, 6.0}
	id2 := v.Add(vec2)
	assert.Equal(t, ID(1), id2, "expected second ID to be 1")

	// Add third vector (should create a new chunk)
	vec3 := Vector{7.0, 8.0, 9.0}
	id3 := v.Add(vec3)
	assert.Equal(t, ID(2), id3, "expected third ID to be 2")

	// Check that we now have 2 chunks
	assert.Len(t, v.chunks, 2, "expected 2 chunks after adding 3 vectors")

	// First chunk should be full
	assert.Len(t, v.chunks[0].Records, chunkSize, "unexpected number of records in first chunk")

	// Second chunk should have 1 record
	assert.Len(t, v.chunks[1].Records, 1, "expected second chunk to have 1 record")

	// Check that the current chunk is the second chunk
	assert.Equal(t, v.chunks[1], v.currentChunk, "currentChunk should be the second chunk")
}

func TestVectors_Delete(t *testing.T) {
	v := NewVectors(5)

	// Add some vectors
	vec1 := Vector{1.0, 2.0, 3.0}
	vec2 := Vector{4.0, 5.0, 6.0}
	vec3 := Vector{7.0, 8.0, 9.0}

	id1 := v.Add(vec1)
	id2 := v.Add(vec2)
	id3 := v.Add(vec3)

	// Check that the IDs are sequential
	assert.Equal(t, ID(0), id1, "Expected ID 0")
	assert.Equal(t, ID(1), id2, "Expected ID 1")
	assert.Equal(t, ID(2), id3, "Expected ID 2")

	// Test deleting an existing vector
	assert.True(t, v.Delete(id2), "Delete(%d) should return true", id2)

	// Try to delete the same vector again (should fail)
	assert.False(t, v.Delete(id2), "Delete(%d) should return false when deleting already deleted vector", id2)

	// Delete non-existent vector
	assert.False(t, v.Delete(ID(99)), "Delete(99) should return false for non-existent vector")
}

func TestVectors_Get(t *testing.T) {
	v := NewVectors(10)

	// Add vectors
	id0 := v.Add(Vector{1, 0, 0})     // ID 0
	id1 := v.Add(Vector{0, 1, 0})     // ID 1
	id2 := v.Add(Vector{0, 0, 1})     // ID 2
	id3 := v.Add(Vector{0.7, 0.7, 0}) // ID 3 - closer to [1,1,0] than the others

	// Search for vectors similar to [1, 1, 0]
	query := Vector{1, 1, 0}
	results := v.Get([]Vector{query}, 4) // Get all 4 vectors

	assert.Len(t, results, 4, "expected 4 results")

	// Convert results to a map to check containment regardless of order
	resultMap := make(map[ID]bool)
	for _, id := range results {
		resultMap[id] = true
	}

	// Verify all vectors are returned
	assert.True(t, resultMap[id0] && resultMap[id1] && resultMap[id2] && resultMap[id3],
		"Not all expected vectors were returned: %v", results)
}

func TestVectors_Compact(t *testing.T) {
	v := NewVectors(3)

	// Add vectors with distinct values to make them easily identifiable
	ids := make([]ID, 4)
	vecs := []Vector{
		{1.0, 0.0, 0.0}, // Vector pointing in x direction
		{0.0, 1.0, 0.0}, // Vector pointing in y direction
		{0.0, 0.0, 1.0}, // Vector pointing in z direction
		{1.0, 1.0, 1.0}, // Vector pointing along the diagonal
	}

	for i, vec := range vecs {
		ids[i] = v.Add(vec)
	}

	// Delete one vector
	v.Delete(ids[1]) // Delete the y-direction vector

	// Perform compact
	v.Compact()

	assert.Len(t, v.chunks, 1, "expected 1 chunk after compact")
	assert.Equal(t, 3, v.header.ChunkSize, "expected compacted Vectors to maintain original chunkSize")

	// Create queries using the original vectors we added (excluding the deleted one)
	remainingVecs := []struct {
		id    ID
		query Vector
	}{
		{ids[0], vecs[0]}, // x direction
		{ids[2], vecs[2]}, // z direction
		{ids[3], vecs[3]}, // diagonal
	}

	for _, rv := range remainingVecs {
		results := v.Get([]Vector{rv.query}, 1)
		assert.Len(t, results, 1, "Expected 1 result for vector %v", rv.query)
		if len(results) > 0 {
			assert.Equal(t, rv.id, results[0], "Expected ID %d for vector %v", rv.id, rv.query)
		}
	}
}

func TestVectors_Repack(t *testing.T) {
	v := NewVectors(3)

	// Add vectors with distinct values to make them easily identifiable
	ids := make([]ID, 4)
	vecs := []Vector{
		{1.0, 0.0, 0.0}, // Vector pointing in x direction
		{0.0, 1.0, 0.0}, // Vector pointing in y direction
		{0.0, 0.0, 1.0}, // Vector pointing in z direction
		{1.0, 1.0, 1.0}, // Vector pointing along the diagonal
	}

	for i, vec := range vecs {
		ids[i] = v.Add(vec)
	}

	// Delete one vector
	v.Delete(ids[1]) // Delete the y-direction vector

	// Perform repack
	newVectors := v.Repack()

	assert.Len(t, newVectors.chunks, 1, "expected 1 chunk after repack")
	assert.Equal(t, v.header.ChunkSize, newVectors.header.ChunkSize, "expected repacked Vectors to have same chunkSize")

	// Create queries using the original vectors we added (excluding the deleted one)
	remainingVecs := []struct {
		id    ID
		query Vector
	}{
		{ids[0], vecs[0]}, // x direction
		{ids[2], vecs[2]}, // z direction
		{ids[3], vecs[3]}, // diagonal
	}

	for _, rv := range remainingVecs {
		results := newVectors.Get([]Vector{rv.query}, 1)
		assert.Len(t, results, 1, "Expected 1 result for vector %v", rv.query)
		if len(results) > 0 {
			assert.Equal(t, rv.id, results[0], "Expected ID %d for vector %v", rv.id, rv.query)
		}
	}
}

func TestVectors_DeleteAndQuery(t *testing.T) {
	v := NewVectors(5)

	// Add vectors with distinct values
	vec1 := Vector{1.0, 0.0, 0.0} // x-axis
	vec2 := Vector{0.0, 1.0, 0.0} // y-axis
	vec3 := Vector{0.0, 0.0, 1.0} // z-axis

	id1 := v.Add(vec1)
	id2 := v.Add(vec2)
	id3 := v.Add(vec3)

	// Initial verification - query should return all vectors
	query := Vector{1.0, 1.0, 1.0} // diagonal vector equally distant from all test vectors
	results := v.Get([]Vector{query}, 3)

	assert.Len(t, results, 3, "Expected 3 vectors initially")

	// Check that all IDs are present
	resultMap := make(map[ID]bool)
	for _, id := range results {
		resultMap[id] = true
	}
	assert.True(t, resultMap[id1] && resultMap[id2] && resultMap[id3],
		"Not all expected vectors were returned before deletion: %v", results)

	// Delete the second vector
	assert.True(t, v.Delete(id2), "Failed to delete vector with ID %d", id2)

	// Query again
	results = v.Get([]Vector{query}, 3)

	// Should have only 2 results now
	assert.Len(t, results, 2, "Expected 2 vectors after deletion")

	// The results should contain id1 and id3, but not id2
	resultMap = make(map[ID]bool)
	for _, id := range results {
		resultMap[id] = true
	}

	assert.True(t, resultMap[id1] && resultMap[id3],
		"Missing expected vectors after deletion: %v", results)
	assert.False(t, resultMap[id2],
		"Deleted vector with ID %d was still returned in query", id2)
}

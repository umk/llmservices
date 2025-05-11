package vectors

import (
	"testing"
)

func TestNewVectors(t *testing.T) {
	chunkSize := 128
	v := NewVectors(chunkSize)

	if v == nil {
		t.Fatal("NewVectors returned nil")
	}

	if v.chunkSize != chunkSize {
		t.Errorf("expected chunkSize %d, got %d", chunkSize, v.chunkSize)
	}

	if len(v.chunks) != 1 {
		t.Errorf("expected 1 chunk initially, got %d", len(v.chunks))
	}

	if v.currentChunk == nil {
		t.Error("currentChunk is nil")
	}

	if v.currentChunk.baseId != 0 {
		t.Errorf("expected baseId 0, got %d", v.currentChunk.baseId)
	}

	if cap(v.currentChunk.records) != chunkSize {
		t.Errorf("expected records capacity %d, got %d", chunkSize, cap(v.currentChunk.records))
	}
}

func TestVectors_Add(t *testing.T) {
	chunkSize := 2 // Small chunk size for testing chunk creation
	v := NewVectors(chunkSize)

	// Add first vector
	vec1 := Vector{1.0, 2.0, 3.0}
	id1 := v.Add(vec1)
	if id1 != 0 {
		t.Errorf("expected first ID to be 0, got %d", id1)
	}

	// Add second vector (still fits in first chunk)
	vec2 := Vector{4.0, 5.0, 6.0}
	id2 := v.Add(vec2)
	if id2 != 1 {
		t.Errorf("expected second ID to be 1, got %d", id2)
	}

	// Add third vector (should create a new chunk)
	vec3 := Vector{7.0, 8.0, 9.0}
	id3 := v.Add(vec3)
	if id3 != 2 {
		t.Errorf("expected third ID to be 2, got %d", id3)
	}

	// Check that we now have 2 chunks
	if len(v.chunks) != 2 {
		t.Errorf("expected 2 chunks after adding 3 vectors, got %d", len(v.chunks))
	}

	// First chunk should be full
	if len(v.chunks[0].records) != chunkSize {
		t.Errorf("expected first chunk to have %d records, got %d", chunkSize, len(v.chunks[0].records))
	}

	// Second chunk should have 1 record
	if len(v.chunks[1].records) != 1 {
		t.Errorf("expected second chunk to have 1 record, got %d", len(v.chunks[1].records))
	}

	// Check that the current chunk is the second chunk
	if v.currentChunk != v.chunks[1] {
		t.Error("expected currentChunk to be the second chunk")
	}
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
	if id1 != 0 || id2 != 1 || id3 != 2 {
		t.Errorf("Expected sequential IDs 0,1,2, got %d,%d,%d", id1, id2, id3)
	}

	// Test deleting an existing vector
	if !v.Delete(id2) {
		t.Errorf("Delete(%d) returned false, expected true", id2)
	}

	// Try to delete the same vector again (should fail)
	if v.Delete(id2) {
		t.Errorf("Delete(%d) returned true when trying to delete already deleted vector", id2)
	}

	// Delete non-existent vector
	if v.Delete(ID(99)) {
		t.Errorf("Delete(99) returned true for non-existent vector")
	}
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

	if len(results) != 4 {
		t.Fatalf("expected 4 results, got %d", len(results))
	}

	// Since similarity is cosine similarity, we can calculate expected results:
	// Similarity between [1,1,0] and [0.7,0.7,0] should be highest
	// Then [1,0,0] and [0,1,0] should be equal
	// [0,0,1] should be lowest (orthogonal)

	// Convert results to a map to check containment regardless of order
	resultMap := make(map[ID]bool)
	for _, id := range results {
		resultMap[id] = true
	}

	// Verify all vectors are returned
	if !resultMap[id0] || !resultMap[id1] || !resultMap[id2] || !resultMap[id3] {
		t.Errorf("Not all expected vectors were returned: %v", results)
	}
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

	if len(v.chunks) != 1 {
		t.Errorf("expected 1 chunk after compact, got %d", len(v.chunks))
	}

	// Check that the new Vectors has the same chunkSize
	if v.chunkSize != 3 {
		t.Errorf("expected compacted Vectors to have chunkSize %d, got %d", v.chunkSize, v.chunkSize)
	}

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
		if len(results) != 1 {
			t.Fatalf("Expected 1 result for vector %v, got %d results", rv.query, len(results))
			continue
		}

		if results[0] != rv.id {
			t.Fatalf("Expected ID %d for vector %v, got %d", rv.id, rv.query, results[0])
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

	if len(newVectors.chunks) != 1 {
		t.Errorf("expected 1 chunk after repack, got %d", len(newVectors.chunks))
	}

	// Check that the new Vectors has the same chunkSize
	if newVectors.chunkSize != v.chunkSize {
		t.Errorf("expected repacked Vectors to have chunkSize %d, got %d", v.chunkSize, newVectors.chunkSize)
	}

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
		if len(results) != 1 {
			t.Fatalf("Expected 1 result for vector %v, got %d results", rv.query, len(results))
			continue
		}

		if results[0] != rv.id {
			t.Fatalf("Expected ID %d for vector %v, got %d", rv.id, rv.query, results[0])
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

	if len(results) != 3 {
		t.Fatalf("Expected 3 vectors initially, got %d", len(results))
	}

	// Check that all IDs are present
	resultMap := make(map[ID]bool)
	for _, id := range results {
		resultMap[id] = true
	}
	if !resultMap[id1] || !resultMap[id2] || !resultMap[id3] {
		t.Errorf("Not all expected vectors were returned before deletion: %v", results)
	}

	// Delete the second vector
	if !v.Delete(id2) {
		t.Fatalf("Failed to delete vector with ID %d", id2)
	}

	// Query again
	results = v.Get([]Vector{query}, 3)

	// Should have only 2 results now
	if len(results) != 2 {
		t.Fatalf("Expected 2 vectors after deletion, got %d", len(results))
	}

	// The results should contain id1 and id3, but not id2
	resultMap = make(map[ID]bool)
	for _, id := range results {
		resultMap[id] = true
	}

	if !resultMap[id1] || !resultMap[id3] {
		t.Errorf("Missing expected vectors after deletion: %v", results)
	}

	if resultMap[id2] {
		t.Errorf("Deleted vector with ID %d was still returned in query", id2)
	}
}

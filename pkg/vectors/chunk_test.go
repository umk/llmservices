package vectors

import (
	"testing"

	mathutil "github.com/umk/llmservices/internal/math"
)

func TestNewChunk(t *testing.T) {
	baseId := ID(100)
	chunkSize := 10
	chunk := newChunk(baseId, chunkSize)

	if chunk.baseId != baseId {
		t.Errorf("newChunk baseId = %v, want %v", chunk.baseId, baseId)
	}
	if len(chunk.records) != 0 {
		t.Errorf("newChunk len(records) = %d, want 0", len(chunk.records))
	}
	if cap(chunk.records) != chunkSize {
		t.Errorf("newChunk cap(records) = %d, want %d", cap(chunk.records), chunkSize)
	}
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
	if id1 != 0 {
		t.Errorf("add(vec1) id = %v, want 0", id1)
	}
	if len(chunk.records) != 1 {
		t.Fatalf("len(records) after first add = %d, want 1", len(chunk.records))
	}
	rec1 := chunk.records[0]
	// Calculate expected norm for vec1
	expectedNorm1 := mathutil.VectorNorm(vec1, nil)
	if rec1.id != 0 || !mathutil.VectorsEq(rec1.vector, vec1) || !mathutil.ApproxEq(rec1.norm, expectedNorm1, 1e-6) {
		t.Errorf("first record = %+v, want {id: 0, vector: %v, norm: %f}", rec1, vec1, expectedNorm1)
	}

	// Add second vector
	id2 := chunk.add(vec2)
	if id2 != 1 {
		t.Errorf("add(vec2) id = %v, want 1", id2)
	}
	if len(chunk.records) != 2 {
		t.Fatalf("len(records) after second add = %d, want 2", len(chunk.records))
	}
	rec2 := chunk.records[1]
	// Calculate expected norm for vec2
	expectedNorm2 := mathutil.VectorNorm(vec2, nil)
	if rec2.id != 1 || !mathutil.VectorsEq(rec2.vector, vec2) || !mathutil.ApproxEq(rec2.norm, expectedNorm2, 1e-6) {
		t.Errorf("second record = %+v, want {id: 1, vector: %v, norm: %f}", rec2, vec2, expectedNorm2)
	}

	// Add third vector
	id3 := chunk.add(vec3)
	if id3 != 2 {
		t.Errorf("add(vec3) id = %v, want 2", id3)
	}
	if len(chunk.records) != 3 {
		t.Fatalf("len(records) after third add = %d, want 3", len(chunk.records))
	}
	rec3 := chunk.records[2]
	// Calculate expected norm for vec3
	expectedNorm3 := mathutil.VectorNorm(vec3, nil)
	if rec3.id != 2 || !mathutil.VectorsEq(rec3.vector, vec3) || !mathutil.ApproxEq(rec3.norm, expectedNorm3, 1e-6) {
		t.Errorf("third record = %+v, want {id: 2, vector: %v, norm: %f}", rec3, vec3, expectedNorm3)
	}

	// Try adding when full
	id4 := chunk.add(vec4)
	if id4 != -1 {
		t.Errorf("add(vec4) when full id = %v, want -1", id4)
	}
	if len(chunk.records) != 3 {
		t.Errorf("len(records) after failed add = %d, want 3", len(chunk.records))
	}
}

func TestVectorsChunk_AddWithBaseId(t *testing.T) {
	baseId := ID(100)
	chunkSize := 2
	chunk := newChunk(baseId, chunkSize)

	vec1 := Vector{1.0, 2.0}
	vec2 := Vector{3.0, 4.0}

	id1 := chunk.add(vec1)
	if id1 != 100 {
		t.Errorf("add(vec1) id = %v, want 100", id1)
	}
	id2 := chunk.add(vec2)
	if id2 != 101 {
		t.Errorf("add(vec2) id = %v, want 101", id2)
	}

	if len(chunk.records) != 2 {
		t.Fatalf("len(records) = %d, want 2", len(chunk.records))
	}
	if chunk.records[0].id != 100 || chunk.records[1].id != 101 {
		t.Errorf("Expected record IDs 100 and 101, got %v and %v", chunk.records[0].id, chunk.records[1].id)
	}
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
		if ids[i] != baseId+ID(i) {
			t.Fatalf("Failed to add vector %d, expected id %v, got %v", i, baseId+ID(i), ids[i])
		}
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
		{name: "Delete non-existent (too low)", idToDelete: 9, expectOk: false, expectNil: false},   // expectNil is false as it's out of bounds
		{name: "Delete non-existent (too high)", idToDelete: 15, expectOk: false, expectNil: false}, // expectNil is false as it's out of bounds
		{name: "Delete another existing", idToDelete: 11, expectOk: true, expectNil: true},
		{name: "Delete the last remaining", idToDelete: 13, expectOk: true, expectNil: true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ok := chunk.delete(tc.idToDelete)
			if ok != tc.expectOk {
				t.Errorf("delete(%v) returned %v, want %v", tc.idToDelete, ok, tc.expectOk)
			}

			// Verify record state (only if ID is within the initial range of this chunk)
			if tc.idToDelete >= baseId && tc.idToDelete < baseId+ID(chunkSize) {
				index := int(tc.idToDelete - baseId)
				if index < 0 || index >= len(chunk.records) {
					// This case should ideally not happen if the ID is within range,
					// but adding a safeguard.
					if tc.expectNil {
						t.Errorf("Index %d out of bounds for id %v, but expected record to be nil", index, tc.idToDelete)
					}
					return // Cannot check chunk.records[index]
				}
				recordIsNil := chunk.records[index].vector == nil
				if recordIsNil != tc.expectNil {
					t.Errorf("record at index %d (id %v) nil state = %v, want %v", index, tc.idToDelete, recordIsNil, tc.expectNil)
				}
			}
		})
	}

	// Final check: ensure only nil records remain for deleted IDs 10, 11, 12, 13, 14
	expectedNils := map[ID]bool{10: true, 11: true, 12: true, 13: true, 14: true}
	for i, rec := range chunk.records {
		id := baseId + ID(i)
		isNil := rec.vector == nil
		shouldBeNil := expectedNils[id]
		if isNil != shouldBeNil {
			t.Errorf("Final state check: record for id %v nil state = %v, want %v", id, isNil, shouldBeNil)
		}
	}
}

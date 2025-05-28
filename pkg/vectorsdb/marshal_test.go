package vectorsdb

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/umk/llmservices/pkg/vectors"
)

// TestMarshalUnmarshalSimple tests basic marshaling and unmarshaling functionality
func TestMarshalUnmarshalSimple(t *testing.T) {
	// Create a test database with string values
	db := NewDatabase[string](3)

	// Add some test data
	record1 := Record[string]{
		Vector: vectors.Vector{0.1, 0.2, 0.3},
		Data:   "test data 1",
	}
	record2 := Record[string]{
		Vector: vectors.Vector{0.4, 0.5, 0.6},
		Data:   "test data 2",
	}

	addedRecord1 := db.Add(record1)
	addedRecord2 := db.Add(record2)

	// Marshal the database
	buf := &bytes.Buffer{}
	err := Marshal(buf, db)
	require.NoError(t, err)

	fmt.Println(buf.Len())

	// Unmarshal the database
	unmarshalled, err := Unmarshal[string](buf)
	require.NoError(t, err)

	// Verify the unmarshalled database header
	assert.Equal(t, db.header.VectorLength, unmarshalled.header.VectorLength)
	assert.Equal(t, db.header.RepackPercent, unmarshalled.header.RepackPercent)
	assert.Equal(t, db.header.ItemsCount, unmarshalled.header.ItemsCount)
	assert.Equal(t, db.header.DeletesCount, unmarshalled.header.DeletesCount)

	// Verify data was preserved
	assert.Equal(t, len(db.Data), len(unmarshalled.Data))
	assert.Equal(t, db.Data[addedRecord1.Id], unmarshalled.Data[addedRecord1.Id])
	assert.Equal(t, db.Data[addedRecord2.Id], unmarshalled.Data[addedRecord2.Id])

	// Test search functionality to verify vectors were preserved correctly
	queryVec := vectors.Vector{0.1, 0.2, 0.3} // Should match record1
	results := unmarshalled.Get([]vectors.Vector{queryVec}, 1)
	require.Len(t, results, 1)
	assert.Equal(t, "test data 1", results[0].Data)
}

// TestMarshalUnmarshalEmpty tests marshaling and unmarshaling an empty database
func TestMarshalUnmarshalEmpty(t *testing.T) {
	// Create an empty database
	db := NewDatabase[string](3)

	// Marshal the database
	buf := &bytes.Buffer{}
	err := Marshal(buf, db)
	require.NoError(t, err)

	// Unmarshal the database
	unmarshalled, err := Unmarshal[string](buf)
	require.NoError(t, err)

	// Verify the unmarshalled database
	assert.Equal(t, db.header.VectorLength, unmarshalled.header.VectorLength)
	assert.Equal(t, db.header.RepackPercent, unmarshalled.header.RepackPercent)
	assert.Equal(t, db.header.ItemsCount, unmarshalled.header.ItemsCount)
	assert.Equal(t, db.header.DeletesCount, unmarshalled.header.DeletesCount)
	assert.Empty(t, unmarshalled.Data)
}

// TestMarshalUnmarshalWithDeletes tests marshaling and unmarshaling with deleted records
func TestMarshalUnmarshalWithDeletes(t *testing.T) {
	// Create a test database
	db := NewDatabase[string](3)

	// Add some test data
	record1 := Record[string]{
		Vector: vectors.Vector{0.1, 0.2, 0.3},
		Data:   "test data 1",
	}
	record2 := Record[string]{
		Vector: vectors.Vector{0.4, 0.5, 0.6},
		Data:   "test data 2",
	}
	record3 := Record[string]{
		Vector: vectors.Vector{0.7, 0.8, 0.9},
		Data:   "test data 3",
	}

	addedRecord1 := db.Add(record1)
	addedRecord2 := db.Add(record2)
	addedRecord3 := db.Add(record3)

	// Delete one record
	db.Delete(addedRecord2.Id)

	// Marshal the database
	buf := &bytes.Buffer{}
	err := Marshal(buf, db)
	require.NoError(t, err)

	// Unmarshal the database
	unmarshalled, err := Unmarshal[string](buf)
	require.NoError(t, err)

	// Verify data was preserved (including the deletion)
	assert.Equal(t, len(db.Data), len(unmarshalled.Data))
	assert.Equal(t, db.Data[addedRecord1.Id], unmarshalled.Data[addedRecord1.Id])
	assert.Equal(t, db.Data[addedRecord3.Id], unmarshalled.Data[addedRecord3.Id])
	_, ok := unmarshalled.Data[addedRecord2.Id]
	assert.False(t, ok, "Deleted record should not be present")
}

// TestStruct is a complex data type for testing
type TestStruct struct {
	Name  string
	Value int
	Tags  []string
}

// TestMarshalUnmarshalComplexType tests marshaling with a more complex data type
func TestMarshalUnmarshalComplexType(t *testing.T) {
	// Create a test database with a complex type
	db := NewDatabase[TestStruct](3)

	// Add some test data
	record1 := Record[TestStruct]{
		Vector: vectors.Vector{0.1, 0.2, 0.3},
		Data:   TestStruct{Name: "item1", Value: 42, Tags: []string{"tag1", "tag2"}},
	}

	addedRecord := db.Add(record1)

	// Marshal the database
	buf := &bytes.Buffer{}
	err := Marshal(buf, db)
	require.NoError(t, err)

	// Unmarshal the database
	unmarshalled, err := Unmarshal[TestStruct](buf)
	require.NoError(t, err)

	// Verify data was preserved
	assert.Equal(t, db.Data[addedRecord.Id], unmarshalled.Data[addedRecord.Id])
	assert.Equal(t, "item1", unmarshalled.Data[addedRecord.Id].Name)
	assert.Equal(t, 42, unmarshalled.Data[addedRecord.Id].Value)
	assert.Equal(t, []string{"tag1", "tag2"}, unmarshalled.Data[addedRecord.Id].Tags)
}

// Error testing helpers
type errorWriter struct{}

func (w *errorWriter) Write(p []byte) (int, error) {
	return 0, errors.New("test write error")
}

type errorReader struct{}

func (r *errorReader) Read(p []byte) (int, error) {
	return 0, errors.New("test read error")
}

// TestMarshalError tests error handling in Marshal
func TestMarshalError(t *testing.T) {
	db := NewDatabase[string](3)
	err := Marshal(&errorWriter{}, db)
	assert.Error(t, err)
}

// TestUnmarshalError tests error handling in Unmarshal
func TestUnmarshalError(t *testing.T) {
	_, err := Unmarshal[string](&errorReader{})
	assert.Error(t, err)
}

// limitedReader returns data for a few reads, then starts returning errors
type limitedReader struct {
	data     []byte
	pos      int
	maxReads int
	reads    int
}

func newLimitedReader(data []byte, maxReads int) *limitedReader {
	return &limitedReader{
		data:     data,
		maxReads: maxReads,
	}
}

func (r *limitedReader) Read(p []byte) (int, error) {
	if r.reads >= r.maxReads {
		return 0, errors.New("max reads exceeded")
	}
	r.reads++

	if r.pos >= len(r.data) {
		return 0, io.EOF
	}

	n := copy(p, r.data[r.pos:])
	r.pos += n
	return n, nil
}

// TestUnmarshalPartialReads tests unmarshaling with partial reads
func TestUnmarshalPartialReads(t *testing.T) {
	// Create and marshal a valid database
	db := NewDatabase[string](3)
	db.Add(Record[string]{
		Vector: vectors.Vector{0.1, 0.2, 0.3},
		Data:   "test",
	})

	buf := &bytes.Buffer{}
	err := Marshal(buf, db)
	require.NoError(t, err)

	data := buf.Bytes()

	// Test with limited reader that fails after a few reads
	// Testing different failure points
	for _, maxReads := range []int{1, 2, 5} {
		_, err = Unmarshal[string](newLimitedReader(data, maxReads))
		assert.Error(t, err)
	}
}

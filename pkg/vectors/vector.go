package vectors

import (
	"slices"
	"sync"

	"github.com/umk/llmservices/internal/math"
	slicesutil "github.com/umk/llmservices/internal/slices"
)

type ID int64

type Vector []float32

type Vectors struct {
	chunkSize int

	chunks       []*vectorsChunk
	currentChunk *vectorsChunk
}

func NewVectors(chunkSize int) *Vectors {
	vectors := &Vectors{
		chunkSize: chunkSize,
		chunks:    make([]*vectorsChunk, 1, 32),
	}

	currentChunk := newChunk(ID(0), chunkSize)

	vectors.chunks[0] = currentChunk
	vectors.currentChunk = currentChunk

	return vectors
}

func (v *Vectors) Add(vector Vector) ID {
	id := v.currentChunk.add(vector)
	if id >= 0 {
		return id
	}

	baseId := ID(len(v.chunks) * v.chunkSize)

	v.currentChunk = newChunk(baseId, v.chunkSize)
	v.chunks = append(v.chunks, v.currentChunk)

	return v.currentChunk.add(vector)
}

func (v *Vectors) Delete(id ID) bool {
	i, _ := slices.BinarySearchFunc(v.chunks, id, func(c *vectorsChunk, id ID) int {
		return int((c.baseId - 1) - id)
	})
	if i == 0 {
		return false
	}
	return v.chunks[i-1].delete(id)
}

func (v *Vectors) Get(vectors []Vector, n int) []ID {
	h := v.getHeaps(vectors, n)
	r := reduceHeaps(h, n)

	ids := make([]ID, len(r))
	for i, hr := range r {
		ids[i] = hr.record.id
	}

	return ids
}

func (v *Vectors) Compact() {
	var destIndex, destRecordIndex int
	destChunk := v.chunks[destIndex]

	// Compact records by iterating over all chunks and their records
	for _, srcChunk := range v.chunks {
		for _, record := range srcChunk.records {
			if record.vector == nil {
				continue
			}

			// When destination chunk is full, move to next and reset index
			if destRecordIndex == cap(destChunk.records) {
				destIndex++
				destChunk = v.chunks[destIndex]
				destRecordIndex = 0
			}

			// Write the valid record to the destination
			destChunk.records[destRecordIndex] = record
			// Set new base ID for the chunk when writing its first record
			if destRecordIndex == 0 {
				destChunk.baseId = record.id
			}
			destRecordIndex++
		}
	}

	v.currentChunk = destChunk

	// Nil out chunks that are no longer used
	for i := destIndex + 1; i < len(v.chunks); i++ {
		v.chunks[i] = nil
	}
	v.chunks = v.chunks[:destIndex+1]

	// Clear trailing nil records from the destination chunk
	for i := destRecordIndex; i < len(destChunk.records); i++ {
		destChunk.records[i] = nil
	}
	destChunk.records = destChunk.records[:destRecordIndex]
}

func (v *Vectors) Repack() *Vectors {
	// Create a new Vectors instance to hold compacted records
	vectors := &Vectors{
		chunkSize: v.chunkSize,
		chunks:    make([]*vectorsChunk, 1, 32),
	}

	// Initialize the first destination chunk
	destChunk := &vectorsChunk{
		records: make([]*chunkRecord, 0, v.chunkSize),
	}
	vectors.chunks[0] = destChunk

	// Iterate over all existing chunks
	for _, srcChunk := range v.chunks {
		// Iterate over all records in the current source chunk
		for _, record := range srcChunk.records {
			if record.vector == nil {
				continue
			}

			// If the current destination chunk is full, allocate a new one
			if len(destChunk.records) == cap(destChunk.records) {
				destChunk = &vectorsChunk{
					records: make([]*chunkRecord, 0, v.chunkSize),
				}
				vectors.chunks = append(vectors.chunks, destChunk)
			}
			// Set the base ID for a chunk when inserting its first record
			if len(destChunk.records) == 0 {
				destChunk.baseId = record.id
			}

			// Append the valid record
			destChunk.records = append(destChunk.records, record)
		}
	}

	// Update the current chunk pointer
	vectors.currentChunk = destChunk

	return vectors
}

func (v *Vectors) getHeaps(vectors []Vector, n int) <-chan minDistanceHeap {
	out := make(chan minDistanceHeap)

	go func() {
		defer close(out)

		var wg sync.WaitGroup

		for _, vector := range vectors {
			tmp := vectorsPool.Get(len(vector))
			norm := math.VectorNorm(vector, *tmp)
			vectorsPool.Put(tmp)

			for i := range v.chunks {
				wg.Add(1)
				go func(chunk *vectorsChunk) {
					defer wg.Done()

					out <- v.getByChunk(chunk, vector, n, norm)
				}(v.chunks[i])
			}
		}

		wg.Wait()
	}()

	return out
}

func (v *Vectors) getByChunk(
	chunk *vectorsChunk, vector Vector, n int, norm float64,
) minDistanceHeap {
	dh := slicesutil.MakeLimitHeap[*minDistanceHeapItem](n)

	tmp := vectorsPool.Get(len(vector))

	for _, r := range chunk.records {
		if r.vector == nil {
			continue
		}

		s := math.CosineSimilarity(vector, r.vector, norm, r.norm, *tmp)

		dh.Push(&minDistanceHeapItem{record: r, similarity: s})
	}

	vectorsPool.Put(tmp)

	return dh
}

func reduceHeaps(in <-chan minDistanceHeap, n int) minDistanceHeap {
	out := make(chan minDistanceHeap, 1)

	go func() {
		defer close(out)

		h := make(minDistanceHeap, 0, n)
		for cur := range in {
			for _, r := range cur {
				h.Push(r)
			}
		}

		out <- h
	}()

	return <-out
}

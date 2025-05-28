package vectors

import (
	"slices"
	"sync"

	"github.com/umk/llmservices/internal/math"
	slicesutil "github.com/umk/llmservices/internal/slices"
)

type ID int64

type Vector []float32

type vectorsHeader struct {
	ChunkSize int
}

type Vectors struct {
	header vectorsHeader

	chunks       []*vectorsChunk
	currentChunk *vectorsChunk
}

func NewVectors(chunkSize int) *Vectors {
	vectors := &Vectors{
		header: vectorsHeader{
			ChunkSize: chunkSize,
		},
		chunks: make([]*vectorsChunk, 1, 32),
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

	baseId := ID(len(v.chunks) * v.header.ChunkSize)

	v.currentChunk = newChunk(baseId, v.header.ChunkSize)
	v.chunks = append(v.chunks, v.currentChunk)

	return v.currentChunk.add(vector)
}

func (v *Vectors) Delete(id ID) bool {
	i, _ := slices.BinarySearchFunc(v.chunks, id, func(c *vectorsChunk, id ID) int {
		return int((c.BaseId - 1) - id)
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
		ids[i] = hr.record.Id
	}

	return ids
}

func (v *Vectors) Compact() {
	var destIndex, destRecordIndex int
	destChunk := v.chunks[destIndex]

	// Compact records by iterating over all chunks and their records
	for _, srcChunk := range v.chunks {
		for _, record := range srcChunk.Records {
			if record.Vector == nil {
				continue
			}

			// When destination chunk is full, move to next and reset index
			if destRecordIndex == cap(destChunk.Records) {
				destIndex++
				destChunk = v.chunks[destIndex]
				destRecordIndex = 0
			}

			// Write the valid record to the destination
			destChunk.Records[destRecordIndex] = record
			// Set new base ID for the chunk when writing its first record
			if destRecordIndex == 0 {
				destChunk.BaseId = record.Id
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
	for i := destRecordIndex; i < len(destChunk.Records); i++ {
		destChunk.Records[i] = nil
	}
	destChunk.Records = destChunk.Records[:destRecordIndex]
}

func (v *Vectors) Repack() *Vectors {
	// Create a new Vectors instance to hold compacted records
	vectors := &Vectors{
		header: vectorsHeader{
			ChunkSize: v.header.ChunkSize,
		},
		chunks: make([]*vectorsChunk, 1, 32),
	}

	// Initialize the first destination chunk
	destChunk := &vectorsChunk{
		Records: make([]*chunkRecord, 0, v.header.ChunkSize),
	}
	vectors.chunks[0] = destChunk

	// Iterate over all existing chunks
	for _, srcChunk := range v.chunks {
		// Iterate over all records in the current source chunk
		for _, record := range srcChunk.Records {
			if record.Vector == nil {
				continue
			}

			// If the current destination chunk is full, allocate a new one
			if len(destChunk.Records) == cap(destChunk.Records) {
				destChunk = &vectorsChunk{
					Records: make([]*chunkRecord, 0, v.header.ChunkSize),
				}
				vectors.chunks = append(vectors.chunks, destChunk)
			}
			// Set the base ID for a chunk when inserting its first record
			if len(destChunk.Records) == 0 {
				destChunk.BaseId = record.Id
			}

			// Append the valid record
			destChunk.Records = append(destChunk.Records, record)
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

	for _, r := range chunk.Records {
		if r.Vector == nil {
			continue
		}

		s := math.CosineSimilarity(vector, r.Vector, norm, r.Norm, *tmp)

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

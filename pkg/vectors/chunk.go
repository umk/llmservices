package vectors

import (
	"slices"

	"github.com/umk/llmservices/internal/math"
)

type vectorsChunk struct {
	baseId  ID
	records []*chunkRecord
}

type chunkRecord struct {
	id     ID
	vector Vector
	norm   float64
}

func newChunk(baseId ID, chunkSize int) *vectorsChunk {
	return &vectorsChunk{
		baseId:  baseId,
		records: make([]*chunkRecord, 0, chunkSize),
	}
}

func (vc *vectorsChunk) add(vector []float32) ID {
	if len(vc.records) == cap(vc.records) {
		return -1
	}

	id := vc.baseId + ID(len(vc.records))

	tmp := vectorsPool.Get(len(vector))

	vc.records = append(vc.records, &chunkRecord{
		id:     id,
		vector: vector,
		norm:   math.VectorNorm(vector, *tmp),
	})

	vectorsPool.Put(tmp)

	return id
}

func (vc *vectorsChunk) delete(id ID) bool {
	if i, ok := slices.BinarySearchFunc(vc.records, id, func(r *chunkRecord, id ID) int {
		return int(r.id - id)
	}); ok {
		r := vc.records[i]
		d := (r.vector != nil)
		r.vector = nil
		return d
	}
	return false
}

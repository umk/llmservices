package vectors

import (
	"slices"

	"github.com/umk/llmservices/internal/math"
)

type vectorsChunk struct {
	BaseID  ID
	Records []*chunkRecord
}

type chunkRecord struct {
	ID     ID
	Vector Vector
	Norm   float64
}

func newChunk(baseID ID, chunkSize int) *vectorsChunk {
	return &vectorsChunk{
		BaseID:  baseID,
		Records: make([]*chunkRecord, 0, chunkSize),
	}
}

func (vc *vectorsChunk) add(vector []float32) ID {
	if len(vc.Records) == cap(vc.Records) {
		return -1
	}

	id := vc.BaseID + ID(len(vc.Records))

	tmp := vectorsPool.Get(len(vector))

	vc.Records = append(vc.Records, &chunkRecord{
		ID:     id,
		Vector: vector,
		Norm:   math.VectorNorm(vector, *tmp),
	})

	vectorsPool.Put(tmp)

	return id
}

func (vc *vectorsChunk) delete(id ID) bool {
	if i, ok := slices.BinarySearchFunc(vc.Records, id, func(r *chunkRecord, id ID) int {
		return int(r.ID - id)
	}); ok {
		r := vc.Records[i]
		d := (r.Vector != nil)
		r.Vector = nil
		return d
	}
	return false
}

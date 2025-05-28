package vectors

import (
	"slices"

	"github.com/umk/llmservices/internal/math"
)

type vectorsChunk struct {
	BaseId  ID
	Records []*chunkRecord
}

type chunkRecord struct {
	Id     ID
	Vector Vector
	Norm   float64
}

func newChunk(baseId ID, chunkSize int) *vectorsChunk {
	return &vectorsChunk{
		BaseId:  baseId,
		Records: make([]*chunkRecord, 0, chunkSize),
	}
}

func (vc *vectorsChunk) add(vector []float32) ID {
	if len(vc.Records) == cap(vc.Records) {
		return -1
	}

	id := vc.BaseId + ID(len(vc.Records))

	tmp := vectorsPool.Get(len(vector))

	vc.Records = append(vc.Records, &chunkRecord{
		Id:     id,
		Vector: vector,
		Norm:   math.VectorNorm(vector, *tmp),
	})

	vectorsPool.Put(tmp)

	return id
}

func (vc *vectorsChunk) delete(id ID) bool {
	if i, ok := slices.BinarySearchFunc(vc.Records, id, func(r *chunkRecord, id ID) int {
		return int(r.Id - id)
	}); ok {
		r := vc.Records[i]
		d := (r.Vector != nil)
		r.Vector = nil
		return d
	}
	return false
}

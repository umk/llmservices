package vectorsdb

import (
	"sync"

	"github.com/umk/llmservices/internal/config"
	"github.com/umk/llmservices/pkg/vectors"
)

type Database[V any] struct {
	mu sync.RWMutex

	vectors *vectors.Vectors
	Data    map[vectors.ID]V

	vectorLength  int
	repackPercent int

	itemsCount   int
	deletesCount int

	repacking bool
}

type Record[V any] struct {
	Id     vectors.ID
	Vector vectors.Vector
	Data   V
}

func NewDatabase[V any](vectorLength int, options ...Option) *Database[V] {
	// Create default configuration
	cfg := &dbConfig{
		repackPercent: config.C.RepackPercent,
	}

	// Apply options
	for _, option := range options {
		option(cfg)
	}

	// Create database with the configured settings
	return &Database[V]{
		vectors:       vectors.NewVectors(128),
		Data:          make(map[vectors.ID]V),
		vectorLength:  vectorLength,
		repackPercent: cfg.repackPercent,
	}
}

func (db *Database[V]) Add(record Record[V]) Record[V] {
	db.mu.Lock()
	defer db.mu.Unlock()

	record.Vector = db.resizeVector(record.Vector)
	record = db.addRecord(record)

	db.itemsCount++
	return record
}

func (db *Database[V]) AddBatch(records []Record[V]) []Record[V] {
	db.mu.Lock()
	defer db.mu.Unlock()

	result := make([]Record[V], len(records))
	for i, record := range records {
		records[i].Vector = db.resizeVector(record.Vector)
		result[i] = db.addRecord(records[i])
	}

	db.itemsCount += len(records)
	return result
}

func (db *Database[V]) Delete(id vectors.ID) {
	db.mu.Lock()
	defer db.mu.Unlock()

	if db.deleteRecord(id) {
		db.increaseDeleteCount(1)
	}
}

func (db *Database[V]) DeleteBatch(ids []vectors.ID) {
	db.mu.Lock()
	defer db.mu.Unlock()

	deletedCount := 0

	for _, id := range ids {
		if db.deleteRecord(id) {
			deletedCount++
		}
	}

	db.increaseDeleteCount(deletedCount)
}

func (db *Database[V]) Get(vecs []vectors.Vector, n int) []Record[V] {
	db.mu.RLock()
	defer db.mu.RUnlock()

	for i, vec := range vecs {
		vecs[i] = db.resizeVector(vec)
	}

	ids := db.vectors.Get(vecs, n)

	r := make([]Record[V], len(ids))
	for i, id := range ids {
		r[i] = Record[V]{
			Id:   id,
			Data: db.Data[id],
		}
	}

	return r
}

// increaseDeleteCount increments the delete count and checks if it exceeds
// the threshold. If it does, it triggers a repack operation in a separate
// goroutine. Must be called from a write lock.
func (db *Database[V]) increaseDeleteCount(count int) {
	db.deletesCount += count

	totalItems := db.itemsCount + db.deletesCount
	if !db.repacking && totalItems > 0 && (db.deletesCount*100/totalItems) > db.repackPercent {
		db.repacking = true
		go func(vectors *vectors.Vectors) {
			db.mu.RLock()
			defer db.mu.RUnlock()

			db.vectors = vectors.Repack()

			db.itemsCount -= db.deletesCount
			db.deletesCount = 0

			db.repacking = false
		}(db.vectors)
	}
}

func (db *Database[V]) addRecord(record Record[V]) Record[V] {
	record.Id = db.vectors.Add(record.Vector)
	record.Vector = nil

	db.Data[record.Id] = record.Data

	return record
}

func (db *Database[V]) deleteRecord(id vectors.ID) bool {
	if db.vectors.Delete(id) {
		delete(db.Data, id)
		return true
	}
	return false
}

func (db *Database[V]) resizeVector(vec vectors.Vector) vectors.Vector {
	switch {
	case len(vec) > db.vectorLength:
		return vec[:db.vectorLength]
	case len(vec) < db.vectorLength:
		adjusted := make(vectors.Vector, db.vectorLength)
		copy(adjusted, vec)
		return adjusted
	default:
		return vec
	}
}

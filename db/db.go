package db

import (
	"sync"

	"github.com/umk/phishell/db/vector"
)

const itemsToDeletesRateToRepack = 10

type Database[V any] struct {
	vectors *vector.Vectors
	data    map[vector.ID]V
	mutex   sync.RWMutex

	itemsCount   int
	deletesCount int // number of nil items in chunks
}

type Record[V any] struct {
	ID     vector.ID
	Vector vector.Vector
	Data   V
}

func NewDatabase[V any]() *Database[V] {
	return &Database[V]{
		vectors: vector.NewVectors(128),
		data:    make(map[vector.ID]V),
	}
}

func (db *Database[V]) Add(record Record[V]) Record[V] {
	db.mutex.Lock()
	defer db.mutex.Unlock()

	record.ID = db.vectors.Add(record.Vector)
	db.data[record.ID] = record.Data

	db.itemsCount++
	return record
}

func (db *Database[V]) AddBatch(records []Record[V]) []Record[V] {
	db.mutex.Lock()
	defer db.mutex.Unlock()

	result := make([]Record[V], len(records))

	for i, record := range records {
		record.ID = db.vectors.Add(record.Vector)
		db.data[record.ID] = record.Data
		db.itemsCount++
		result[i] = record
	}
	return result
}

func (db *Database[V]) Delete(id vector.ID) {
	db.mutex.Lock()
	defer db.mutex.Unlock()

	if !db.vectors.Delete(id) {
		return
	}

	delete(db.data, id)

	db.increaseDeletePressure(1)
}

func (db *Database[V]) DeleteBatch(ids []vector.ID) {
	db.mutex.Lock()
	defer db.mutex.Unlock()

	deletedCount := 0

	for _, id := range ids {
		if db.vectors.Delete(id) {
			delete(db.data, id)
			deletedCount++
		}
	}

	db.increaseDeletePressure(deletedCount)
}

func (db *Database[V]) Get(vectors []vector.Vector, n int) []Record[V] {
	db.mutex.RLock()
	defer db.mutex.RUnlock()
	ids := db.vectors.Get(vectors, n)

	r := make([]Record[V], len(ids))
	for i, id := range ids {
		r[i] = Record[V]{
			ID:   id,
			Data: db.data[id],
		}
	}
	return r
}

func (db *Database[V]) increaseDeletePressure(count int) {
	db.deletesCount += count

	if db.deletesCount > (db.itemsCount / itemsToDeletesRateToRepack) {
		go func() {
			db.mutex.RLock()
			defer db.mutex.RUnlock()

			db.vectors = db.vectors.Repack()

			db.itemsCount -= db.deletesCount
			db.deletesCount = 0
		}()
	}
}

package db

import "github.com/umk/phishell/db/vector"

type Db[V any] struct {
	vectors *vector.Vectors
	data    map[vector.ID]V
}

type Record[V any] struct {
	ID     vector.ID
	Vector vector.Vector
	Data   V
}

func NewDb[V any]() *Db[V] {
	return &Db[V]{
		vectors: vector.NewVectors(128),
		data:    make(map[vector.ID]V),
	}
}

func (db *Db[V]) Add(record Record[V]) Record[V] {
	record.ID = db.vectors.Add(record.Vector)
	db.data[record.ID] = record.Data

	return record
}

func (db *Db[V]) Delete(id vector.ID) {
	db.vectors.Delete(id)
	delete(db.data, id)
}

func (db *Db[V]) Get(vectors []vector.Vector, n int) []Record[V] {
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

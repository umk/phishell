package vector

import "slices"

type vectorsChunk struct {
	baseID  ID
	records []*record
}

type record struct {
	id     ID
	vector Vector
	norm   float64
}

func searchChunk(c *vectorsChunk, id ID) int {
	return int(id - c.baseID)
}

func searchRecord(r *record, id ID) int {
	return int(id - r.id)
}

func newChunk(baseID ID, chunkSize int) *vectorsChunk {
	return &vectorsChunk{
		baseID:  baseID,
		records: make([]*record, 0, chunkSize),
	}
}

func (vc *vectorsChunk) add(vector []float32) ID {
	if len(vc.records) == cap(vc.records) {
		return -1
	}

	id := vc.baseID + ID(len(vc.records))

	tmp := vectorsPool.Get(len(vector))

	vc.records = append(vc.records, &record{
		id:     id,
		vector: vector,
		norm:   vectorNorm(vector, *tmp),
	})

	vectorsPool.Put(tmp)

	return id
}

func (vc *vectorsChunk) delete(id ID) {
	if i, ok := slices.BinarySearchFunc(vc.records, id, searchRecord); ok {
		vc.records[i] = nil
	}
}

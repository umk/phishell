package vector

import (
	"container/heap"
	"slices"
	"sync"

	"github.com/umk/phishell/util/slicesx"
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
	currentChunk := v.currentChunk

	for {
		id := currentChunk.add(vector)
		if id >= 0 {
			return id
		}

		if v.currentChunk == currentChunk {
			baseID := ID(len(v.chunks) * v.chunkSize)

			currentChunk = newChunk(baseID, v.chunkSize)

			id := v.currentChunk.add(vector)

			v.chunks = append(v.chunks, v.currentChunk)
			v.currentChunk = currentChunk

			return id
		}

		currentChunk = v.currentChunk
	}
}

func (v *Vectors) Delete(id ID) {
	i, _ := slices.BinarySearchFunc(v.chunks, id, searchChunk)
	v.chunks[i].delete(id)
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

func (v *Vectors) getHeaps(vectors []Vector, n int) <-chan maxDistanceHeap {
	out := make(chan maxDistanceHeap)

	go func() {
		defer close(out)

		var wg sync.WaitGroup

		for _, vector := range vectors {
			tmp := vectorsPool.Get(len(vector))
			norm := vectorNorm(vector, *tmp)
			vectorsPool.Put(tmp)

			for i := 0; i < len(v.chunks); i++ {
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
) maxDistanceHeap {
	dh := make(maxDistanceHeap, 0, n)

	heap.Init(&dh)

	tmp := vectorsPool.Get(len(vector))

	count := len(chunk.records)
	for i := 0; i < count; i++ {
		r := chunk.records[i]

		if r == nil {
			continue
		}

		s := cosineSimilarity(vector, r.vector, norm, r.norm, *tmp)
		slicesx.PushOrReplace(&dh, &maxDistanceHeapItem{record: r, similarity: s})
	}

	vectorsPool.Put(tmp)

	return dh
}

func reduceHeaps(in <-chan maxDistanceHeap, n int) maxDistanceHeap {
	out := make(chan maxDistanceHeap, 1)

	go func() {
		defer close(out)

		h := make(maxDistanceHeap, 0, n)
		for cur := range in {
			for _, r := range cur {
				slicesx.PushOrReplace(&h, r)
			}
		}

		out <- h
	}()

	return <-out
}

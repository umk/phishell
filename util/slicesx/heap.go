package slicesx

import "container/heap"

type SliceHeap[T SliceHeapItem[T]] []T

type SliceHeapItem[T any] interface {
	Less(another T) bool
}

func (h SliceHeap[T]) Len() int           { return len(h) }
func (h SliceHeap[T]) Less(i, j int) bool { return h[i].Less(h[j]) }
func (h SliceHeap[T]) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }

func (h *SliceHeap[T]) Push(x any) {
	*h = append(*h, x.(T))
}

func (h *SliceHeap[T]) Pop() any {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}

func PushOrReplace[T SliceHeapItem[T]](h *SliceHeap[T], item T) {
	if len(*h) < cap(*h) {
		heap.Push(h, item)
	} else {
		(*h)[0] = item
		heap.Fix(h, 0)
	}
}

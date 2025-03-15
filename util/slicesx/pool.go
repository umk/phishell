package slicesx

import "sync"

type SlicesPool[T any] struct {
	p    sync.Pool
	size int
}

func NewSlicesPool[T any](size int) *SlicesPool[T] {
	return &SlicesPool[T]{
		p: sync.Pool{
			New: func() any {
				slice := make([]T, 0, size)
				return &slice
			},
		},
		size: size,
	}
}

func (sp *SlicesPool[T]) Get(size int) *[]T {
	if size > sp.size {
		s := make([]T, size)
		return &s
	} else {
		s := sp.p.Get().(*[]T)
		*s = (*s)[:size]
		return s
	}
}

func (sp *SlicesPool[T]) Put(s *[]T) {
	if cap(*s) == sp.size {
		sp.p.Put(s)
	}
}

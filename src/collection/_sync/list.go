package _sync

import (
	"sync"
)

type Slice[T any] struct {
	list []T
	mu   sync.RWMutex
}

func NewSlice[T any]() *Slice[T] {
	return &Slice[T]{
		mu: sync.RWMutex{},
	}
}

func (s *Slice[T]) Append(vals ...T) {
	s.mu.Lock()
	s.list = append(s.list, vals...)
	s.mu.Unlock()
}

func (s *Slice[T]) List() []T {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.list
}

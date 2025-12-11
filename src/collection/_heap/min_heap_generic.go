package _heap

import (
	"container/heap"
)

func NewGeneric[T IOrdered[T]]() IMinHeap[T] {
	return newMinHeapGeneric[T]()
}

type arrGeneric[T IOrdered[T]] []T

func newArrGeneric[T IOrdered[T]](size int) *arrGeneric[T] {
	arr := make(arrGeneric[T], 0, size)
	return &arr
}

func (a *arrGeneric[T]) Len() int {
	return len(*a)
}

func (a *arrGeneric[T]) Less(i, j int) bool {
	return (*a)[i].Less((*a)[j])
}

func (a *arrGeneric[T]) Swap(i, j int) {
	(*a)[i], (*a)[j] = (*a)[j], (*a)[i]
}

func (a *arrGeneric[T]) Push(x any) {
	*a = append(*a, x.(T))
}

func (a *arrGeneric[T]) Pop() any {
	if len(*a) == 0 {
		return nil
	}
	last := (*a)[len(*a)-1]
	*a = (*a)[:len(*a)-1]
	return last
}

type minHeapGeneric[T IOrdered[T]] struct {
	arr *arrGeneric[T]
}

func newMinHeapGeneric[T IOrdered[T]]() *minHeapGeneric[T] {
	return &minHeapGeneric[T]{
		arr: newArrGeneric[T](16),
	}
}

func (h *minHeapGeneric[T]) Push(item T) {
	heap.Push(h.arr, item)
}

func (h *minHeapGeneric[T]) Pop() (T, bool) {
	if len(*h.arr) == 0 {
		var zero T
		return zero, false
	}
	return heap.Pop(h.arr).(T), true
}

func (h *minHeapGeneric[T]) Peek() (T, bool) {
	if len(*h.arr) == 0 {
		var zero T
		return zero, false
	}
	return (*h.arr)[0], true
}

func (h *minHeapGeneric[T]) Len() int {
	return len(*h.arr)
}

func (h *minHeapGeneric[T]) Slice() []T {
	return *h.arr
}

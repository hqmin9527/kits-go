package _heap

import (
	"cmp"
	"container/heap"
)

// 基础类型的最小堆实现

func New[T cmp.Ordered]() IMinHeap[T] {
	return newMinHeapPrimitive[T]()
}

type arrPrimitive[T cmp.Ordered] []T

func newArrPrimitive[T cmp.Ordered](size int) *arrPrimitive[T] {
	arr := make(arrPrimitive[T], 0, size)
	return &arr
}

func (a *arrPrimitive[T]) Len() int {
	return len(*a)
}

func (a *arrPrimitive[T]) Less(i, j int) bool {
	return (*a)[i] < (*a)[j]
}

func (a *arrPrimitive[T]) Swap(i, j int) {
	(*a)[i], (*a)[j] = (*a)[j], (*a)[i]
}

func (a *arrPrimitive[T]) Push(x any) {
	*a = append(*a, x.(T))
}

func (a *arrPrimitive[T]) Pop() any {
	if len(*a) == 0 {
		return nil
	}
	last := (*a)[len(*a)-1]
	*a = (*a)[:len(*a)-1]
	return last
}

type minHeapPrimitive[T cmp.Ordered] struct {
	arr *arrPrimitive[T]
}

func newMinHeapPrimitive[T cmp.Ordered]() *minHeapPrimitive[T] {
	return &minHeapPrimitive[T]{
		arr: newArrPrimitive[T](16),
	}
}

func (h *minHeapPrimitive[T]) Push(item T) {
	heap.Push(h.arr, item)
}

func (h *minHeapPrimitive[T]) Pop() (T, bool) {
	if len(*h.arr) == 0 {
		var zero T
		return zero, false
	}
	return heap.Pop(h.arr).(T), true
}

func (h *minHeapPrimitive[T]) Peek() (T, bool) {
	if len(*h.arr) == 0 {
		var zero T
		return zero, false
	}
	return (*h.arr)[0], true
}

func (h *minHeapPrimitive[T]) Len() int {
	return len(*h.arr)
}

func (h *minHeapPrimitive[T]) Slice() []T {
	return *h.arr
}

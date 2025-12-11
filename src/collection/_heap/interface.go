package _heap

type IOrdered[T any] interface {
	Less(T) bool
}

type IMinHeap[T any] interface {
	Len() int
	Push(item T)
	Pop() (T, bool)
	Peek() (T, bool)
	Slice() []T
}

type IMaxN[T any] interface {
	Push(item T)
	Pop() (T, bool)
	Slice() []T
	Len() int
	SortedSlice() []T
	Max() (T, bool)
	Min() (T, bool)
}

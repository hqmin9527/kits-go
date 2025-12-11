package _heap

// 泛型实现

import (
	"sort"

	"github.com/hqmin9527/kits-go/src/utils"
)

func NewMaxNGeneric[T IOrdered[T]](n int) IMaxN[T] {
	if n <= 0 {
		return nil
	}

	return &maxNGeneric[T]{
		h: NewGeneric[T](),
		n: n,
	}
}

type maxNGeneric[T IOrdered[T]] struct {
	h   IMinHeap[T]
	max T
	n   int
}

func (m *maxNGeneric[T]) Push(item T) {
	// 数量不足，直接push
	if m.h.Len() < m.n {
		m.pushItem(item)
		return
	}
	// 判断堆顶元素（最小元素）
	top, _ := m.h.Peek()
	if top.Less(item) {
		m.h.Pop()
		m.pushItem(item)
	}
}

func (m *maxNGeneric[T]) pushItem(item T) {
	m.h.Push(item)
	if m.Len() == 1 || m.max.Less(item) {
		m.max = item
	}
}

func (m *maxNGeneric[T]) Pop() (T, bool) {
	// pop的都是最小的元素，不用更新max
	return m.h.Pop()
}

func (m *maxNGeneric[T]) Slice() []T {
	return m.h.Slice()
}

func (m *maxNGeneric[T]) Len() int {
	return m.h.Len()
}

func (m *maxNGeneric[T]) SortedSlice() []T {
	// 复制一份避免影响原先数据结构
	s := utils.CloneSlice(m.Slice())
	sort.Slice(s, func(i, j int) bool {
		return s[i].Less(s[j])
	})
	return s
}

func (m *maxNGeneric[T]) Max() (T, bool) {
	if m.Len() == 0 {
		var zero T
		return zero, false
	}
	return m.max, true
}

func (m *maxNGeneric[T]) Min() (T, bool) {
	return m.h.Peek()
}

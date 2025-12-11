package _heap

// 基本类型的最大N堆

import (
	"cmp"
	"sort"

	"github.com/hqmin9527/kits-go/src/utils"
)

func NewMaxN[T cmp.Ordered](n int) IMaxN[T] {
	if n <= 0 {
		return nil
	}

	return &maxNPrimitive[T]{
		h: New[T](),
		n: n,
	}
}

type maxNPrimitive[T cmp.Ordered] struct {
	h   IMinHeap[T]
	max T
	n   int
}

func (m *maxNPrimitive[T]) Push(item T) {
	// 数量不足，直接push
	if m.h.Len() < m.n {
		m.pushItem(item)
		return
	}
	// 判断堆顶元素（最小元素）
	top, _ := m.h.Peek()
	if top < item {
		m.h.Pop()
		m.pushItem(item)
	}
}

func (m *maxNPrimitive[T]) pushItem(item T) {
	m.h.Push(item)
	if m.Len() == 1 || m.max < item {
		m.max = item
	}
}

func (m *maxNPrimitive[T]) Pop() (T, bool) {
	// pop的都是最小的元素，不用更新max
	return m.h.Pop()
}

func (m *maxNPrimitive[T]) Slice() []T {
	return m.h.Slice()
}

func (m *maxNPrimitive[T]) Len() int {
	return m.h.Len()
}

func (m *maxNPrimitive[T]) SortedSlice() []T {
	// 复制一份避免影响原先数据结构
	s := utils.CloneSlice(m.Slice())
	sort.Slice(s, func(i, j int) bool {
		return s[i] < s[j]
	})
	return s
}

func (m *maxNPrimitive[T]) Max() (T, bool) {
	if m.Len() == 0 {
		var zero T
		return zero, false
	}
	return m.max, true
}

func (m *maxNPrimitive[T]) Min() (T, bool) {
	return m.h.Peek()
}

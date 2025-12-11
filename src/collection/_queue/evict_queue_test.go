package _queue

import (
	"fmt"
	"testing"
)

func TestEvictQueue(t *testing.T) {
	q := NewEvictQueue[int](3)
	items := buildLinkedItems(10)

	q.EnQueue(items[1])
	q.EnQueue(items[2])
	q.EnQueue(items[3])
	printQueue(q) // [1 2 3]
	q.EnQueue(items[4])
	printQueue(q) // [2 3 4]
	q.EnQueue(items[5])
	printQueue(q) // [3 4 5]
	v, ok := q.DeQueue()
	fmt.Println(v.Value(), ok) // 3 true
	printQueue(q)              // [4 5]
	v, ok = q.DeQueue()
	fmt.Println(v.Value(), ok) // 4 true
	printQueue(q)              // [5]
	v, ok = q.DeQueue()
	fmt.Println(v.Value(), ok) // 5 true
	printQueue(q)              // []
	v, ok = q.DeQueue()
	fmt.Println(v.Value(), ok) // 0 false
	printQueue(q)              // []
	v, ok = q.DeQueue()
	fmt.Println(v.Value(), ok) // 0 false
	printQueue(q)              // []
}

func buildLinkedItems(n int) []*LinkedItem[int] {
	items := make([]*LinkedItem[int], n)
	for i := range items {
		items[i] = NewLinkedItem(i)
	}
	return items
}

func printQueue(q *EvictQueue[int]) {
	res := make([]int, 0)
	q.TraverseReverse(func(item *LinkedItem[int]) bool {
		res = append(res, item.Value())
		return true
	})

	fmt.Println(res)
}

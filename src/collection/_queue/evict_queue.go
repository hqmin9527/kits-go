package _queue

// 限制长度的队列，使用双向链表实现
// head --next--> tail --next--> nil
// tail --prev--> head --prev--> nil

type LinkedItem[T any] struct {
	value T
	next  *LinkedItem[T]
	prev  *LinkedItem[T]
}

func NewLinkedItem[T any](value T) *LinkedItem[T] {
	return &LinkedItem[T]{
		value: value,
	}
}

func (item *LinkedItem[T]) Value() T {
	return item.value
}

type EvictQueue[T any] struct {
	cap  int
	len  int
	head *LinkedItem[T] // index of in;
	tail *LinkedItem[T] // index of out;
}

func NewEvictQueue[T any](cap int) *EvictQueue[T] {
	q := new(EvictQueue[T])
	q.cap = cap
	q.len = 0
	return q
}

// 从head加入一个元素，如满了，从tail减去一个元素
func (q *EvictQueue[T]) EnQueue(item *LinkedItem[T]) (evict *LinkedItem[T]) {
	// DeQueue同时len--
	if q.IsFull() {
		evict, _ = q.DeQueue()
	}

	if q.len == 0 {
		q.head = item
		q.tail = item
	} else {
		q.head.prev = item
		item.next = q.head
		q.head = item
	}

	q.len++
	return
}

// 从tail减去一个元素，如为空，返回 (nil, false)
func (q *EvictQueue[T]) DeQueue() (*LinkedItem[T], bool) {
	if q.IsEmpty() {
		var zero T
		return NewLinkedItem(zero), false
	}

	var res *LinkedItem[T]

	res = q.tail
	if q.len == 1 {
		q.head = nil
		q.tail = nil
	} else {
		q.tail = q.tail.prev
		q.tail.next = nil
	}
	// 断开出列节点的连接
	res.prev = nil
	res.next = nil

	q.len--
	return res, true
}

func (q *EvictQueue[T]) IsEmpty() bool {
	return q.len == 0
}

func (q *EvictQueue[T]) IsFull() bool {
	return q.len == q.cap
}

func (q *EvictQueue[T]) Len() int {
	return q.len
}

func (q *EvictQueue[T]) Cap() int {
	return q.cap
}

// from head to tail
func (q *EvictQueue[T]) Traverse(visitor func(item *LinkedItem[T]) bool) {
	Traverse(q.head, visitor)
}

// from tail to head
func (q *EvictQueue[T]) TraverseReverse(visitor func(item *LinkedItem[T]) bool) {
	TraverseReverse(q.tail, visitor)
}

// til tail
func Traverse[T any](item *LinkedItem[T], visitor func(item *LinkedItem[T]) bool) {
	for item != nil {
		if !visitor(item) {
			break
		}
		item = item.next
	}
}

// til head
func TraverseReverse[T any](item *LinkedItem[T], visitor func(item *LinkedItem[T]) bool) {
	for item != nil {
		if !visitor(item) {
			break
		}
		item = item.prev
	}
}

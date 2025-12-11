package _queue

import (
	"sync"
	"time"
)

type TimestampedValue[T any] struct {
	Value     T
	Timestamp int64
}
type CircularQueue[T any] struct {
	items    []TimestampedValue[T]
	start    int
	end      int // end points to the next empty slot
	len      int
	capacity int
	expire   time.Duration
	mu       sync.Mutex
}

func NewCircularQueue[T any](capacity int, expire time.Duration) *CircularQueue[T] {
	return &CircularQueue[T]{
		items:    make([]TimestampedValue[T], capacity),
		capacity: capacity,
		expire:   expire,
	}
}
func (q *CircularQueue[T]) Push(item T) []T {
	q.mu.Lock()
	defer q.mu.Unlock()
	// Remove values that are older than expireTime
	currentTime := time.Now().UnixMilli()
	ts := currentTime - q.expire.Milliseconds()
	for q.len > 0 && ts > q.items[q.start].Timestamp {
		q.start = (q.start + 1) % q.capacity
		q.len--
	}
	newItem := TimestampedValue[T]{
		Value:     item,
		Timestamp: currentTime,
	}
	// If the queue is not full, insert the value
	if q.len < q.capacity {
		q.items[q.end] = newItem
		q.end = (q.end + 1) % q.capacity
		q.len++
	} else {
		// When the queue is full, we override the oldest value
		q.items[q.start] = newItem
		q.start = (q.start + 1) % q.capacity
		q.end = (q.end + 1) % q.capacity
	}

	// queue is full
	if q.len == q.capacity {
		allValues := make([]T, q.len)
		for i := 0; i < q.len; i++ {
			index := (q.start + i) % q.capacity
			allValues[i] = q.items[index].Value
		}
		q.clear()
		return allValues
	}
	return nil
}

func (q *CircularQueue[T]) clear() {
	// Reset the CircularQueue
	q.start = 0
	q.end = 0
	q.len = 0
}

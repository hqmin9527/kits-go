package _queue

import (
	"fmt"
	"testing"
	"time"
)

func TestCircularQueue(t *testing.T) {
	queue := NewCircularQueue[string](5, time.Second*10)
	// Example usage of the queue
	queue.Push("1")
	queue.Push("2")
	//time.Sleep(3 * time.Second)
	Print(queue)
	queue.Push("3")
	queue.Push("4")
	//time.Sleep(3 * time.Second)
	Print(queue)
	v := queue.Push("5")
	fmt.Println("Value returned:", v)
	queue.Push("6")
	//time.Sleep(3 * time.Second)
	Print(queue)
	queue.Push("7")
	v = queue.Push("8")
	fmt.Println("Value returned:", v)
	time.Sleep(3 * time.Second)
	Print(queue)
	queue.Push("9")
	Print(queue)
}

func Print(q *CircularQueue[string]) {
	q.mu.Lock()
	defer q.mu.Unlock()
	fmt.Println("Queue contents:")
	index := q.start
	for i := 0; i < q.len; i++ {
		item := q.items[index]
		fmt.Printf("Index %d: %v (stored at %v)\n", index, item.Value, item.Timestamp)
		index = (index + 1) % q.capacity
	}
}

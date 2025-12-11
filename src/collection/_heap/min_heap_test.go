package _heap

import "testing"

type Int int

func (i Int) Less(j Int) bool {
	return i < j
}

func TestMinHeapGeneric(t *testing.T) {
	h := NewGeneric[Int]()
	h.Push(Int(2))
	h.Push(Int(5))
	h.Push(Int(3))
	h.Push(Int(7))
	h.Push(Int(6))
	h.Push(Int(8))
	h.Push(Int(4))

	for {
		v, ok := h.Pop()
		if !ok {
			break
		}
		t.Log(v)
	}
}

func TestMinHeapPrimitive(t *testing.T) {
	h := New[int]()
	h.Push(2)
	h.Push(5)
	h.Push(3)
	h.Push(7)
	h.Push(6)
	h.Push(8)
	h.Push(4)

	for {
		v, ok := h.Pop()
		if !ok {
			break
		}
		t.Log(v)
	}
}

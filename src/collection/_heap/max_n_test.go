package _heap

import (
	"testing"
)

func TestMaxNGeneric(t *testing.T) {
	m := NewMaxNGeneric[Int](2)
	m.Push(Int(2))
	m.Push(Int(5))
	m.Push(Int(3))
	m.Push(Int(7))
	m.Push(Int(6))
	m.Push(Int(8))
	m.Push(Int(4))

	for {
		v, ok := m.Pop()
		if !ok {
			break
		}
		t.Log(v)
	}
}

func TestMaxNPrimitive(t *testing.T) {
	m := NewMaxN[int](2)
	m.Push(2)
	m.Push(5)
	m.Push(3)
	m.Push(7)
	m.Push(6)
	m.Push(8)
	m.Push(4)

	for {
		v, ok := m.Pop()
		if !ok {
			break
		}
		t.Log(v)
	}
}

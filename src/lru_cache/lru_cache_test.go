package lru_cache

import (
	"fmt"
	"reflect"
	"testing"
)

type Value string

func (v Value) Key() any {
	return v
}

func TestLruCache_String(t *testing.T) {
	var expectedEviction any
	echo := func(v any) {
		fmt.Println("on evict value:", v)
		if v != expectedEviction {
			t.Errorf("expect eviction: %v, got %v\n", expectedEviction, v)
		}
	}
	var res *[]any
	var reverse *[]any
	genVisitor := func(res *[]any) func(v any) bool {
		return func(v any) bool {
			*res = append(*res, v)
			return true
		}
	}

	res = &[]any{}
	reverse = &[]any{}
	cache := NewLruCache(3, echo)

	// 空Cache的遍历
	cache.RangeFromLatest(genVisitor(res))
	if expected := []any{}; !reflect.DeepEqual(expected, *res) {
		t.Errorf("expected: %#v, got: %#v\n", []any{}, *res)
	}
	cache.RangeFromEarliest(genVisitor(reverse))
	if expected := []any{}; !reflect.DeepEqual(expected, *reverse) {
		t.Errorf("expected: %#v, got: %#v\n", []any{}, *reverse)
	}

	// A=>B=>A=>C=>D
	A, B, C, D := Value("A"), Value("B"), Value("C"), Value("D")
	// 淘汰B
	expectedEviction = B
	cache.Add(A, B, A, C, D)
	// 最近：DCA
	cache.RangeFromLatest(genVisitor(res))
	if expected := []any{D, C, A}; !reflect.DeepEqual(expected, *res) {
		t.Errorf("expected: %v, got: %v\n", expected, *res)
	}
	// 最早：ACD
	cache.RangeFromEarliest(genVisitor(reverse))
	if expected := []any{A, C, D}; !reflect.DeepEqual(expected, *reverse) {
		t.Errorf("expected: %v, got: %v\n", expected, *reverse)
	}

	expectedEviction = A
	cache.Evict()
	expectedEviction = C
	cache.Evict()
	expectedEviction = D
	cache.Evict()

	cache.Evict()

	// 空Cache
	res = &[]any{}
	cache.RangeFromLatest(genVisitor(res))
	if expected := []any{}; !reflect.DeepEqual(expected, *res) {
		t.Errorf("expected: %#v, got: %#v\n", []any{}, *res)
	}
}

type Value2 struct {
	key   string
	value int
}

func (v Value2) Key() any {
	return v.key
}

func TestLruCache_Overwrite(t *testing.T) {
	A1 := Value2{"A", 1}
	A2 := Value2{"A", 2}

	cache := NewLruCache(3, nil)
	cache.Add(A1, A2)

	var res []Value2
	cache.RangeFromEarliest(func(value any) bool {
		res = append(res, value.(Value2))
		return true
	})

	if expected := []Value2{A2}; !reflect.DeepEqual(expected, res) {
		t.Errorf("expected: %v, got: %v\n", expected, res)
	}
}

package _map

import (
	"fmt"
	"testing"
)

func TestMapNumberPtr(t *testing.T) {
	m := map[string]any{
		"age": 12.0,
	}
	age := GetNumberPtr[int](m, "age")
	fmt.Println(*age)
}

func TestMapDirectPtr(t *testing.T) {
	m := map[string]any{
		"name": "jack",
	}
	name := GetDirectPtr[string](m, "name")
	fmt.Println(*name)
}

package _map

import (
	"fmt"
	"testing"
)

func TestGetIntPtr(t *testing.T) {
	m := map[string]any{
		"age": 12.0,
	}
	mm := Map(m)
	fmt.Println(*mm.GetIntPtr("age"))
}

package _slice

import (
	"fmt"
	"testing"
)

func TestGetSliceStr(t *testing.T) {
	s := []any{"aa"}
	ss := []any{s}
	b := GetSliceStr(ss, 0)
	fmt.Println(b)
}

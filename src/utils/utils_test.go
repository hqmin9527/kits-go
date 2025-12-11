package utils

import (
	"fmt"
	"testing"
)

func TestIf(t *testing.T) {
	a := If(true, "是", "否")
	fmt.Println(a)
}

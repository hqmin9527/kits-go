package _slice

import (
	"fmt"
	"testing"
)

func TestSetAdd(t *testing.T) {
	arr := []string{"A"}
	arr = SetAdd(arr, "B")
	fmt.Println(arr)
}

func TestPartition(t *testing.T) {
	var arr []string
	ps := Partition(arr, 100)
	fmt.Println("ps:", ps)

	arr = append(arr, "A", "B", "C")
	ps = Partition(arr, 2)
	fmt.Println("ps:", ps)
}

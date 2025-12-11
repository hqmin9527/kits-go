package retry

import (
	"fmt"
	"testing"
)

func TestRetry(t *testing.T) {
	Do(numberOut, []int{0, 1, 1, 2, 2, 2, 2, 2, 2, 2, 2, 2})
}

func TestTickDo(t *testing.T) {
	TickDo(numberOut, 2)
}

func TestRetryOnce(t *testing.T) {
	Do(printOnce, []int{0, 1, 1, 2, 2, 2, 2, 2, 2, 2, 2, 2})
	TickDo(printOnce, 10)
}

var outNumber = 5
var i = 0

func printOnce() bool {
	fmt.Println("success")
	return true
}

func numberOut() bool {
	i++
	fmt.Println(i)
	return i > outNumber
}

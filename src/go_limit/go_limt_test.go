package go_limit

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/pkg/errors"
)

func TestGoLimit_RunError(t *testing.T) {

	limit := New(10)
	for i := 0; i < 10; i++ {
		ti := i
		limit.RunError(func() error {
			return errors.New(strconv.Itoa(ti))
		})
	}

	limit.Wait()
	errs := limit.ListErrors()
	fmt.Println(errs)
	firstErr := limit.FirstError()
	fmt.Println(firstErr)

	if len(errs) != 10 {
		t.Error("length of errs should be 10\n")
	}

	if errs[0] != firstErr {
		t.Errorf("errs[0] should equal firstErr errs[0]: %v, firsrErr: %v\n", errs[0], firstErr)
	}
}

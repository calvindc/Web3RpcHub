package errorslist

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestError(t *testing.T) {
	a := assert.New(t)
	var els ErrorList

	for i := 0; i < 3; i++ {
		els.Errs = append(els.Errs, errors.New(fmt.Sprintf("%d", i)))
	}

	fmt.Println(els.Error())
	a.Equal("Err(1/3): 0\nErr(2/3): 1\nErr(3/3): 2\n", els.Error())
}

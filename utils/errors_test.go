package utils

import (
	"errors"
	"github.com/AWoelfel/go-multi-ssh/tests/assert"
	"testing"
)

func TestFromErrors(t *testing.T) {

	t.Run("nil", func(t *testing.T) {})
	t.Run("nil errors", func(t *testing.T) {})
	t.Run("mixed errors", func(t *testing.T) {})
	t.Run("nested errors", func(t *testing.T) {})
	t.Run("errors", func(t *testing.T) {})

}

func TestErrors_Error(t *testing.T) {

	a := errors.New("A")
	b := errors.New("B")
	c := errors.New("C")

	x := Errors{a, b, c}

	assert.EqualValues(t, "A\nB\nC", x.Error())

}

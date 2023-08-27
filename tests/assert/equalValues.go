package assert

import "testing"

func EqualValues[T comparable](t *testing.T, expected T, actual T) bool {

	if expected != actual {
		t.Logf("expected %v but got %v", expected, actual)
		t.Fail()
	}

	return !t.Failed()
}

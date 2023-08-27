package assert

import (
	"testing"
)

func EqualArrayValues[T comparable](t *testing.T, expected []T, actual []T) bool {
	return EqualArrayValuesWithComparator(t, expected, actual, EqualValues[T])

}

func EqualArrayValuesWithComparator[T comparable](t *testing.T, expected []T, actual []T, cmp func(t *testing.T, expected T, actual T) bool) bool {

	if len(expected) != len(actual) {
		t.Logf("expected a array of length %d got one with length %d", len(expected), len(actual))
		t.Fail()
		return false
	}

	for i := 0; i < len(expected); i++ {
		cmp(t, expected[i], actual[i])
	}

	return !t.Failed()

}

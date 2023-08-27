package assert

import "testing"

func ArrayLen[T any](t *testing.T, expectedLen int, actual []T) bool {
	if len(actual) != expectedLen {
		t.Logf("expected a array of length %d got one with length %d", expectedLen, len(actual))
		t.Fail()
	}

	return !t.Failed()
}

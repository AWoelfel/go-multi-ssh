package assert

import "testing"

func NoError(t *testing.T, err error) bool {

	if err != nil {
		t.Logf("expected no error got %q", err.Error())
		t.Fail()
	}

	return !t.Failed()
}

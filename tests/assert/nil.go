package assert

import "testing"

func Nil(t *testing.T, obj interface{}) bool {

	if obj != nil {
		t.Logf("expected nil got %q", obj)
		t.Fail()
	}

	return !t.Failed()
}

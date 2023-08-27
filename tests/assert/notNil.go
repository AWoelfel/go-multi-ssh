package assert

import "testing"

func NotNil(t *testing.T, obj interface{}) bool {

	if obj == nil {
		t.Log("expected an object got nil")
		t.Fail()
	}

	return !t.Failed()
}

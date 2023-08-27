package connection

import (
	"github.com/AWoelfel/go-multi-ssh/tests/assert"
	"testing"
)

func TestConnection(t *testing.T) {
	target := ClientContext{
		Namespace: "ns",
		Pod:       "po",
		Container: "co",
	}

	assert.EqualValues(t, "ns/po:co", target.Label())
}

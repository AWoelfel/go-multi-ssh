package connection

import (
	"github.com/AWoelfel/go-multi-ssh/tests/assert"
	"testing"
)

func TestConnection(t *testing.T) {
	t.Run("slash prefixed", func(t *testing.T) {
		target := ClientContext{
			ID:        "0123456789ABCDEF",
			Container: "/some-container",
		}

		assert.EqualValues(t, "0123456789AB (some-container)", target.Label())
	})
	t.Run("non slash prefixed", func(t *testing.T) {
		target := ClientContext{
			ID:        "0123456789ABCDEF",
			Container: "some-container",
		}

		assert.EqualValues(t, "0123456789AB (some-container)", target.Label())
	})
}

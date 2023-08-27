package output

import (
	"github.com/AWoelfel/go-multi-ssh/tests/assert"
	"github.com/muesli/termenv"
	"io"
	"math/rand"
	"testing"
)

type testRemoteMeta struct {
	l string
	c termenv.ANSI256Color
}

func (t *testRemoteMeta) Label() string {
	return t.l
}

func (t *testRemoteMeta) Color() termenv.ANSI256Color {
	return t.c
}

func TestLineByLineChannel(t *testing.T) {

	t.Run("properly ended", func(t *testing.T) {
		r, w := io.Pipe()
		testRemote := &testRemoteMeta{l: "test", c: termenv.ANSI256Color(int(int64(rand.Intn(256))))}
		c := LineByLineChannel(r, testRemote)

		go func() {
			w.Write([]byte("A"))
			w.Write([]byte("\n"))
			w.Write([]byte("B"))
			w.Write([]byte("B\nCC"))
			w.Write([]byte("C"))
			w.Write([]byte("\n"))
			w.Close()
		}()

		var result []LineByLineChannelPayload
		for e := range c {
			result = append(result, e)
		}

		assert.AssertObjectsEqual(t,
			[]LineByLineChannelPayload{
				{
					Remote: testRemote,
					Line:   "A",
				},
				{
					Remote: testRemote,
					Line:   "BB",
				},
				{
					Remote: testRemote,
					Line:   "CCC",
				},
			},
			result,
		)
	})

	t.Run("nonproperly ended", func(t *testing.T) {
		r, w := io.Pipe()
		testRemote := &testRemoteMeta{l: "test", c: termenv.ANSI256Color(int(int64(rand.Intn(256))))}
		c := LineByLineChannel(r, testRemote)

		go func() {
			w.Write([]byte("A"))
			w.Write([]byte("\n"))
			w.Write([]byte("B"))
			w.Write([]byte("B\nCC"))
			w.Close()
		}()

		var result []LineByLineChannelPayload
		for e := range c {
			result = append(result, e)
		}

		assert.AssertObjectsEqual(t,
			[]LineByLineChannelPayload{
				{
					Remote: testRemote,
					Line:   "A",
				},
				{
					Remote: testRemote,
					Line:   "BB",
				},
				{
					Remote: testRemote,
					Line:   "CC",
				},
			},
			result,
		)
	})

}

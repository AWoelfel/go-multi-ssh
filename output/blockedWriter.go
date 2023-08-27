package output

import (
	"github.com/muesli/termenv"
	"io"
	"strings"
	"sync"
)

type blockedWriter struct {
	sync.Mutex
	wg         sync.WaitGroup
	Out        io.Writer
	lastWriter RemoteMeta
}

func NewBlockedWriter(out io.Writer) *blockedWriter {
	return &blockedWriter{
		Out:        out,
		lastWriter: nil,
	}
}

func (bw *blockedWriter) AttachSource(sourceChannel <-chan LineByLineChannelPayload) {
	bw.wg.Add(1)
	go func(c <-chan LineByLineChannelPayload) {
		for e := range c {
			bw.WriteLine(e.Remote, e.Line)
		}
		bw.wg.Done()
	}(sourceChannel)
}

func (bw *blockedWriter) WriteLine(conn RemoteMeta, line string) {
	bw.Lock()
	defer bw.Unlock()

	if bw.lastWriter != conn {
		clientHeader := conn.Label()
		clientHeaderLine := strings.Repeat("-", len(clientHeader))

		bw.Out.Write([]byte("\n"))
		bw.Out.Write([]byte(termenv.String(clientHeader).Foreground(conn.Color()).Bold().String()))
		bw.Out.Write([]byte("\n"))
		bw.Out.Write([]byte(termenv.String(clientHeaderLine).Foreground(conn.Color()).Bold().String()))
		bw.Out.Write([]byte("\n"))
	}

	bw.Out.Write([]byte(line))
	bw.Out.Write([]byte("\n"))
	bw.lastWriter = conn
}

func (bw *blockedWriter) Wait() {
	bw.wg.Wait()
}

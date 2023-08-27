package output

import (
	"fmt"
	"github.com/muesli/termenv"
	"io"
	"sync"
)

type inlineWriter struct {
	sync.Mutex
	wg  sync.WaitGroup
	Out io.Writer
}

func NewInlineWriter(out io.Writer) *inlineWriter {
	return &inlineWriter{
		Out: out,
	}
}

func (i *inlineWriter) AttachSource(sourceChannel <-chan LineByLineChannelPayload) {
	i.wg.Add(1)
	go func(c <-chan LineByLineChannelPayload) {
		for e := range c {
			i.WriteLine(e.Remote, e.Line)
		}
		i.wg.Done()
	}(sourceChannel)
}

func (i *inlineWriter) WriteLine(conn RemoteMeta, line string) {
	i.Lock()
	defer i.Unlock()
	i.Out.Write([]byte(termenv.String(fmt.Sprintf("%s : ", conn.Label())).Foreground(conn.Color()).Bold().String()))
	i.Out.Write([]byte(line))
	i.Out.Write([]byte("\n"))
}

func (i *inlineWriter) Wait() {
	i.wg.Wait()
}

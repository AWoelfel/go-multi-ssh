package output

import (
	"io"
	"sync"
)

type unlabeledWriter struct {
	sync.Mutex
	wg  sync.WaitGroup
	Out io.Writer
}

func NewUnlabeledWriter(out io.Writer) *unlabeledWriter {
	return &unlabeledWriter{Out: out}
}

func (bw *unlabeledWriter) AttachSource(sourceChannel <-chan LineByLineChannelPayload) {
	bw.wg.Add(1)
	go func(c <-chan LineByLineChannelPayload) {
		for e := range c {
			bw.WriteLine(e.Remote, e.Line)
		}
		bw.wg.Done()
	}(sourceChannel)
}

func (bw *unlabeledWriter) WriteLine(_ RemoteMeta, line string) {
	bw.Lock()
	defer bw.Unlock()

	bw.Out.Write([]byte(line))
	bw.Out.Write([]byte("\n"))
}

func (bw *unlabeledWriter) Wait() {
	bw.wg.Wait()
}

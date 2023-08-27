package tests

import (
	"github.com/AWoelfel/go-multi-ssh/output"
	"sync"
)

type FakeOutputSink struct {
	sync.Mutex
	wg sync.WaitGroup

	Lines map[string][]string
}

func (f *FakeOutputSink) AttachSource(sourceChannel <-chan output.LineByLineChannelPayload) {
	f.wg.Add(1)
	go func(c <-chan output.LineByLineChannelPayload) {
		for e := range c {
			f.WriteLine(e.Remote, e.Line)
		}
		f.wg.Done()
	}(sourceChannel)
}

func (f *FakeOutputSink) WriteLine(source output.RemoteMeta, line string) {
	f.Lock()
	defer f.Unlock()

	var connectionLines []string

	if _, found := f.Lines[source.Label()]; found {
		connectionLines = f.Lines[source.Label()]
	}

	connectionLines = append(connectionLines, line)
	f.Lines[source.Label()] = connectionLines
}

func (f *FakeOutputSink) Wait() {
	f.wg.Wait()
}

func NewFakeOutputSink() *FakeOutputSink {
	return &FakeOutputSink{Lines: make(map[string][]string)}
}

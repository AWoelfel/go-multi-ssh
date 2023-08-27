package connectionTest

import (
	"fmt"
	"github.com/AWoelfel/go-multi-ssh/output"
	"github.com/muesli/termenv"
	"io"
	"sync"
)

type ConnectionStatusWriter struct {
	sync.Mutex
	wg  sync.WaitGroup
	Out io.Writer
}

func (i *ConnectionStatusWriter) AttachSource(sourceChannel <-chan output.LineByLineChannelPayload) {
	i.wg.Add(1)
	go func(c <-chan output.LineByLineChannelPayload) {
		for e := range c {
			i.WriteLine(e.Remote, e.Line)
		}
		i.wg.Done()
	}(sourceChannel)
}

func (i *ConnectionStatusWriter) WriteLine(conn output.RemoteMeta, line string) {
	i.Lock()
	defer i.Unlock()

	i.Out.Write([]byte(line))
	i.Out.Write([]byte("     "))
	i.Out.Write([]byte(termenv.String(fmt.Sprintf("%s", conn.Label())).Foreground(conn.Color()).Bold().String()))
	i.Out.Write([]byte("\n"))
}

func (i *ConnectionStatusWriter) Wait() {
	i.wg.Wait()
}

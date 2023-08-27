package output

import (
	"bufio"
	"github.com/muesli/termenv"
	"io"
)

type RemoteMeta interface {
	Label() string
	Color() termenv.ANSI256Color
}

type LineByLineChannelPayload struct {
	Remote RemoteMeta
	Line   string
}

func LineByLineChannel(source io.Reader, remote RemoteMeta) <-chan LineByLineChannelPayload {

	c := make(chan LineByLineChannelPayload)

	go func(source io.Reader, payload RemoteMeta, targetChannel chan LineByLineChannelPayload) {
		s := bufio.NewScanner(source)
		s.Split(bufio.ScanLines)
		for s.Scan() {
			targetChannel <- LineByLineChannelPayload{Remote: payload, Line: s.Text()}
		}
		close(c)
	}(source, remote, c)

	return c
}

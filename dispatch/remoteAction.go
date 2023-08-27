package dispatch

import (
	"context"
	"github.com/AWoelfel/go-multi-ssh/output"
)

type RemoteAction interface {
	Remote() output.RemoteMeta
	Run(ctx context.Context, stdOut output.RemoteActionOutputSink, stdErr output.RemoteActionOutputSink) (err error)
}

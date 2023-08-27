package output

import (
	"context"
	"github.com/AWoelfel/go-multi-ssh/config"
	"io"
	"os"
)

func DefaultContextWithWriters(ctx context.Context, label config.OutputLabel) context.Context {
	return ContextWithWriters(ctx, label, os.Stdout, os.Stderr)
}

func ContextWithWriters(ctx context.Context, label config.OutputLabel, stdout io.Writer, stderr io.Writer) context.Context {

	switch label {
	case config.BlockOutputLabel:
		ctx = WithOutputSink(ctx, StdOutChannel, NewBlockedWriter(stdout))
		return WithOutputSink(ctx, StdErrChannel, NewBlockedWriter(stderr))
	case config.InlineOutputLabel:
		ctx = WithOutputSink(ctx, StdOutChannel, NewInlineWriter(stdout))
		return WithOutputSink(ctx, StdErrChannel, NewInlineWriter(stderr))
	case config.NoOutputLabel:
		fallthrough
	default:
		ctx = WithOutputSink(ctx, StdOutChannel, NewUnlabeledWriter(stdout))
		return WithOutputSink(ctx, StdErrChannel, NewUnlabeledWriter(stderr))
	}

}

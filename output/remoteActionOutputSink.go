package output

import "context"

type RemoteActionOutputSink interface {
	AttachSource(sourceChannel <-chan LineByLineChannelPayload)
	WriteLine(source RemoteMeta, line string)
	Wait()
}

type TargetChanel int

const (
	StdOutChannel TargetChanel = iota
	StdErrChannel
)

type outWriterContextKey TargetChanel

func WithOutputSink(ctx context.Context, channel TargetChanel, target RemoteActionOutputSink) context.Context {
	return context.WithValue(ctx, channel, target)
}

func OutputSinkFromContext(ctx context.Context, channel TargetChanel) RemoteActionOutputSink {
	return ctx.Value(channel).(RemoteActionOutputSink)
}

package connectionTest

import (
	"context"
	"fmt"
	"github.com/AWoelfel/go-multi-ssh/cmd/docker/config"
	"github.com/AWoelfel/go-multi-ssh/cmd/docker/connection"
	"github.com/AWoelfel/go-multi-ssh/cmd/docker/dockerClient"
	"github.com/AWoelfel/go-multi-ssh/connectionTest"
	"github.com/AWoelfel/go-multi-ssh/dispatch"
	"github.com/AWoelfel/go-multi-ssh/output"
)

type connectionTestAction struct {
	client *connection.ClientContext
}

func (c *connectionTestAction) Remote() output.RemoteMeta {
	return c.client
}

func (c *connectionTestAction) Run(_ context.Context, stdOut output.RemoteActionOutputSink, _ output.RemoteActionOutputSink) error {

	stdOut.WriteLine(c.client, connectionTest.Ok)

	return nil
}

func Execute(ctx context.Context) error {

	dClient := dockerClient.ClientFromContext(ctx)

	ctx = output.WithOutputSink(ctx, output.StdOutChannel, &connectionTest.ConnectionStatusWriter{Out: dClient.Out()})
	ctx = output.WithOutputSink(ctx, output.StdErrChannel, &connectionTest.ConnectionStatusWriter{Out: dClient.Err()})

	cmdConfig := config.FromContext(ctx)

	targetHosts, err := cmdConfig.Clients(ctx)
	if err != nil {
		return fmt.Errorf("unable to resolve targets (%w)", err)
	}

	var allClients = make([]dispatch.RemoteAction, len(targetHosts), len(targetHosts))
	for i := 0; i < len(targetHosts); i++ {
		allClients[i] = &connectionTestAction{&targetHosts[i]}
	}

	return dispatch.Dispatch(ctx, allClients)
}

package connectionTest

import (
	"context"
	"fmt"
	"github.com/AWoelfel/go-multi-ssh/cmd/ssh/config"
	connection2 "github.com/AWoelfel/go-multi-ssh/cmd/ssh/connection"
	"github.com/AWoelfel/go-multi-ssh/connectionTest"
	"github.com/AWoelfel/go-multi-ssh/dispatch"
	"github.com/AWoelfel/go-multi-ssh/output"
	"github.com/AWoelfel/go-multi-ssh/utils"
	"golang.org/x/crypto/ssh"
	"os"
)

type connectionTestAction struct {
	client *connection2.ClientContext
}

func (c *connectionTestAction) Remote() output.RemoteMeta {
	return c.client
}

func (c *connectionTestAction) Run(_ context.Context, stdOut output.RemoteActionOutputSink, _ output.RemoteActionOutputSink) (err error) {
	var client *ssh.Client
	client, err = connection2.OpenClient(c.client)
	defer func() { err = utils.FromErrors(err, client.Close()) }()

	line := connectionTest.Ok

	if err != nil {
		line = connectionTest.NotOk
	}

	stdOut.WriteLine(c.client, line)

	if err != nil {
		return nil
	}

	return
}

func Execute(ctx context.Context) error {
	ctx = output.WithOutputSink(ctx, output.StdOutChannel, &connectionTest.ConnectionStatusWriter{Out: os.Stdout})
	ctx = output.WithOutputSink(ctx, output.StdErrChannel, &connectionTest.ConnectionStatusWriter{Out: os.Stderr})

	cmdConfig := config.FromContext(ctx)

	targetHosts, err := cmdConfig.Clients()
	if err != nil {
		return fmt.Errorf("unable to resolve targets (%w)", err)
	}

	var allClients = make([]dispatch.RemoteAction, len(targetHosts), len(targetHosts))
	for i := 0; i < len(targetHosts); i++ {
		allClients[i] = &connectionTestAction{connection2.NewClientContext(targetHosts[i])}
	}

	return dispatch.Dispatch(ctx, allClients)
}

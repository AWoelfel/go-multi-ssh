package single

import (
	"context"
	"fmt"
	"github.com/AWoelfel/go-multi-ssh/cmd/docker/config"
	"github.com/AWoelfel/go-multi-ssh/cmd/docker/connection"
	"github.com/AWoelfel/go-multi-ssh/cmd/docker/dockerClient"
	"github.com/AWoelfel/go-multi-ssh/dispatch"
	"github.com/AWoelfel/go-multi-ssh/output"
	"github.com/AWoelfel/go-multi-ssh/utils"
	"github.com/docker/cli/cli/command/container"
	"io"
)

type singleCommandAction struct {
	client *connection.ClientContext
	cmd    []string
}

func (c *singleCommandAction) Remote() output.RemoteMeta {
	return c.client
}

func (c *singleCommandAction) Run(ctx context.Context, stdOut output.RemoteActionOutputSink, stdErr output.RemoteActionOutputSink) (err error) {

	dClientFactory := dockerClient.DockerCliFactoryFromContext(ctx)

	stdoutR, stdoutW := io.Pipe()
	stderrR, stderrW := io.Pipe()

	defer func() {
		err = utils.FromErrors(err, stdoutW.Close())
		err = utils.FromErrors(err, stderrW.Close())
	}()

	dClient, err := dClientFactory.AttachCli(ctx, stdoutW, stderrW)
	if err != nil {
		return fmt.Errorf("unable to create docker cli (%w)", err)
	}

	stdOut.AttachSource(output.LineByLineChannel(stdoutR, c.client))
	stdErr.AttachSource(output.LineByLineChannel(stderrR, c.client))

	exec := container.NewExecOptions()
	exec.Container = c.client.ID
	exec.Command = c.cmd

	err = container.RunExec(dClient, exec)
	return

}

func Execute(ctx context.Context, args []string) error {

	cmdConfig := config.FromContext(ctx)

	targetHosts, err := cmdConfig.Clients(ctx)
	if err != nil {
		return fmt.Errorf("unable to resolve targets (%w)", err)
	}

	var allClients = make([]dispatch.RemoteAction, len(targetHosts), len(targetHosts))
	for i := 0; i < len(targetHosts); i++ {
		allClients[i] = &singleCommandAction{&targetHosts[i], args}
	}

	return dispatch.Dispatch(ctx, allClients)
}

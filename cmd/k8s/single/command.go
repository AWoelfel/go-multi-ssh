package single

import (
	"context"
	"fmt"
	"github.com/AWoelfel/go-multi-ssh/cmd/k8s/config"
	"github.com/AWoelfel/go-multi-ssh/cmd/k8s/connection"
	"github.com/AWoelfel/go-multi-ssh/cmd/k8s/k8sClient"
	"github.com/AWoelfel/go-multi-ssh/dispatch"
	"github.com/AWoelfel/go-multi-ssh/output"
	"github.com/AWoelfel/go-multi-ssh/utils"
	"io"
	"k8s.io/client-go/tools/remotecommand"
	"strings"
)

type singleCommandAction struct {
	client *connection.ClientContext
	cmd    string
}

func (c *singleCommandAction) Remote() output.RemoteMeta {
	return c.client
}

func (c *singleCommandAction) Run(ctx context.Context, stdOut output.RemoteActionOutputSink, stdErr output.RemoteActionOutputSink) (err error) {

	execFactory := k8sClient.ExecutorFactoryFromContext(ctx)

	var exec remotecommand.Executor
	exec, err = execFactory(ctx, c.client.Namespace, c.client.Pod, c.client.Container, c.cmd)
	if err != nil {
		return
	}

	stdoutR, stdoutW := io.Pipe()
	stderrR, stderrW := io.Pipe()

	defer func() {
		err = utils.FromErrors(err, stdoutW.Close())
		err = utils.FromErrors(err, stderrW.Close())
	}()

	stdOut.AttachSource(output.LineByLineChannel(stdoutR, c.client))
	stdErr.AttachSource(output.LineByLineChannel(stderrR, c.client))

	err = exec.StreamWithContext(ctx,
		remotecommand.StreamOptions{
			Stdin:  nil,
			Stdout: stdoutW,
			Stderr: stderrW,
		})
	return

}

func Execute(ctx context.Context, args []string) error {

	cmd := strings.Join(args, " ")

	cmdConfig := config.FromContext(ctx)

	targetHosts, err := cmdConfig.Clients(ctx)
	if err != nil {
		return fmt.Errorf("unable to resolve targets (%w)", err)
	}

	var allClients = make([]dispatch.RemoteAction, len(targetHosts), len(targetHosts))
	for i := 0; i < len(targetHosts); i++ {
		allClients[i] = &singleCommandAction{&targetHosts[i], cmd}
	}

	return dispatch.Dispatch(ctx, allClients)
}

package connectionTest

import (
	"context"
	"fmt"
	"github.com/AWoelfel/go-multi-ssh/cmd/k8s/config"
	"github.com/AWoelfel/go-multi-ssh/cmd/k8s/connection"
	"github.com/AWoelfel/go-multi-ssh/cmd/k8s/k8sClient"
	"github.com/AWoelfel/go-multi-ssh/connectionTest"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"os"

	"github.com/AWoelfel/go-multi-ssh/dispatch"
	"github.com/AWoelfel/go-multi-ssh/output"
)

type connectionTestAction struct {
	client *connection.ClientContext
}

func (c *connectionTestAction) Remote() output.RemoteMeta {
	return c.client
}

func (c *connectionTestAction) Run(ctx context.Context, stdOut output.RemoteActionOutputSink, stdErr output.RemoteActionOutputSink) (err error) {

	kClient := k8sClient.ClientFromContext(ctx)

	targetPod, err := kClient.CoreV1().Pods(c.client.Namespace).Get(ctx, c.client.Pod, metaV1.GetOptions{})
	if err != nil {
		return
	}

	line := connectionTest.NotOk

	for i := 0; i < len(targetPod.Status.ContainerStatuses); i++ {
		if targetPod.Status.ContainerStatuses[i].Name == c.client.Container {
			if targetPod.Status.ContainerStatuses[i].Ready {
				line = connectionTest.Ok
			}
			break
		}
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

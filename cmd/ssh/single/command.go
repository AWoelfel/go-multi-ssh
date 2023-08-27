package single

import (
	"context"
	"fmt"
	"github.com/AWoelfel/go-multi-ssh/cmd/ssh/config"
	connection2 "github.com/AWoelfel/go-multi-ssh/cmd/ssh/connection"
	"github.com/AWoelfel/go-multi-ssh/dispatch"
	"github.com/AWoelfel/go-multi-ssh/output"
	"github.com/AWoelfel/go-multi-ssh/utils"
	"golang.org/x/crypto/ssh"
	"io"
	"strings"
)

type singleCommandAction struct {
	cmd    string
	client *connection2.ClientContext
}

func (c *singleCommandAction) Remote() output.RemoteMeta {
	return c.client
}

func (c *singleCommandAction) Run(_ context.Context, stdOut output.RemoteActionOutputSink, stdErr output.RemoteActionOutputSink) (err error) {

	var client *ssh.Client
	client, err = connection2.OpenClient(c.client)
	if err != nil {
		return
	}

	defer func() { err = utils.FromErrors(err, client.Close()) }()

	var session *ssh.Session
	session, err = client.NewSession()

	if err != nil {
		return
	}

	defer func() {
		if sessionCloseError := session.Close(); sessionCloseError != io.EOF {
			err = utils.FromErrors(err, sessionCloseError)
		}
	}()

	var remoteStdOutReader io.Reader
	remoteStdOutReader, err = session.StdoutPipe()

	if err != nil {
		err = fmt.Errorf("unable to attach to remote stdout (%w)", err)
		return
	}

	stdOut.AttachSource(output.LineByLineChannel(remoteStdOutReader, c.client))

	var remoteStdErrReader io.Reader
	remoteStdErrReader, err = session.StderrPipe()

	if err != nil {
		err = fmt.Errorf("unable to attach to remote stdout (%w)", err)
		return
	}

	stdErr.AttachSource(output.LineByLineChannel(remoteStdErrReader, c.client))

	err = session.Run(c.cmd)

	return
}

func Execute(ctx context.Context, args []string) error {

	cmd := strings.Join(args, " ")

	cmdConfig := config.FromContext(ctx)

	targetHosts, err := cmdConfig.Clients()
	if err != nil {
		return fmt.Errorf("unable to resolve targets (%w)", err)
	}

	var allClients = make([]dispatch.RemoteAction, len(targetHosts), len(targetHosts))
	for i := 0; i < len(targetHosts); i++ {
		allClients[i] = &singleCommandAction{cmd: cmd, client: connection2.NewClientContext(targetHosts[i])}
	}

	return dispatch.Dispatch(ctx, allClients)
}

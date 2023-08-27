package dockerClient

import (
	"context"
	"errors"
	"fmt"
	"github.com/docker/cli/cli/command"
	"io"
)

type dockerClientContextKey int

var dockerClientContextKeyVal dockerClientContextKey = 0

func WithClient(ctx context.Context, dockerClient command.Cli) context.Context {
	return context.WithValue(ctx, dockerClientContextKeyVal, dockerClient)
}

func ClientFromContext(ctx context.Context) command.Cli {
	return ctx.Value(dockerClientContextKeyVal).(command.Cli)
}

type dockerCliFactoryContextKey int

var dockerCliFactoryContextKeyVal dockerCliFactoryContextKey = 0

type DockerCliFactory interface {
	AttachCli(ctx context.Context, stdOut io.Writer, stdErr io.Writer) (command.Cli, error)
}

type defaultDockerCliFactory struct {
}

func (d *defaultDockerCliFactory) AttachCli(ctx context.Context, stdOut io.Writer, stdErr io.Writer) (command.Cli, error) {
	dClient := ClientFromContext(ctx)

	//hijack docker stdErr/stdOut
	hijackedDockerClient, success := dClient.(*command.DockerCli)
	if !success {
		return nil, errors.New("unable to convert docker cli")
	}

	hijackedDockerClientCopy := *hijackedDockerClient

	err := hijackedDockerClientCopy.Apply(command.WithOutputStream(stdOut))
	if err != nil {
		return nil, fmt.Errorf("unable to set output stream (%w)", err)
	}
	err = hijackedDockerClientCopy.Apply(command.WithErrorStream(stdErr))
	if err != nil {
		return nil, fmt.Errorf("unable to set error stream (%w)", err)
	}

	return &hijackedDockerClientCopy, nil
}

func DockerCliFactoryFromContext(ctx context.Context) DockerCliFactory {

	existing := ctx.Value(dockerCliFactoryContextKeyVal)
	if existing != nil {
		return existing.(DockerCliFactory)
	}

	//fallback to default factory
	return &defaultDockerCliFactory{}
}

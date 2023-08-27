package single

import (
	"context"
	"errors"
	"github.com/AWoelfel/go-multi-ssh/cmd/k8s/connection"
	"github.com/AWoelfel/go-multi-ssh/cmd/k8s/k8sClient"
	"github.com/AWoelfel/go-multi-ssh/tests"
	"github.com/AWoelfel/go-multi-ssh/tests/assert"
	"github.com/muesli/termenv"
	"k8s.io/client-go/kubernetes/fake"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/remotecommand"
	"testing"
)

type testExecutor struct {
	t             *testing.T
	executorError error
}

func (t *testExecutor) Stream(options remotecommand.StreamOptions) error {
	panic("you are not supposed to be called")
}

func (t *testExecutor) StreamWithContext(_ context.Context, options remotecommand.StreamOptions) error {

	assert.Nil(t.t, options.Stdin)
	assert.NotNil(t.t, options.Stderr)
	assert.NotNil(t.t, options.Stdout)

	return t.executorError
}

func TestExecCmd(t *testing.T) {

	ctx := context.Background()

	expectedServer := "https://host-from-cluster.com"
	expectedToken := "token-from-cluster"
	expectedTokenFile := "tokenfile-from-cluster"
	expectedCAFile := "/path/to/ca-from-cluster.crt"

	restConfig := &rest.Config{
		Host:            expectedServer,
		BearerToken:     expectedToken,
		BearerTokenFile: expectedTokenFile,
		TLSClientConfig: rest.TLSClientConfig{
			CAFile: expectedCAFile,
		},
	}

	ctx = k8sClient.WithClient(ctx, fake.NewSimpleClientset())
	ctx = k8sClient.WithConfig(ctx, restConfig)

	t.Run("test k8s cmd exec", func(t *testing.T) {

		testStdOut := tests.NewFakeOutputSink()
		testStdErr := tests.NewFakeOutputSink()

		target := singleCommandAction{
			client: &connection.ClientContext{
				Namespace: "namespace",
				Pod:       "pod",
				Container: "container",
				Col:       termenv.ANSI256Color(25),
			},
			cmd: "command",
		}

		execFactory := func(ctx context.Context, namespace string, pod string, container string, cmd string) (remotecommand.Executor, error) {
			assert.EqualValues(t, "namespace", namespace)
			assert.EqualValues(t, "pod", pod)
			assert.EqualValues(t, "container", container)
			assert.EqualValues(t, "command", cmd)
			return &testExecutor{t, nil}, nil
		}

		ctx := k8sClient.WithExecutorFactory(ctx, execFactory)

		err := target.Run(ctx, testStdOut, testStdErr)

		assert.NoError(t, err)

	})

	t.Run("test k8s cmd exec (execution fails)", func(t *testing.T) {

		testStdOut := tests.NewFakeOutputSink()
		testStdErr := tests.NewFakeOutputSink()

		target := singleCommandAction{
			client: &connection.ClientContext{
				Namespace: "namespace",
				Pod:       "pod",
				Container: "container",
				Col:       termenv.ANSI256Color(25),
			},
			cmd: "command",
		}

		execFactory := func(ctx context.Context, namespace string, pod string, container string, cmd string) (remotecommand.Executor, error) {
			assert.EqualValues(t, "namespace", namespace)
			assert.EqualValues(t, "pod", pod)
			assert.EqualValues(t, "container", container)
			assert.EqualValues(t, "command", cmd)
			return &testExecutor{t, errors.New("some-execution-error")}, nil
		}

		ctx := k8sClient.WithExecutorFactory(ctx, execFactory)

		err := target.Run(ctx, testStdOut, testStdErr)

		assert.EqualValues(t, "some-execution-error", err.Error())

	})

}

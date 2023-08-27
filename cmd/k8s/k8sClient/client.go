package k8sClient

import (
	"context"
	"fmt"
	v1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/remotecommand"
	"log"
	"os"
	"path/filepath"
)

type kubernetesClientContextKey int

var kubernetesClientContextKeyVal kubernetesClientContextKey = 0

type kubernetesRestConfigContextKey int

var kubernetesRestConfigContextKeyVal kubernetesRestConfigContextKey = 0

type kubernetesCommandExecutorContextKey int

var kubernetesCommandExecutorContextKeyVal kubernetesCommandExecutorContextKey = 0

func AttachClient(ctx context.Context) (context.Context, error) {
	configPath := ""

	if k8sConfig, found := os.LookupEnv("KUBECONFIG"); found {
		configPath = k8sConfig
	}

	if len(configPath) == 0 {
		userHome, err := os.UserHomeDir()
		if err != nil {
			return nil, fmt.Errorf("unable to resolve current user home dir (%w)", err)
		}
		configPath = filepath.Join(userHome, ".kube", "config")
	}

	config, err := clientcmd.BuildConfigFromFlags("", configPath)
	if err != nil {
		log.Panicln("failed to create K8s config")
	}

	k8sClient, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, fmt.Errorf("unable to create k8s client (%w)", err)
	}

	ctx = WithConfig(ctx, config)
	ctx = WithClient(ctx, k8sClient)

	return ctx, nil
}

func WithClient(ctx context.Context, k8sClient kubernetes.Interface) context.Context {
	return context.WithValue(ctx, kubernetesClientContextKeyVal, k8sClient)
}

func WithConfig(ctx context.Context, restConfig *rest.Config) context.Context {
	return context.WithValue(ctx, kubernetesRestConfigContextKeyVal, restConfig)
}

func ClientFromContext(ctx context.Context) kubernetes.Interface {
	return ctx.Value(kubernetesClientContextKeyVal).(kubernetes.Interface)
}

func RestConfigFromContext(ctx context.Context) *rest.Config {
	return ctx.Value(kubernetesRestConfigContextKeyVal).(*rest.Config)
}

type ExecutorFactory func(ctx context.Context, namespace string, pod string, container string, cmd string) (remotecommand.Executor, error)

func defaultExecutorFactory(ctx context.Context, namespace string, pod string, container string, cmd string) (remotecommand.Executor, error) {
	kClient := ClientFromContext(ctx)
	restConfig := RestConfigFromContext(ctx)

	req := kClient.CoreV1().RESTClient().Post().Resource("pods").Name(pod).
		Namespace(namespace).SubResource("exec")
	option := &v1.PodExecOptions{
		Command:   []string{cmd},
		Stdin:     false,
		Stdout:    true,
		Stderr:    true,
		TTY:       false,
		Container: container,
	}

	req.VersionedParams(option, scheme.ParameterCodec)
	return remotecommand.NewSPDYExecutor(restConfig, "POST", req.URL())
}

func ExecutorFactoryFromContext(ctx context.Context) ExecutorFactory {

	existing := ctx.Value(kubernetesCommandExecutorContextKeyVal)

	if existing != nil {
		return existing.(ExecutorFactory)
	}

	return defaultExecutorFactory
}

func WithExecutorFactory(ctx context.Context, exFactory ExecutorFactory) context.Context {
	return context.WithValue(ctx, kubernetesCommandExecutorContextKeyVal, exFactory)
}

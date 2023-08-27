package config

import (
	"context"
	"fmt"
	"github.com/AWoelfel/go-multi-ssh/cmd/k8s/connection"
	"github.com/AWoelfel/go-multi-ssh/cmd/k8s/k8sClient"
	"github.com/muesli/termenv"
	coreV1 "k8s.io/api/core/v1"
	"math/rand"
	"strings"

	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const defaultContainerAnnotation = "kubectl.kubernetes.io/default-container"

// Clients returns a prefiltered set of target clients based on the Configuration.
func (c *Configuration) Clients(ctx context.Context) ([]connection.ClientContext, error) {

	kClient := k8sClient.ClientFromContext(ctx)

	labelSelector := strings.Join(c.SearchLabels, ",")
	podList, err := kClient.CoreV1().Pods(c.Namespace).List(ctx, metaV1.ListOptions{LabelSelector: labelSelector})
	if err != nil {
		return nil, fmt.Errorf("unable to list pods (%w)", err)
	}
	var result []connection.ClientContext

	for i := 0; i < len(podList.Items); i++ {

		pod := podList.Items[i]

		if pod.Status.Phase == coreV1.PodSucceeded ||
			pod.Status.Phase == coreV1.PodFailed {
			continue
		}

		targetContainerName := pod.Spec.Containers[0].Name

		if defaultContainer, found := pod.Annotations[defaultContainerAnnotation]; found {
			targetContainerName = defaultContainer
		}

		result = append(result, connection.ClientContext{
			Namespace: pod.Namespace,
			Pod:       pod.Name,
			Container: targetContainerName,
			Col:       termenv.ANSI256Color(rand.Intn(256)),
		})
	}

	return result, nil
}

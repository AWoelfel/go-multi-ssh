package config

import (
	"context"
	"github.com/AWoelfel/go-multi-ssh/cmd/k8s/connection"
	"github.com/AWoelfel/go-multi-ssh/cmd/k8s/k8sClient"
	"github.com/AWoelfel/go-multi-ssh/tests/assert"
	coreV1 "k8s.io/api/core/v1"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/fake"
	"testing"
)

func TestConfiguration_Clients(t *testing.T) {

	ctx := context.Background()

	clientComparator := func(t *testing.T, expected, actual connection.ClientContext) bool {
		assert.EqualValues(t, expected.Namespace, actual.Namespace)
		assert.EqualValues(t, expected.Pod, actual.Pod)
		assert.EqualValues(t, expected.Container, actual.Container)
		return !t.Failed()
	}

	ctx = k8sClient.WithClient(
		ctx,
		fake.NewSimpleClientset(
			&coreV1.Pod{
				ObjectMeta: metaV1.ObjectMeta{
					Namespace: "foo",
					Name:      "some-pod",
				},
				Spec: coreV1.PodSpec{
					Containers: []coreV1.Container{
						{Name: "some-container"},
					},
				},
			},
			&coreV1.Pod{
				ObjectMeta: metaV1.ObjectMeta{
					Namespace: "foo",
					Name:      "some-other-pod",
					Labels:    map[string]string{"component": "selective-A"},
				},
				Spec: coreV1.PodSpec{
					Containers: []coreV1.Container{
						{Name: "some-other-container"},
					},
				},
			},
			&coreV1.Pod{
				ObjectMeta: metaV1.ObjectMeta{
					Namespace:   "bar",
					Name:        "some-bar-pod",
					Labels:      map[string]string{"component": "selective-A", "app": "selective-B"},
					Annotations: map[string]string{"kubectl.kubernetes.io/default-container": "some-bar-container"},
				},
				Spec: coreV1.PodSpec{
					Containers: []coreV1.Container{
						{Name: "do-not-pick-this-container"},
						{Name: "some-bar-container"},
					},
				},
			},
			&coreV1.Pod{
				ObjectMeta: metaV1.ObjectMeta{
					Namespace: "foo",
					Name:      "some-PodFailed-pod",
				},
				Spec: coreV1.PodSpec{
					Containers: []coreV1.Container{
						{Name: "some-container"},
					},
				},
				Status: coreV1.PodStatus{
					Phase: coreV1.PodFailed,
				},
			},
			&coreV1.Pod{
				ObjectMeta: metaV1.ObjectMeta{
					Namespace: "foo",
					Name:      "some-PodSucceeded-pod",
				},
				Spec: coreV1.PodSpec{
					Containers: []coreV1.Container{
						{Name: "some-container"},
					},
				},
				Status: coreV1.PodStatus{
					Phase: coreV1.PodSucceeded,
				},
			},
		),
	)

	t.Run("all pods of all namespaces", func(t *testing.T) {

		cfg := DefaultConfig

		clients, err := cfg.Clients(ctx)
		assert.NoError(t, err)
		assert.EqualArrayValuesWithComparator(
			t,
			[]connection.ClientContext{
				{Namespace: "bar", Pod: "some-bar-pod", Container: "some-bar-container"},
				{Namespace: "foo", Pod: "some-other-pod", Container: "some-other-container"},
				{Namespace: "foo", Pod: "some-pod", Container: "some-container"},
			},
			clients,
			clientComparator)

	})

	t.Run("all pods of single namespace", func(t *testing.T) {
		cfg := DefaultConfig
		cfg.Namespace = "foo"

		clients, err := cfg.Clients(ctx)
		assert.NoError(t, err)
		assert.EqualArrayValuesWithComparator(
			t,
			[]connection.ClientContext{
				{Namespace: "foo", Pod: "some-other-pod", Container: "some-other-container"},
				{Namespace: "foo", Pod: "some-pod", Container: "some-container"},
			},
			clients,
			clientComparator)

	})

	t.Run("selective pods of single namespace", func(t *testing.T) {
		cfg := DefaultConfig
		cfg.Namespace = "foo"
		cfg.SearchLabels = []string{"component=selective-A"}

		clients, err := cfg.Clients(ctx)
		assert.NoError(t, err)
		assert.EqualArrayValuesWithComparator(
			t,
			[]connection.ClientContext{
				{Namespace: "foo", Pod: "some-other-pod", Container: "some-other-container"},
			},
			clients,
			clientComparator)

	})

	t.Run("selective pods of all namespaces", func(t *testing.T) {
		cfg := DefaultConfig
		cfg.SearchLabels = []string{"component=selective-A"}

		clients, err := cfg.Clients(ctx)
		assert.NoError(t, err)
		assert.EqualArrayValuesWithComparator(
			t,
			[]connection.ClientContext{
				{Namespace: "bar", Pod: "some-bar-pod", Container: "some-bar-container"},
				{Namespace: "foo", Pod: "some-other-pod", Container: "some-other-container"},
			},
			clients,
			clientComparator)

	})

	t.Run("selective pods (multiple labelselectors)", func(t *testing.T) {
		cfg := DefaultConfig
		cfg.SearchLabels = []string{"component=selective-A", "app=selective-B"}

		clients, err := cfg.Clients(ctx)
		assert.NoError(t, err)
		assert.EqualArrayValuesWithComparator(
			t,
			[]connection.ClientContext{
				{Namespace: "bar", Pod: "some-bar-pod", Container: "some-bar-container"},
			},
			clients,
			clientComparator)
	})

}

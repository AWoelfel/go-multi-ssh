package config

import (
	"context"
	"fmt"
	"github.com/AWoelfel/go-multi-ssh/cmd/docker/connection"
	"github.com/AWoelfel/go-multi-ssh/cmd/docker/dockerClient"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/muesli/termenv"
	"math/rand"
)

// Clients returns a prefiltered set of target clients based on the Configuration.
func (c *Configuration) Clients(ctx context.Context) ([]connection.ClientContext, error) {

	dClient := dockerClient.ClientFromContext(ctx)

	labelSelector := make([]filters.KeyValuePair, len(c.SearchLabels), len(c.SearchLabels))
	for i := 0; i < len(c.SearchLabels); i++ {
		labelSelector[i] = filters.Arg("label", c.SearchLabels[i])
	}

	containers, err := dClient.Client().ContainerList(ctx, types.ContainerListOptions{
		Filters: filters.NewArgs(labelSelector...),
	})
	if err != nil {
		return nil, fmt.Errorf("unable to list containers (%w)", err)
	}
	var result []connection.ClientContext

	for i := 0; i < len(containers); i++ {

		container := containers[i]
		if container.State != "running" {
			continue
		}

		result = append(result, connection.ClientContext{
			ID:        container.ID,
			Container: container.Names[0],
			Col:       termenv.ANSI256Color(rand.Intn(256)),
		})
	}

	return result, nil
}

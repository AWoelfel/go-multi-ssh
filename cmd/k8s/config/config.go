package config

import (
	"context"
	coreV1 "k8s.io/api/core/v1"
)

var DefaultConfig Configuration = Configuration{
	SearchLabels: nil,
	Namespace:    coreV1.NamespaceAll,
}

type Configuration struct {
	SearchLabels []string
	Namespace    string
}

type configurationContextKey int

var configurationContextKeyValue configurationContextKey = 0

func (c *Configuration) WithContext(ctx context.Context) context.Context {
	return context.WithValue(ctx, configurationContextKeyValue, c)
}

func FromContext(ctx context.Context) *Configuration {
	return ctx.Value(configurationContextKeyValue).(*Configuration)
}

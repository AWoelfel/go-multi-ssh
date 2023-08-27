package config

import (
	"context"
)

var DefaultConfig Configuration = Configuration{
	SearchLabels: nil,
}

type Configuration struct {
	SearchLabels []string
}

type configurationContextKey int

var configurationContextKeyValue configurationContextKey = 0

func (c *Configuration) WithContext(ctx context.Context) context.Context {
	return context.WithValue(ctx, configurationContextKeyValue, c)
}

func FromContext(ctx context.Context) *Configuration {
	return ctx.Value(configurationContextKeyValue).(*Configuration)
}

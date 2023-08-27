package config

import (
	"context"
	"path/filepath"
)

const DefaultIndexFile = "ssh_config"

var DefaultConfig Configuration = Configuration{
	IndexFile:   filepath.Join(".", DefaultIndexFile),
	IncludeTags: nil,
	ExcludeTags: nil,
}

type Configuration struct {
	IncludeTags []string
	ExcludeTags []string
	IndexFile   string
}

type configurationContextKey int

var configurationContextKeyValue configurationContextKey = 0

func (c *Configuration) WithContext(ctx context.Context) context.Context {
	return context.WithValue(ctx, configurationContextKeyValue, c)
}

func FromContext(ctx context.Context) *Configuration {
	return ctx.Value(configurationContextKeyValue).(*Configuration)
}

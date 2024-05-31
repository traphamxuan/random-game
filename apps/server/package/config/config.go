package config

import (
	"context"

	"github.com/caarlos0/env"
	"github.com/traphamxuan/gobs"
)

type Configuration struct {
}

var _ gobs.IService = (*Configuration)(nil)

func (c *Configuration) Init(context.Context, *gobs.Component) error {
	return nil
}

func NewConfig(ctx context.Context) *Configuration {
	return &Configuration{}
}

func (c *Configuration) ParseConfig(result interface{}) error {
	return env.Parse(result)
}

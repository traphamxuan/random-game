package config

import (
	"context"

	"github.com/caarlos0/env"
	servicemanager "github.com/traphamxuan/random-game/package/service_manager"
)

type Configuration struct {
}

var _ servicemanager.IService = (*Configuration)(nil)

func NewConfig(ctx context.Context) *Configuration {
	return &Configuration{}
}

func (c *Configuration) Name() string {
	return "Configuration"
}

func (c *Configuration) Setup(ctx context.Context) error {
	return nil
}

func (c *Configuration) Start(ctx context.Context) error {
	return nil
}

func (c *Configuration) Stop(ctx context.Context) error {
	return nil
}

func (c *Configuration) ParseConfig(result interface{}) error {
	return env.Parse(result)
}

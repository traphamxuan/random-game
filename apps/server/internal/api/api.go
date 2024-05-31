package api

import (
	"context"

	"game-random-api/package/logger"

	"github.com/traphamxuan/gobs"
)

type AppInstance struct {
	log *logger.Logger
	bs  *gobs.Bootstrap
}

func NewAPI(ctx context.Context) *AppInstance {
	api := &AppInstance{}
	api.log = &logger.Logger{}
	if err := api.log.Setup(ctx); err != nil {
		panic(err)
	}
	bs := gobs.NewBootstrap(gobs.DefaultConfig)
	bs.AddOrPanic(&Router{})
	api.bs = bs
	return api
}

func (a *AppInstance) Setup(ctx context.Context) error {
	return a.bs.Setup(ctx)
}

func (a *AppInstance) Start(ctx context.Context) error {
	return a.bs.Start(ctx)
}

func (a *AppInstance) Interrupt(ctx context.Context) {
	a.log.Warn("Interrupting app instance")
}

func (a *AppInstance) Stop(ctx context.Context) {
	a.bs.Stop(ctx)
}

func (a *AppInstance) Deinit(ctx context.Context) {
	a.bs.Deinit(ctx)
}

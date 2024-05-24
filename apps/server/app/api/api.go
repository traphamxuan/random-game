package api

import (
	"context"

	"github.com/traphamxuan/random-game/package/cache"
	"github.com/traphamxuan/random-game/package/config"
	"github.com/traphamxuan/random-game/package/logger"
	"github.com/traphamxuan/random-game/package/orm"
	"github.com/traphamxuan/random-game/package/s3"
	servicemanager "github.com/traphamxuan/random-game/package/service_manager"
)

type AppInstance struct {
	log           *logger.Logger
	sm            *servicemanager.ServiceManager
	setupPriority []string
	startPriority []string
	stopPriority  []string
}

func NewAPI(ctx context.Context) AppInstance {
	sm := servicemanager.NewServiceManager()
	log := logger.NewLogger(ctx)
	sm.Add(log)
	sm.Add(config.NewConfig(ctx))
	sm.Add(cache.NewRedis(ctx, sm))
	sm.Add(orm.NewOrm(ctx, sm))
	sm.Add(s3.NewS3(ctx, sm))
	sm.Add(NewRouter(ctx, sm))

	return AppInstance{
		log:           log,
		sm:            sm,
		setupPriority: []string{"Logger", "Config"},
		startPriority: []string{"Logger", "Config", "Redis", "Orm", "S3"},
		stopPriority:  []string{"Router"},
	}
}

func (a *AppInstance) Setup(ctx context.Context) error {
	onSetup := func(ctx context.Context, s servicemanager.IService, key string, err error) {
		if err != nil {
			a.log.Error("Failed to setup service %s: %v", key, err)
		} else {
			a.log.Info("Service %s setup successfully", key)
		}
	}
	return a.sm.Setup(ctx, &onSetup, a.startPriority...)
}

func (a *AppInstance) Start(ctx context.Context) error {
	onStart := func(ctx context.Context, s servicemanager.IService, key string, err error) {
		if err != nil {
			a.log.Error("Failed to start service %s: %v", key, err)
		} else {
			a.log.Info("Service %s started successfully", key)
		}
	}
	return a.sm.Start(ctx, &onStart, a.startPriority...)
}

func (a *AppInstance) Interrupt(ctx context.Context) {
	a.log.Warn("Interrupting app instance")
}

func (a *AppInstance) Stop(ctx context.Context) {
	onStop := func(ctx context.Context, s servicemanager.IService, key string, err error) {
		if err != nil {
			a.log.Error("Failed to stop service %s: %v", key, err)
		} else {
			a.log.Info("Service %s stopped successfully", key)
		}
	}
	a.sm.Stop(ctx, &onStop, a.stopPriority...)
}

func (a *AppInstance) Deinit(ctx context.Context) {}

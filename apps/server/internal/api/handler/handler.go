package handler

import (
	"context"
	"game-random-api/internal/api/service"
	"game-random-api/package/logger"
	"reflect"

	"github.com/labstack/echo/v4"
	"github.com/traphamxuan/gobs"
)

type IHandler interface {
	SetRoute(context.Context, *echo.Group) error
}

type Handler struct {
	log            *logger.Logger
	Authentication *Authentication
	HealthCheck    *HealthCheck
}

var handlers = []gobs.IService{
	&service.Service{},
	&logger.Logger{},
	&Authentication{},
	&HealthCheck{},
}

var _ gobs.IService = (*Handler)(nil)

// Init implements gobs.IService.
func (h *Handler) Init(ctx context.Context, sb *gobs.Component) error {
	sb.Deps = handlers
	onSetup := func(ctx context.Context, dependencies []gobs.IService, _ []gobs.CustomService) error {
		handlers = dependencies
		h.log = dependencies[1].(*logger.Logger)
		h.Authentication = dependencies[2].(*Authentication)
		h.HealthCheck = dependencies[3].(*HealthCheck)
		return nil
	}
	sb.OnSetup = &onSetup
	return nil
}

func (h *Handler) SetRoute(ctx context.Context, eg *echo.Group) error {
	v := reflect.ValueOf(h).Elem()
	for i := 0; i < v.NumField(); i++ {
		// Ignore if field is private
		if v.Type().Field(i).PkgPath != "" {
			continue
		}
		if handler, ok := v.Field(i).Interface().(IHandler); ok {
			h.log.Infof("Setting route for %s", v.Type().Field(i).Name)
			if err := handler.SetRoute(ctx, eg); err != nil {
				return err
			}
		}
	}
	return nil
}

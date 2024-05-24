package controller

import (
	"context"

	"github.com/labstack/echo/v4"
	"github.com/traphamxuan/random-game/app/api/middleware"
	servicemanager "github.com/traphamxuan/random-game/package/service_manager"
)

type IController interface {
	Setup(e *echo.Echo) error
}

type Controller struct {
	endpoints []IController
}

func NewController(ctx context.Context,
	mw *middleware.Middleware,
	sm *servicemanager.ServiceManager,
) *Controller {
	return &Controller{
		endpoints: []IController{
			NewUser(ctx, mw, sm),
		},
	}
}

func (c *Controller) Setup(e *echo.Echo) error {
	for _, ep := range c.endpoints {
		if err := ep.Setup(e); err != nil {
			return err
		}
	}
	return nil
}

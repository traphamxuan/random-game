package handler

import (
	"context"
	"game-random-api/package/logger"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/traphamxuan/gobs"
)

type HealthCheck struct {
	log *logger.Logger
}

var _ gobs.IService = (*HealthCheck)(nil)
var _ IHandler = (*HealthCheck)(nil)

func (a *HealthCheck) Init(ctx context.Context, sb *gobs.Component) error {
	sb.Deps = []gobs.IService{
		&logger.Logger{},
	}
	onSetup := func(ctx context.Context, dependencies []gobs.IService, _ []gobs.CustomService) error {
		a.log = dependencies[0].(*logger.Logger)
		return nil
	}
	sb.OnSetup = &onSetup
	return nil
}

// Route implements IHandler.
func (a *HealthCheck) SetRoute(ctx context.Context, e *echo.Group) error {
	g := e.Group("/ping")
	g.GET("", a.Ping)
	return nil
}

func (a *HealthCheck) Ping(c echo.Context) error {
	return c.JSON(http.StatusOK, "OK")
}

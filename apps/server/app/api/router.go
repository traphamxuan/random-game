package api

import (
	"context"
	"net/http"

	"github.com/labstack/echo/v4"
	eMiddleware "github.com/labstack/echo/v4/middleware"
	"github.com/traphamxuan/random-game/app/api/controller"
	"github.com/traphamxuan/random-game/app/api/middleware"
	"github.com/traphamxuan/random-game/app/api/validator"
	servicemanager "github.com/traphamxuan/random-game/package/service_manager"
)

type Router struct {
	middleware *middleware.Middleware
	validator  *validator.Validator
	controller *controller.Controller
	engine     *echo.Echo
}

var _ servicemanager.IService = (*Router)(nil)

func NewRouter(ctx context.Context, sm *servicemanager.ServiceManager) *Router {
	v := validator.NewValidator(ctx, sm)
	m := middleware.NewMiddleware(ctx, sm)
	c := controller.NewController(ctx, m, sm)
	return &Router{
		validator:  v,
		middleware: m,
		controller: c,
		engine:     echo.New(),
	}
}

// Setup implements IAppInstance.
func (r *Router) Setup(ctx context.Context) error {

	e := r.engine
	e.HideBanner = true
	e.HidePort = true
	e.Use(eMiddleware.CORSWithConfig(eMiddleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowHeaders: []string{"*"},
		AllowMethods: []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete, http.MethodOptions, http.MethodPatch},
	}))
	e.Validator = r.validator
	e.HTTPErrorHandler = r.middleware.ErrorHandling.Handler()

	if err := r.validator.Setup(ctx); err != nil {
		return err
	}
	if err := r.middleware.Setup(ctx); err != nil {
		return err
	}
	if err := r.controller.Setup(e); err != nil {
		return err
	}
	return nil
}

// Start implements IAppInstance.
func (r *Router) Start(ctx context.Context) error {
	return r.engine.Start(":8080")
}

// Stop implements IAppInstance.
func (r *Router) Stop(ctx context.Context) error {
	return r.engine.Shutdown(ctx)
}

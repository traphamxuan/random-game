package handler

import (
	"context"
	"game-random-api/internal/api/common"
	"game-random-api/internal/api/dto"
	"game-random-api/internal/api/service/authentication"
	"game-random-api/internal/api/validator"
	"game-random-api/package/logger"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/traphamxuan/gobs"
)

type Authentication struct {
	log  *logger.Logger
	auth *authentication.Authentication
}

var _ gobs.IService = (*Authentication)(nil)
var _ IHandler = (*Authentication)(nil)

func (a *Authentication) Init(ctx context.Context, sb *gobs.Component) error {
	sb.Deps = []gobs.IService{
		&logger.Logger{},
		&authentication.Authentication{},
	}
	onSetup := func(ctx context.Context, dependencies []gobs.IService, _ []gobs.CustomService) error {
		a.log = dependencies[0].(*logger.Logger)
		a.auth = dependencies[1].(*authentication.Authentication)
		return nil
	}
	sb.OnSetup = &onSetup
	return nil
}

// Route implements IHandler.
func (a *Authentication) SetRoute(ctx context.Context, e *echo.Group) error {
	g := e.Group("/auth")
	g.PUT("", a.Signin)
	g.DELETE("", a.Logout)
	g.GET("", a.Refresh)
	return nil
}

func (a *Authentication) Signin(c echo.Context) error {
	var creds dto.Credentials

	if err := validator.BindAndValidate(c, &creds); nil != err {
		a.log.Error(err)
		return err
	}

	resp, err := a.auth.Signin(c.Request().Context(), creds)
	if nil != err {
		return err
	}

	return c.JSON(http.StatusOK, resp)
}

func (a *Authentication) Refresh(c echo.Context) error {
	tokenStr, err := common.ExtractToken(c.Request().Header["Authorization"][0])
	if nil != err {
		return err
	}

	resp, err := a.auth.RefreshToken(c.Request().Context(), tokenStr)
	if nil != err {
		return err
	}

	return c.JSON(http.StatusOK, resp)
}

func (a *Authentication) Logout(c echo.Context) error {
	return c.JSON(http.StatusOK, nil)
}

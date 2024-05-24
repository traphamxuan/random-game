package controller

import (
	"context"

	"github.com/labstack/echo/v4"
	"github.com/traphamxuan/random-game/app/api/middleware"
	"github.com/traphamxuan/random-game/package/logger"
	servicemanager "github.com/traphamxuan/random-game/package/service_manager"
)

type User struct {
	mw  *middleware.Middleware
	log *logger.Logger
}

var _ IController = (*User)(nil)

func NewUser(ctx context.Context,
	mw *middleware.Middleware,
	sm *servicemanager.ServiceManager,
) *User {
	return &User{
		mw:  mw,
		log: servicemanager.GetServiceOrPanic[*logger.Logger](sm, "Logger", "User"),
	}
}

func (u *User) Setup(e *echo.Echo) error {
	g := e.Group("/api/user")
	g.GET("/", u.GetMany)
	g.Use(u.mw.Authentication.Handler())
	g.GET("/:id", u.GetDetail)
	g.POST("/", u.Create)
	g.PUT("/:id", u.Update)
	g.DELETE("/:id", u.Delete)

	return nil
}

func (u *User) GetMany(c echo.Context) error {
	return c.JSON(200, "GetMany")
}

func (u *User) GetDetail(c echo.Context) error {
	return c.JSON(200, "GetDetail")
}

func (u *User) Create(c echo.Context) error {
	return c.JSON(200, "Create")
}

func (u *User) Update(c echo.Context) error {
	return c.JSON(200, "Update")
}

func (u *User) Delete(c echo.Context) error {
	return c.JSON(200, "Delete")
}

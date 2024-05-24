package middleware

import (
	"context"

	"github.com/golang-jwt/jwt/v5"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/traphamxuan/random-game/package/config"
	servicemanager "github.com/traphamxuan/random-game/package/service_manager"
)

type Authentication struct {
	config    *config.Configuration
	jwtSecret string
}

func NewAuthentication(ctx context.Context, sm *servicemanager.ServiceManager) *Authentication {
	return &Authentication{
		config: servicemanager.GetServiceOrPanic[*config.Configuration](sm, "Config", "Authentication"),
	}
}

type JWTSecret struct {
	Secret string `env:"JWT_SECRET" envDefault:"mysecretjwt"`
}

func (a *Authentication) Setup(ctx context.Context) error {
	var cfg JWTSecret
	if err := a.config.ParseConfig(&cfg); err != nil {
		return err
	}
	a.jwtSecret = cfg.Secret
	return nil
}

func (a *Authentication) Handler() echo.MiddlewareFunc {
	config := echojwt.Config{
		NewClaimsFunc: func(c echo.Context) jwt.Claims {
			return new(jwt.RegisteredClaims)
		},
		SigningKey: []byte(a.jwtSecret),
	}
	return echojwt.WithConfig(config)
}

package middleware

import (
	"context"

	"game-random-api/package/config"
	"game-random-api/package/logger"

	"github.com/golang-jwt/jwt/v5"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/traphamxuan/gobs"
)

type Authentication struct {
	log    *logger.Logger
	config *JWTSecret
}

var _ gobs.IService = (*Authentication)(nil)

func (a *Authentication) Init(ctx context.Context, sb *gobs.Component) error {
	sb.Deps = []gobs.IService{
		&logger.Logger{},
		&config.Configuration{},
	}
	onSetup := func(ctx context.Context, dependencies []gobs.IService, _ []gobs.CustomService) error {
		a.log = dependencies[0].(*logger.Logger)
		config := dependencies[1].(*config.Configuration)
		var cfg JWTSecret
		if err := config.ParseConfig(&cfg); err != nil {
			return err
		}
		a.config = &cfg
		return nil
	}
	sb.OnSetup = &onSetup
	return nil
}

type JWTSecret struct {
	Secret string `env:"JWT_SECRET" envDefault:"mysecretjwt"`
}

func (a *Authentication) Handler() echo.MiddlewareFunc {
	config := echojwt.Config{
		NewClaimsFunc: func(c echo.Context) jwt.Claims {
			return new(jwt.RegisteredClaims)
		},
		SigningKey: []byte(a.config.Secret),
	}
	return echojwt.WithConfig(config)
}

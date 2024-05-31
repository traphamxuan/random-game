package middleware

import (
	"context"
	"game-random-api/package/logger"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/traphamxuan/gobs"
)

type MWLogger struct {
	log *logger.Logger
}

var _ gobs.IService = (*MWLogger)(nil)

// Init implements gobs.IService.
func (l *MWLogger) Init(ctx context.Context, sb *gobs.Component) error {
	sb.Deps = []gobs.IService{
		&logger.Logger{},
	}
	onSetup := func(ctx context.Context, dependencies []gobs.IService, _ []gobs.CustomService) error {
		l.log = dependencies[0].(*logger.Logger)
		return nil
	}
	sb.OnSetup = &onSetup
	return nil
}

func (l *MWLogger) Handler() echo.MiddlewareFunc {
	return middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogURI:     true,
		LogStatus:  true,
		LogMethod:  true,
		LogLatency: true,
		LogError:   true,
		LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
			l.log.Debugf("%s %s %d %s",
				v.Method,
				v.URI,
				v.Status,
				v.Latency.String(),
			)
			return nil
		},
	})
}

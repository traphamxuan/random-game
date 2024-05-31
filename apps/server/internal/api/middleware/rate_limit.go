package middleware

import (
	"context"
	"game-random-api/package/logger"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/traphamxuan/gobs"
	"golang.org/x/time/rate"
)

type RateLimit struct {
	log *logger.Logger
}

var _ gobs.IService = (*RateLimit)(nil)

// Init implements gobs.IService.
func (r *RateLimit) Init(ctx context.Context, sb *gobs.Component) error {
	sb.Deps = []gobs.IService{&logger.Logger{}}
	onSetup := func(ctx context.Context, dependencies []gobs.IService, _ []gobs.CustomService) error {
		r.log = dependencies[0].(*logger.Logger)
		return nil
	}
	sb.OnSetup = &onSetup
	return nil
}

func (r *RateLimit) Handler(limit float64, burst int, duration time.Duration) echo.MiddlewareFunc {
	config := middleware.RateLimiterConfig{
		Skipper: middleware.DefaultSkipper,
		Store: middleware.NewRateLimiterMemoryStoreWithConfig(
			middleware.RateLimiterMemoryStoreConfig{Rate: rate.Limit(limit), Burst: burst, ExpiresIn: duration},
		),
		IdentifierExtractor: func(ctx echo.Context) (string, error) {
			id := ctx.RealIP()
			return id, nil
		},
		ErrorHandler: func(context echo.Context, err error) error {
			return context.JSON(http.StatusForbidden, nil)
		},
		DenyHandler: func(context echo.Context, identifier string, err error) error {
			return context.JSON(http.StatusTooManyRequests, nil)
		},
	}

	return middleware.RateLimiterWithConfig(config)
}

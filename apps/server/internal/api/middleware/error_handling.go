package middleware

import (
	"context"
	"errors"
	"net"
	"net/http"

	"game-random-api/package/logger"

	"github.com/labstack/echo/v4"
	"github.com/traphamxuan/gobs"
)

type ErrorHandling struct {
	log *logger.Logger
}

var _ gobs.IService = (*ErrorHandling)(nil)

func (r *ErrorHandling) Init(ctx context.Context, sb *gobs.Component) error {
	sb.Deps = []gobs.IService{&logger.Logger{}}
	onSetup := func(ctx context.Context, dependencies []gobs.IService, _ []gobs.CustomService) error {
		r.log = dependencies[0].(*logger.Logger)
		return nil
	}
	sb.OnSetup = &onSetup
	return nil
}

func (e *ErrorHandling) Handler() echo.HTTPErrorHandler {
	return func(err error, c echo.Context) {
		if errors.Is(err, net.ErrClosed) {
			e.log.Debug("Connection is closed by client")
			return
		}

		if errors.Is(err, echo.ErrInternalServerError) {
			e.log.Errorf("Error: %v", err)
		}

		_err := c.JSON(http.StatusInternalServerError, err)
		if _err != nil {
			c.Echo().Logger.Error(_err)
		}
	}
}

package middleware

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/traphamxuan/random-game/package/logger"
	servicemanager "github.com/traphamxuan/random-game/package/service_manager"
)

type ErrorHandling struct {
	log *logger.Logger
}

func NewErrorHandling(ctx context.Context, sm *servicemanager.ServiceManager) *ErrorHandling {
	return &ErrorHandling{
		log: servicemanager.GetServiceOrPanic[*logger.Logger](sm, "Logger", "ErrorHandling"),
	}
}

func (e *ErrorHandling) Handler() echo.HTTPErrorHandler {
	return func(err error, c echo.Context) {
		if errors.Is(err, net.ErrClosed) {
			e.log.Debug("Connection is closed by client")
			return
		}

		e.log.Error(fmt.Sprintf("Error: %v", err))

		_err := c.JSON(http.StatusInternalServerError, err)
		if _err != nil {
			c.Echo().Logger.Error(_err)
		}
	}
}

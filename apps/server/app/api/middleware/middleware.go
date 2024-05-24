package middleware

import (
	"context"

	servicemanager "github.com/traphamxuan/random-game/package/service_manager"
)

type IMiddleware interface {
	Setup(ctx context.Context) error
}

type Middleware struct {
	Authentication *Authentication
	ErrorHandling  *ErrorHandling
	RateLimit      *RateLimit
}

func NewMiddleware(ctx context.Context, sm *servicemanager.ServiceManager) *Middleware {
	return &Middleware{
		Authentication: NewAuthentication(ctx, sm),
		ErrorHandling:  NewErrorHandling(ctx, sm),
		RateLimit:      NewRateLimit(ctx),
	}
}

func (m *Middleware) Setup(ctx context.Context) error {
	for _, mw := range []IMiddleware{m.Authentication} {
		if err := mw.Setup(ctx); err != nil {
			return err
		}
	}
	return nil
}

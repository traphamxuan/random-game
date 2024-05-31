package middleware

import (
	"context"

	"github.com/traphamxuan/gobs"
)

type Middleware struct {
	Authentication *Authentication
	ErrorHandling  *ErrorHandling
	RateLimit      *RateLimit
	Log            *MWLogger
}

var middlewares = []gobs.IService{
	&Authentication{},
	&ErrorHandling{},
	&RateLimit{},
	&MWLogger{},
}

var _ gobs.IService = (*Middleware)(nil)

func (m *Middleware) Init(ctx context.Context, sb *gobs.Component) error {
	sb.Deps = middlewares
	onSetup := func(ctx context.Context, dependencies []gobs.IService, _ []gobs.CustomService) error {
		middlewares = dependencies
		m.Authentication = dependencies[0].(*Authentication)
		m.ErrorHandling = dependencies[1].(*ErrorHandling)
		m.RateLimit = dependencies[2].(*RateLimit)
		m.Log = dependencies[3].(*MWLogger)
		return nil
	}
	sb.OnSetup = &onSetup
	return nil
}

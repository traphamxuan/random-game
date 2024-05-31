package service

import (
	"context"
	"game-random-api/internal/api/service/authentication"

	"github.com/traphamxuan/gobs"
)

type Service struct {
	Authentication *authentication.Authentication
	JwtToken       *authentication.JwtToken
}

var services = []gobs.IService{
	&authentication.Authentication{},
	&authentication.JwtToken{},
}

var _ gobs.IService = (*Service)(nil)

// Init implements gobs.IService.
func (s *Service) Init(ctx context.Context, sb *gobs.Component) error {
	sb.Deps = services
	onSetup := func(ctx context.Context, dependencies []gobs.IService, _ []gobs.CustomService) error {
		services = dependencies
		s.Authentication = dependencies[0].(*authentication.Authentication)
		s.JwtToken = dependencies[1].(*authentication.JwtToken)
		return nil
	}
	sb.OnSetup = &onSetup
	return nil
}

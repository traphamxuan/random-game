package validator

import (
	"context"

	"game-random-api/package/logger"
	vl "game-random-api/package/validator"

	"github.com/traphamxuan/gobs"
)

type Validator struct {
	*BaseValidator
	UserName *UserName
}

var validators = []gobs.IService{
	&logger.Logger{},
	&vl.CustomValidator{},
	&UserName{},
}

var _ gobs.IService = (*Validator)(nil)

func (cv *Validator) Init(ctx context.Context, sb *gobs.Component) error {
	sb.Deps = validators
	onSetup := func(ctx context.Context, dependencies []gobs.IService, _ []gobs.CustomService) error {
		cv.BaseValidator = &BaseValidator{
			log: dependencies[0].(*logger.Logger),
			v:   dependencies[1].(*vl.CustomValidator),
		}
		cv.UserName = dependencies[2].(*UserName)
		validators = dependencies
		return cv.Configure(ctx)
	}
	sb.OnSetup = &onSetup
	return nil
}

// Setup implements gobs.IService.
func (cv *Validator) Configure(ctx context.Context) error {
	for _, v := range validators {
		if v, ok := v.(IValidator); ok {
			cv.v.RegisterValidation(v.Tag(), v.Validate)
		}
	}
	return nil
}

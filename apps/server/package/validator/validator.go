package validator

import (
	"context"
	"game-random-api/package/logger"

	v "github.com/go-playground/validator/v10"
	"github.com/traphamxuan/gobs"
)

type CustomValidator struct {
	log       *logger.Logger
	validator *v.Validate
}

var _ gobs.IService = (*CustomValidator)(nil)

// Init implements gobs.IService.
func (cv *CustomValidator) Init(ctx context.Context, sb *gobs.Component) error {
	sb.Deps = []gobs.IService{
		&logger.Logger{},
	}
	onSetup := func(ctx context.Context, dependencies []gobs.IService, _ []gobs.CustomService) error {
		cv.log = dependencies[0].(*logger.Logger)
		return cv.Setup(ctx)
	}
	sb.OnSetup = &onSetup
	return nil
}

func (cv *CustomValidator) Setup(ctx context.Context) error {
	cv.validator = v.New()
	return nil
}

func (cv *CustomValidator) RegisterValidation(tag string, fn v.Func, callValidationEvenIfNull ...bool) error {
	return cv.validator.RegisterValidation(tag, fn, callValidationEvenIfNull...)
}

func (cv *CustomValidator) Validate(i interface{}) error {
	if err := cv.validator.Struct(i); err != nil {
		// Optionally, you could return the error to give each route more control over the status code
		return err
	}
	return nil
}

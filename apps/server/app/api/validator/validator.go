package validator

import (
	"context"

	v "github.com/go-playground/validator/v10"
	servicemanager "github.com/traphamxuan/random-game/package/service_manager"
)

type IValidator interface {
	Validate(i v.FieldLevel) bool
	Tag() string
}

type Validator struct {
	validators []IValidator
	validator  *v.Validate
}

var _ servicemanager.IService = (*Validator)(nil)

func NewValidator(ctx context.Context, sm *servicemanager.ServiceManager) *Validator {
	return &Validator{
		validators: []IValidator{
			NewUserName(ctx),
		},
	}
}

func (cv *Validator) Validate(i interface{}) error {
	if err := cv.validator.Struct(i); err != nil {
		// Optionally, you could return the error to give each route more control over the status code
		return err
	}
	return nil
}

// Setup implements servicemanager.IService.
func (cv *Validator) Setup(ctx context.Context) error {
	cv.validator = v.New()
	for _, v := range cv.validators {
		cv.validator.RegisterValidation(v.Tag(), v.Validate)
	}
	return nil
}

// Start implements servicemanager.IService.
func (cv *Validator) Start(ctx context.Context) error {
	return nil
}

// Stop implements servicemanager.IService.
func (cv *Validator) Stop(ctx context.Context) error {
	return nil
}

package validator

import (
	"context"

	"game-random-api/package/logger"
	valid "game-random-api/package/validator"
	"game-random-api/utils"

	v "github.com/go-playground/validator/v10"
	"github.com/traphamxuan/gobs"
)

type UserName BaseValidator

// Init implements IValidator.
func (uv *UserName) Init(ctx context.Context, sb *gobs.Component) error {
	sb.Deps = []gobs.IService{
		&logger.Logger{},
		&valid.CustomValidator{},
	}
	onSetup := func(ctx context.Context, dependencies []gobs.IService, _ []gobs.CustomService) error {
		uv.log = dependencies[0].(*logger.Logger)
		uv.v = dependencies[1].(*valid.CustomValidator)
		return nil
	}
	sb.OnSetup = &onSetup
	return nil
}

var _ IValidator = (*UserName)(nil)

func (uv *UserName) Tag() string {
	return "username"
}

func (uv *UserName) Validate(fl v.FieldLevel) bool {
	return utils.UserNameRegex.MatchString(fl.Field().String())
}

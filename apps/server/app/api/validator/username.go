package validator

import (
	"context"

	v "github.com/go-playground/validator/v10"
	"github.com/traphamxuan/random-game/app/utils"
)

type UserName struct {
}

var _ IValidator = (*UserName)(nil)

func NewUserName(ctx context.Context) *UserName {
	return &UserName{}
}

func (uv *UserName) Tag() string {
	return "username"
}

func (uv *UserName) Validate(fl v.FieldLevel) bool {
	return utils.UserNameRegex.MatchString(fl.Field().String())
}

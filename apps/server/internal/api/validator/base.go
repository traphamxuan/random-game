package validator

import (
	"game-random-api/package/logger"
	"game-random-api/package/validator"

	v "github.com/go-playground/validator/v10"
)

type IValidator interface {
	Validate(i v.FieldLevel) bool
	Tag() string
}
type BaseValidator struct {
	log *logger.Logger
	v   *validator.CustomValidator
}

package common

import (
	"errors"
	"game-random-api/internal/orm/schema"
)

type iContext interface {
	Get(string) (any, bool)
}

func GetUserCtx(ctx iContext) (ctxUser schema.User, err error) {
	if ctxVal, ok := ctx.Get("ctxUser"); ok {
		if ctxUser, ok = ctxVal.(schema.User); ok {
			return ctxUser, nil
		}
	}
	return ctxUser, errors.New("invalid user context")
}

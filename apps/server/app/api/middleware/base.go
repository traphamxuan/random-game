package middleware

import (
	"github.com/traphamxuan/random-game/package/cache"
	"github.com/traphamxuan/random-game/package/logger"
	"github.com/traphamxuan/random-game/package/orm"
)

type BaseMiddleware struct {
	log *logger.Logger
	db  *orm.Orm
	rdb *cache.Redis
}

func NewBaseMiddleware(log *logger.Logger, db *orm.Orm, rdb *cache.Redis) *BaseMiddleware {
	return &BaseMiddleware{
		log: log,
		db:  db,
		rdb: rdb,
	}
}

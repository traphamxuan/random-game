package orm

import (
	"context"
	"game-random-api/package/logger"

	"gorm.io/gorm"
)

type SortBy string

const (
	SortByASC  SortBy = "ASC"
	SortByDESC SortBy = "DESC"
)

type BaseORM struct {
	db  *gorm.DB
	log *logger.Logger
}

func (b *BaseORM) Delete(value interface{}, conds ...interface{}) *gorm.DB {
	return b.db.Delete(value, conds...)
}

func (o *BaseORM) StartTransaction(c context.Context, trx_func func(*Orm) error) error {
	return o.db.WithContext(c).
		Transaction(func(tx *gorm.DB) error {
			bOrm := &BaseORM{
				db:  tx,
				log: o.log,
			}
			orm := ormFromBase(bOrm)
			return trx_func(orm)
		})
}

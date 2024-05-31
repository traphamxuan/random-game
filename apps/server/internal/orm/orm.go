package orm

import (
	"context"
	"game-random-api/package/db"
	"game-random-api/package/logger"

	"github.com/traphamxuan/gobs"
)

type Orm struct {
	*BaseORM
	User *User
}

var _ gobs.IService = (*Orm)(nil)

func (o *Orm) Init(ctx context.Context, sb *gobs.Component) error {
	sb.Deps = []gobs.IService{
		&logger.Logger{},
		&db.Gorm{},
	}
	onSetup := func(ctx context.Context, dependencies []gobs.IService, _ []gobs.CustomService) error {
		o.BaseORM = &BaseORM{
			log: dependencies[0].(*logger.Logger),
			db:  dependencies[1].(*db.Gorm).Db,
		}
		return nil
	}
	sb.OnSetup = &onSetup
	return nil
}

func NewOrm(ctx context.Context) *Orm {
	bo := &BaseORM{}
	return ormFromBase(bo)
}

func ormFromBase(bOrm *BaseORM) *Orm {
	o := &Orm{
		BaseORM: bOrm,
		User:    (*User)(bOrm),
	}
	return o
}

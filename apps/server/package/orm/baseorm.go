package orm

import (
	"github.com/traphamxuan/random-game/package/logger"
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

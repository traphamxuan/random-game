package orm

type Models struct {
	User User
}

func NewModels(bOrm *BaseORM) Models {
	return Models{
		User: User{bOrm},
	}
}

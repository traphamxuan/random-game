package orm

import (
	"context"
	"errors"

	"github.com/fatih/structs"
	"github.com/traphamxuan/random-game/common"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type User struct {
	*BaseORM
}

type UserTable struct {
	ID       string `gorm:"column:id;type:UUID;NOT NULL;primaryKey;default:uuid_generate_v4()"`
	Email    string `gorm:"column:email;type:VARCHAR(255);NOT NULL;unique"`
	Password string `gorm:"column:password;type:VARCHAR(255);NOT NULL"`
	Name     string `gorm:"column:name;type:VARCHAR(255);NULL"`
}

func (UserTable) TableName() string {
	return "users"
}

type QueryUserFilter struct {
	ID    *string `structs:"id,omitempty"`
	Email *string `structs:"email,omitempty"`
	Name  *string `structs:"name,omitempty"`
}

func (o *QueryUserFilter) Map() map[string]interface{} {
	return structs.Map(o)
}

func (o *User) GetMany(ctx context.Context, filter QueryUserFilter, page common.Page) (
	[]UserTable, error,
) {
	var Users []UserTable
	tx := o.db.WithContext(ctx).Model(&UserTable{}).
		Where(filter.Map()).
		Offset(page.Offset).
		Limit(page.Limit)
	if page.OrderBy != "" {
		tx = tx.Order(page.OrderBy + " " + string(page.SortBy))
	}
	if err := tx.Find(&Users).Error; err != nil {
		return nil, err
	}
	return Users, nil
}

func (o *User) GetOne(
	ctx context.Context,
	filter QueryUserFilter,
) (*UserTable, error) {
	var user UserTable
	if err := o.db.WithContext(ctx).Model(&UserTable{}).
		Where(filter.Map()).
		First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (o *User) Count(
	ctx context.Context,
	filter QueryUserFilter,
) (int64, error) {
	var count int64
	if err := o.db.WithContext(ctx).Model(&UserTable{}).
		Where(filter.Map()).
		Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

type UserCreateInput struct {
	Email    string `structs:"email"`
	Name     string `structs:"name"`
	Password string `structs:"password"`
}

func (o *User) BulkCreate(
	ctx context.Context,
	data []UserCreateInput,
	playID string,
) ([]UserTable, error) {
	var Users []UserTable
	for _, input := range data {
		Users = append(Users, UserTable{
			Email:    input.Email,
			Name:     input.Name,
			Password: input.Password,
		})
	}
	if err := o.db.WithContext(ctx).Model(&UserTable{}).CreateInBatches(Users, 100).Error; err != nil {
		return nil, err
	}
	return Users, nil
}

type UserUpsertInput struct {
	ID *string
	UserCreateInput
}

func (o *User) BulkUpsert(
	ctx context.Context,
	inputs []UserUpsertInput,
) ([]UserTable, error) {
	Users := make([]UserTable, len(inputs))
	for idx, input := range inputs {
		var id string
		if input.ID != nil {
			id = *input.ID
		}
		Users[idx] = UserTable{
			ID:       id,
			Name:     input.Name,
			Email:    input.Email,
			Password: input.Password,
		}
	}
	if err := o.db.WithContext(ctx).Clauses(clause.OnConflict{
		Columns: []clause.Column{
			{Name: "id"},
			{Name: "source_id"},
			{Name: "target_id"},
		},
		DoUpdates: clause.AssignmentColumns([]string{
			"quantity",
		}),
	}).CreateInBatches(&Users, 100).Error; err != nil {
		return nil, err
	}
	return Users, nil
}

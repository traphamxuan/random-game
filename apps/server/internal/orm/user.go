package orm

import (
	"context"
	"game-random-api/internal/orm/schema"
	"game-random-api/utils"

	"github.com/fatih/structs"
	"gorm.io/gorm/clause"
)

type User BaseORM

type QueryUserFilter struct {
	ID       *string `structs:"id,omitempty"`
	Name     *string `structs:"name,omitempty"`
	Avatar   *string `structs:"avatar,omitempty"`
	Email    *string `structs:"email,omitempty"`
	Password *string `structs:"password,omitempty"`
}

func (o *QueryUserFilter) Map() map[string]interface{} {
	return structs.Map(o)
}

type UserCreateInput struct {
	Name     *string
	Avatar   *string
	Email    string
	Password string
}

type UserUpdateInput struct {
	Name   *string `structs:"name,omitempty"`
	Avatar *string `structs:"avatar,omitempty"`
}

func (o *UserUpdateInput) Map() map[string]interface{} {
	return structs.Map(o)
}

type UserUpsertInput struct {
	ID       *string
	Name     *string
	Avatar   *string
	Email    string
	Password string
}

func (o *User) GetOne(
	ctx context.Context,
	filter QueryUserFilter,
) (*schema.User, error) {
	var product schema.User
	if err := o.db.WithContext(ctx).
		Model(&schema.User{}).
		Where(filter.Map()).
		First(&product).Error; err != nil {
		return nil, err
	}
	return &product, nil
}

func (o *User) GetMany(
	ctx context.Context,
	filter QueryUserFilter,
	page utils.Page,
	relations ...string,
) ([]schema.User, error) {
	var products []schema.User

	tx := o.db.WithContext(ctx).Model(&schema.User{})
	for _, relation := range relations {
		tx = tx.Preload(relation)
	}
	tx = tx.Where(filter.Map()).
		Offset(page.Offset).
		Limit(page.Limit)
	if page.OrderBy != "" {
		tx = tx.Order(page.OrderBy + " " + string(page.SortBy))
	}
	if err := tx.Find(&products).Error; err != nil {
		return nil, err
	}
	return products, nil
}

func (o *User) Count(
	ctx context.Context,
	filter QueryUserFilter,
) (int64, error) {
	var count int64
	if err := o.db.WithContext(ctx).Model(&schema.User{}).
		Where(filter.Map()).
		Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

func (o *User) CreateOne(
	ctx context.Context,
	data UserCreateInput,
) (*schema.User, error) {
	product := schema.User{
		Name:     data.Name,
		Avatar:   data.Avatar,
		Email:    data.Email,
		Password: data.Password,
	}
	if err := o.db.WithContext(ctx).
		Model(&schema.User{}).
		Create(&product).Error; err != nil {
		return nil, err
	}
	return &product, nil
}

func (o *User) BulkCreate(
	ctx context.Context,
	data []UserCreateInput,
) ([]schema.User, error) {
	var products []schema.User
	for _, input := range data {
		products = append(products, schema.User{
			Name:     input.Name,
			Avatar:   input.Avatar,
			Email:    input.Email,
			Password: input.Password,
		})
	}
	if err := o.db.WithContext(ctx).Model(&schema.User{}).
		CreateInBatches(products, 100).
		Error; err != nil {
		return nil, err
	}
	return products, nil
}

func (o *User) Update(
	ctx context.Context,
	filter QueryUserFilter,
	input UserUpdateInput,
	isReturned bool,
) ([]schema.User, error) {
	results := []schema.User{}
	tx := o.db.WithContext(ctx).Model(&results)
	if isReturned {
		tx = tx.Clauses(clause.Returning{})
	}
	tx = tx.Where(filter.Map()).Updates(input.Map())
	if err := tx.Error; err != nil {
		return nil, err
	}
	return results, nil
}

func (o *User) BulkUpsert(
	ctx context.Context,
	inputs []UserUpsertInput,
) ([]schema.User, error) {
	products := make([]schema.User, len(inputs))
	for idx, input := range inputs {
		var id string
		if input.ID != nil {
			id = *input.ID
		}
		products[idx] = schema.User{
			ID:       id,
			Name:     input.Name,
			Avatar:   input.Avatar,
			Password: input.Password,
			Email:    input.Email,
		}
	}
	if err := o.db.WithContext(ctx).Clauses(
		clause.OnConflict{
			Columns: []clause.Column{
				{Name: "id"},
			},
			DoUpdates: clause.AssignmentColumns([]string{
				"name",
				"avatar",
				"email",
				"password",
			}),
		},
		clause.OnConflict{
			Columns: []clause.Column{
				{Name: "email"},
			},
			DoUpdates: clause.AssignmentColumns([]string{
				"avatar",
				"name",
				"password",
			}),
		}).CreateInBatches(&products, 100).Error; err != nil {
		return nil, err
	}
	return products, nil
}

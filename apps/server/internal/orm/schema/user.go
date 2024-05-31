package schema

import (
	"time"
)

type User struct {
	ID        string    `gorm:"column:id;type:UUID;NOT NULL;primaryKey"`
	Name      *string   `gorm:"column:name;type:varchar(255);NULL"`
	Avatar    *string   `gorm:"column:avatar;type:text;NULL"`
	Email     string    `gorm:"column:email;type:varchar(255);NOT NULL"`
	Password  string    `gorm:"column:password;type:text;NOT NULL"`
	UpdatedAt time.Time `gorm:"column:updatedAt;type:TIMESTAMP WITHOUT TIME ZONE;NOT NULL"`
	// Token     string    `gorm:"column:token;type:text;NOT NULL;uniqueIndex:user_uk_token"`
}

func (User) TableName() string {
	return "users"
}

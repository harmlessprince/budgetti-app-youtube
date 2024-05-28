package models

import "gorm.io/gorm"

type UserModel struct {
	gorm.Model
	FirstName *string `gorm:"type:varchar(200)"`
	LastName  *string `gorm:"type:varchar(200)"`
	Email     string  `gorm:"type:varchar(200);not null;unique"`
	Gender    *string `gorm:"type:varchar(50)"`
	Password  string  `gorm:"type:varchar(200);not null"`
}

func (receiver UserModel) TableName() string {
	return "users"
}

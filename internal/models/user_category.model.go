package models

type UserCategoryModel struct {
	UserID     uint `gorm:"primaryKey;column:user_id" json:"user_id"`
	CategoryID uint `gorm:"primaryKey;column:category_id" json:"category_id"`
}

func (UserCategoryModel) TableName() string {
	return "user_categories"
}

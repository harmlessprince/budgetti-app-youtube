package common

import "gorm.io/gorm"

func WhereUserIDScope(UserID uint) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("user_id = ?", UserID)
	}
}

func WhereUserIDWithTableNameScope(UserID uint, tableName string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where(tableName+".user_id = ?", UserID) // select table_name.user_id = 1
	}
}

func LoadCategories(db *gorm.DB) *gorm.DB {
	return db.
		Joins("LEFT JOIN budget_categories bc ON bc.budget_model_id = budgets.id").
		Joins("LEFT JOIN categories c ON c.id = bc.category_model_id")
}

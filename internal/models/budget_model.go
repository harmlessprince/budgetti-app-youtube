package models

import "time"

type BudgetModel struct {
	BaseModel
	Title       string           `gorm:"index;type:varchar(255);not null" json:"title"`
	Slug        string           `gorm:"index;type:varchar(255);not null;uniqueIndex:unique_user_id_slug_year_month" json:"slug"`
	Description *string          `gorm:"type:text" json:"description"`
	UserID      uint             `gorm:"not null;column:user_id;uniqueIndex:unique_user_id_slug_year_month" json:"user_id"`
	Amount      float64          `gorm:"type:decimal(10,2);not null" json:"amount"`
	Categories  []*CategoryModel `gorm:"constraint:OnDelete:CASCADE;many2many:budget_categories;" json:"categories"`
	Date        time.Time        `gorm:"type:datetime;not null" json:"date"`
	Month       uint             `gorm:"type:TINYINT UNSIGNED;not null;index:idx_month_year;uniqueIndex:unique_user_id_slug_year_month" json:"month"`
	Year        uint16           `gorm:"type:INT UNSIGNED;not null;index:idx_month_year;uniqueIndex:unique_user_id_slug_year_month" json:"yea"`
}

// slug, year, month, user_id

// uniqueIndex unique_user_id_slug_year_month
// unique (user_id, slug, month, year)
func (BudgetModel) TableName() string {
	return "budgets"
}

//Categories  []*CategoryModel `gorm:"constraint:OnDelete:CASCADE;many2many:budget_categories;joinForeignKey:BudgetID;joinReferences:CategoryID" json:"categories"`

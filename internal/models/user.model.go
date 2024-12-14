package models

type UserModel struct {
	BaseModel
	FirstName  *string         `gorm:"type:varchar(200)" json:"first_name"`
	LastName   *string         `gorm:"type:varchar(200)" json:"last_name"`
	Email      string          `gorm:"type:varchar(200);not null;unique" json:"email"`
	Gender     *string         `gorm:"type:varchar(50)" json:"gender"`
	Password   string          `gorm:"type:varchar(200);not null" json:"-"`
	Categories []CategoryModel `gorm:"many2many:user_categories;" json:"categories"`
	Budgets    []BudgetModel   `gorm:"foreignKey:UserID" json:"-"`
}

func (receiver UserModel) TableName() string {
	return "users"
}

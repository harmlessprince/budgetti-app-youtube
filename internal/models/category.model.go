package models

type CategoryModel struct {
	BaseModel
	Name     string `gorm:"unique;type:varchar(200);not null" json:"name"`
	Slug     string `gorm:"type:varchar(200);unique;not null" json:"slug"`
	IsCustom bool   `gorm:"type:bool;default:false" json:"is_custom"`
}

func (CategoryModel) TableName() string {
	return "categories"
}

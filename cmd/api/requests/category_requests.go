package requests

type CreateCategoryRequest struct {
	Name     string `json:"name" validate:"required"`
	IsCustom bool   `default:"true" json:"is_custom"`
}

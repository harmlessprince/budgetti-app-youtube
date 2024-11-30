package services

import (
	"errors"
	"github.com/harmlessprince/bougette-backend/cmd/api/requests"
	"github.com/harmlessprince/bougette-backend/common"
	"github.com/harmlessprince/bougette-backend/internal/app_errors"
	"github.com/harmlessprince/bougette-backend/internal/models"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"log"
	"strings"
)

type CategoryService struct {
	DB *gorm.DB
}

func NewCategoryService(db *gorm.DB) *CategoryService {
	return &CategoryService{DB: db}
}
func (c CategoryService) List(categories []*models.CategoryModel, pagination *common.Pagination) (*common.Pagination, error) {
	c.DB.Scopes(pagination.Paginate()).Find(&categories)
	pagination.Items = categories
	return pagination, nil
}

func (c CategoryService) Create(data *requests.CreateCategoryRequest) (*models.CategoryModel, error) {
	slug := strings.ToLower(data.Name)
	slug = strings.Replace(slug, " ", "_", -1)
	categoryCreated := &models.CategoryModel{
		Slug:     slug,
		Name:     data.Name,
		IsCustom: data.IsCustom,
	}
	result := c.DB.Where(models.CategoryModel{Slug: slug, Name: data.Name}).FirstOrCreate(categoryCreated)
	if result.Error != nil {
		log.Println("failed to create", slug)
		log.Println(result.Error.Error())
		if errors.Is(result.Error, gorm.ErrDuplicatedKey) {
			return categoryCreated, nil
		}
		return nil, errors.New("failed to create category")
	}
	return categoryCreated, nil
}

func (c CategoryService) GetById(id uint) (*models.CategoryModel, error) {
	var category models.CategoryModel
	result := c.DB.First(&category, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, app_errors.NewNotFoundError("Category not found")
		}
		return nil, errors.New("failed to fetch category")
	}
	return &category, nil
}

func (c CategoryService) DeleteById(id uint) error {
	var category *models.CategoryModel
	category, err := c.GetById(id)
	if err != nil {
		return err
	}
	c.DB.Delete(&category)
	return nil
}

func (c CategoryService) AssociateUserToCategories(user *models.UserModel, categories []*models.CategoryModel) error {
	if user != nil && categories != nil && len(categories) > 0 {
		var userCategories []*models.UserCategoryModel
		for _, category := range categories {
			userCategories = append(userCategories, &models.UserCategoryModel{
				UserID:     user.ID,
				CategoryID: category.ID,
			})
		}
		err := c.DB.Clauses(clause.OnConflict{DoNothing: true}).Create(userCategories)
		if err != nil {
			return err.Error
		}
	}
	return nil
}

func (c CategoryService) GetMultipleCategories(categoryIds []uint) ([]*models.CategoryModel, error) {
	var categories []*models.CategoryModel
	result := c.DB.Where("id IN ?", categoryIds).Find(&categories)
	if result.Error != nil {
		return nil, result.Error
	}
	return categories, nil
}

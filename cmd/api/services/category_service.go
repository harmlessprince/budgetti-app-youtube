package services

import (
	"errors"
	"github.com/harmlessprince/bougette-backend/cmd/api/requests"
	"github.com/harmlessprince/bougette-backend/internal/app_errors"
	"github.com/harmlessprince/bougette-backend/internal/models"
	"gorm.io/gorm"
	"log"
	"strings"
)

type CategoryService struct {
	DB *gorm.DB
}

func NewCategoryService(db *gorm.DB) *CategoryService {
	return &CategoryService{DB: db}
}
func (c CategoryService) List() ([]*models.CategoryModel, error) {
	var categories []*models.CategoryModel
	result := c.DB.Find(&categories)
	if result.Error != nil {
		return nil, errors.New("failed to fetch all categories")
	}
	return categories, nil
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

//func (c CategoryService) getById(id int)  ()  {
//
//}
//
//func (c CategoryService) delete(id int)  ()  {
//
//}
//
//func (c CategoryService) getById(id int)  ()  {
//
//}

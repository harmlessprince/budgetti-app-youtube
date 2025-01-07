package handlers

import (
	"errors"
	"fmt"
	"github.com/harmlessprince/bougette-backend/cmd/api/requests"
	"github.com/harmlessprince/bougette-backend/cmd/api/services"
	"github.com/harmlessprince/bougette-backend/common"
	"github.com/harmlessprince/bougette-backend/internal/models"
	"github.com/labstack/echo/v4"
)

func (h *Handler) ListCategories(c echo.Context) error {
	var categories []*models.CategoryModel
	categoryService := services.NewCategoryService(h.DB)
	paginator := common.NewPaginator(categories, c.Request(), h.DB)
	paginatedCategory, err := categoryService.List(categories, paginator)
	if err != nil {
		return common.SendInternalServerErrorResponse(c, err.Error())
	}
	return common.SendSuccessResponse(c, "ok", paginatedCategory)
}

func (h *Handler) CreateCategory(c echo.Context) error {
	_, ok := c.Get("user").(models.UserModel)
	if !ok {
		return common.SendInternalServerErrorResponse(c, "User authentication failed")
	}
	// bind request body
	payload := new(requests.CreateCategoryRequest)
	if err := h.BindBodyRequest(c, payload); err != nil {
		return common.SendBadRequestResponse(c, err.Error())
	}

	// validation
	validationErrors := h.ValidateBodyRequest(c, *payload)

	if validationErrors != nil {
		return common.SendFailedValidationResponse(c, validationErrors)
	}

	categoryService := services.NewCategoryService(h.DB)

	category, err := categoryService.Create(payload)
	if err != nil {
		return common.SendInternalServerErrorResponse(c, err.Error())
	}
	return common.SendSuccessResponse(c, "Category created successfully", category)
}

func (h *Handler) DeleteCategory(c echo.Context) error {
	_, ok := c.Get("user").(models.UserModel)
	if !ok {
		return errors.New("User authentication failed")
	}
	var categoryId requests.IDParamRequest
	err := (&echo.DefaultBinder{}).BindPathParams(c, &categoryId)
	if err != nil {
		return common.SendBadRequestResponse(c, err.Error())
	}
	categoryService := services.NewCategoryService(h.DB)
	err = categoryService.DeleteById(categoryId.ID)
	if err != nil {
		return err
	}
	return common.SendSuccessResponse(c, "Category Deleted", nil)
}

func (h *Handler) AssociateUserToCategories(c echo.Context) error {
	user, _ := c.Get("user").(models.UserModel)
	payload := new(requests.AssociateUserToCategoryRequest)
	if err := h.BindBodyRequest(c, payload); err != nil {
		return common.SendBadRequestResponse(c, err.Error())
	}

	// validation
	validationErrors := h.ValidateBodyRequest(c, *payload)

	if validationErrors != nil {
		return common.SendFailedValidationResponse(c, validationErrors)
	}

	categoryService := services.NewCategoryService(h.DB)

	categories, err := categoryService.GetMultipleCategories(payload.Categories)
	if err != nil {
		return common.SendInternalServerErrorResponse(c, err.Error())
	}
	err = categoryService.AssociateUserToCategories(&user, categories)
	if err != nil {
		return common.SendInternalServerErrorResponse(c, "An error occurred while associating user to categories")
	}
	totalCategories := len(categories)
	return common.SendSuccessResponse(c, fmt.Sprintf("%d Categories Associated", totalCategories), nil)
}

func (h *Handler) ListUserCategories(c echo.Context) error {
	user, _ := c.Get("user").(models.UserModel)
	var categories []*models.CategoryModel
	query := h.DB.Model(models.CategoryModel{})
	query = query.InnerJoins("INNER JOIN user_categories ON categories.id = user_categories.category_id")
	query = query.Where("user_categories.user_id = ?", user.ID)

	categoryService := services.NewCategoryService(query)
	paginator := common.NewPaginator(categories, c.Request(), query)
	paginatedCategory, err := categoryService.List(categories, paginator)
	if err != nil {
		return common.SendInternalServerErrorResponse(c, err.Error())
	}
	return common.SendSuccessResponse(c, "ok", paginatedCategory)
}

func (h *Handler) CreateCustomUserCategory(c echo.Context) error {
	user, ok := c.Get("user").(models.UserModel)
	if !ok {
		return common.SendInternalServerErrorResponse(c, "User authentication failed")
	}
	// bind request body
	payload := new(requests.CreateCategoryRequest)
	if err := h.BindBodyRequest(c, payload); err != nil {
		return common.SendBadRequestResponse(c, err.Error())
	}

	// validation
	validationErrors := h.ValidateBodyRequest(c, *payload)

	if validationErrors != nil {
		return common.SendFailedValidationResponse(c, validationErrors)
	}

	categoryService := services.NewCategoryService(h.DB)

	category, err := categoryService.Create(payload)
	if err != nil {
		return common.SendInternalServerErrorResponse(c, err.Error())
	}
	categories := []*models.CategoryModel{category}
	err = categoryService.AssociateUserToCategories(&user, categories)
	if err != nil {
		return common.SendInternalServerErrorResponse(c, "An error occurred while associating user to categories")
	}
	return common.SendSuccessResponse(c, "Category created successfully", category)
}

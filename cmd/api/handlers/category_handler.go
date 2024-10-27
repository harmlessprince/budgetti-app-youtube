package handlers

import (
	"errors"
	"github.com/harmlessprince/bougette-backend/cmd/api/requests"
	"github.com/harmlessprince/bougette-backend/cmd/api/services"
	"github.com/harmlessprince/bougette-backend/common"
	"github.com/harmlessprince/bougette-backend/internal/app_errors"
	"github.com/harmlessprince/bougette-backend/internal/models"
	"github.com/labstack/echo/v4"
	"log"
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
		return common.SendInternalServerErrorResponse(c, "User authentication failed")
	}
	var categoryId requests.IDParamRequest
	err := (&echo.DefaultBinder{}).BindPathParams(c, &categoryId)
	if err != nil {
		return common.SendBadRequestResponse(c, err.Error())
	}
	log.Println("DeleteCategory", "categoryId", categoryId)
	categoryService := services.NewCategoryService(h.DB)
	err = categoryService.DeleteById(categoryId.ID)
	if err != nil {
		if errors.Is(err, app_errors.NewNotFoundError(err.Error())) {
			return common.SendNotFoundResponse(c, err.Error())
		}
		return common.SendBadRequestResponse(c, err.Error())
	}
	return common.SendSuccessResponse(c, "Category Deleted", nil)
}

//func (h *Handler) UpdateCategory(c echo.Context) error {
//	_, ok := c.Get("user").(models.UserModel)
//}

package handlers

import (
	"errors"
	"github.com/harmlessprince/bougette-backend/cmd/api/requests"
	"github.com/harmlessprince/bougette-backend/cmd/api/services"
	"github.com/harmlessprince/bougette-backend/common"
	"github.com/harmlessprince/bougette-backend/internal/models"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func (h *Handler) ListBudget(c echo.Context) error {
	user, _ := c.Get("user").(models.UserModel)

	var budgets []*models.BudgetModel
	budgetService := services.NewBudgetService(h.DB)
	query := h.DB.Preload("Categories").Scopes(common.WhereUserIDScope(user.ID))
	paginator := common.NewPaginator(budgets, c.Request(), query)
	paginatedBudget, err := budgetService.List(query, budgets, paginator)
	if err != nil {
		return common.SendInternalServerErrorResponse(c, err.Error())
	}
	return common.SendSuccessResponse(c, "ok", paginatedBudget)
}

func (h *Handler) CreateBudget(c echo.Context) error {
	user, _ := c.Get("user").(models.UserModel)

	// bind request body
	payload := new(requests.StoreBudgetRequest)
	if err := h.BindBodyRequest(c, payload); err != nil {
		return common.SendBadRequestResponse(c, err.Error())
	}

	// validation
	validationErrors := h.ValidateBodyRequest(c, *payload)

	if validationErrors != nil {
		return common.SendFailedValidationResponse(c, validationErrors)
	}

	budgetService := services.NewBudgetService(h.DB)
	categoryService := services.NewCategoryService(h.DB)

	categories, err := categoryService.GetMultipleCategories(payload.Categories)

	if err != nil {
		c.Logger().Error(err)
		return common.SendInternalServerErrorResponse(c, "Budget could not be created")
	}

	createdBudget := &models.BudgetModel{}

	err = h.DB.Transaction(func(tx *gorm.DB) error {
		budgetService.DB = tx
		categoryService.DB = tx
		createdBudget, err = budgetService.Create(payload, user.ID)
		if err != nil {
			c.Logger().Error(err)
			return common.SendInternalServerErrorResponse(c, "Budget could not be created, try again later")
		}

		err = tx.Model(createdBudget).Association("Categories").Replace(categories)

		if err != nil {
			c.Logger().Error(err)
			return common.SendInternalServerErrorResponse(c, "Budget could not be created")
		}

		return nil
	})

	if err != nil {
		return common.SendInternalServerErrorResponse(c, err.Error())
	}

	createdBudget.Categories = categories

	return common.SendSuccessResponse(c, "Budget created successfully", createdBudget)
}

func (h *Handler) UpdateBudget(c echo.Context) error {
	user, _ := c.Get("user").(models.UserModel)

	var budgetID requests.IDParamRequest
	err := (&echo.DefaultBinder{}).BindPathParams(c, &budgetID)
	if err != nil {
		return common.SendBadRequestResponse(c, err.Error())
	}

	// bind request body
	payload := new(requests.UpdateBudgetRequest)
	if err := h.BindBodyRequest(c, payload); err != nil {
		return common.SendBadRequestResponse(c, err.Error())
	}

	// validation
	validationErrors := h.ValidateBodyRequest(c, *payload)

	if validationErrors != nil {
		return common.SendFailedValidationResponse(c, validationErrors)
	}

	budgetService := services.NewBudgetService(h.DB)
	categoryService := services.NewCategoryService(h.DB)

	budget, err := budgetService.GetById(budgetID.ID)

	if err != nil {
		return err
	}
	if user.ID != budget.UserID {
		return common.SendNotFoundResponse(c, "Budget not found")
	}

	var categories []*models.CategoryModel

	if payload.Categories != nil {
		categories, _ = categoryService.GetMultipleCategories(payload.Categories)
	}

	err = h.DB.Transaction(func(tx *gorm.DB) error {
		budgetService.DB = tx
		categoryService.DB = tx
		updatedBudget, err := budgetService.Update(budget, payload, budgetID.ID)
		if err != nil {
			return common.SendBadRequestResponse(c, err.Error())
		}
		if len(categories) > 0 {
			err = budgetService.DB.Model(updatedBudget).Association("Categories").Replace(categories)
			if err != nil {
				return common.SendInternalServerErrorResponse(c, "Budget could not be updated")
			}
			updatedBudget.Categories = categories
		}
		budget = updatedBudget
		return nil
	})

	return common.SendSuccessResponse(c, "Budget updated successfully", budget)
}

func (h *Handler) DeleteBudget(c echo.Context) error {
	user, _ := c.Get("user").(models.UserModel)
	var budgetID requests.IDParamRequest
	err := (&echo.DefaultBinder{}).BindPathParams(c, &budgetID)
	if err != nil {
		return common.SendBadRequestResponse(c, err.Error())
	}
	budgetService := services.NewBudgetService(h.DB)
	budget, _ := budgetService.GetById(budgetID.ID)
	if budget == nil {
		return common.SendNotFoundResponse(c, "Budget not found")
	}
	if user.ID != budget.UserID {
		return common.SendNotFoundResponse(c, "Budget not found")
	}
	query := h.DB.Scopes(common.WhereUserIDScope(user.ID))
	err = h.DB.Model(&budget).Association("Categories").Clear()
	if err != nil {
		return errors.New("Budget could not be deleted")
	}
	query.Delete(&models.BudgetModel{}, budget.ID)
	return common.SendSuccessResponse(c, "Budget Deleted Successfully ", nil)
}

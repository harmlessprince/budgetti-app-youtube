package services

import (
	"errors"
	"fmt"
	"github.com/harmlessprince/bougette-backend/cmd/api/requests"
	"github.com/harmlessprince/bougette-backend/common"
	"github.com/harmlessprince/bougette-backend/internal/app_errors"
	"github.com/harmlessprince/bougette-backend/internal/models"
	"github.com/labstack/gommon/log"
	"gorm.io/gorm"
	"strings"
	"time"
)

type BudgetService struct {
	DB *gorm.DB
}

func NewBudgetService(db *gorm.DB) *BudgetService {
	return &BudgetService{DB: db}
}

func (b *BudgetService) Create(payload *requests.StoreBudgetRequest, UserID uint) (*models.BudgetModel, error) {
	slug := strings.ToLower(payload.Title)
	slug = strings.Replace(slug, " ", "_", -1)
	model := &models.BudgetModel{
		Amount:      payload.Amount,
		UserID:      UserID,
		Title:       payload.Title,
		Slug:        slug,
		Description: payload.Description,
	}
	if payload.Date == "" {
		currentDate := time.Now()
		model.Date = currentDate
	}
	budgetMonth := uint(model.Date.Month())
	budgetYear := uint16(model.Date.Year())

	model.Month = budgetMonth
	model.Year = budgetYear

	budgetExist, err := b.budgetExistForYearAndMonthAndSlugAndUserID(model.UserID, model.Month, model.Year, model.Slug)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			result := b.DB.Create(model)
			if result.Error != nil {
				return nil, result.Error
			}
			return model, nil
		}
		return nil, err
	}
	return budgetExist, nil
}

func (b *BudgetService) budgetExistForYearAndMonthAndSlugAndUserID(UserID uint, month uint, year uint16, slug string) (*models.BudgetModel, error) {
	retrievedBudget := models.BudgetModel{}
	result := b.DB.Where("user_id = ? AND month = ? AND year = ? AND Slug = ?", UserID, month, year, slug).First(&retrievedBudget)
	if result.Error != nil {
		return nil, result.Error
	}
	return &retrievedBudget, nil
}
func (b *BudgetService) countForYearAndMonthAndSlugAndUserIDExcludeBudgetID(UserID uint, month uint, year uint16, slug string, budgetID uint) int64 {
	var count int64
	log.Info("Passed parameter", UserID, month, year, slug, budgetID)
	b.DB.Model(models.BudgetModel{}).Where("user_id = ? AND month = ? AND year = ? AND Slug = ? AND id <> ?", UserID, month, year, slug, budgetID).Count(&count)
	return count
}

func (b *BudgetService) List(query *gorm.DB, budgets []*models.BudgetModel, paginator *common.Pagination) (*common.Pagination, error) {
	query.Scopes(paginator.Paginate()).Find(&budgets)
	paginator.Items = budgets
	return paginator, nil
}

func (b *BudgetService) GetById(id uint) (*models.BudgetModel, error) {
	var budget models.BudgetModel
	result := b.DB.First(&budget, id) // select * from budgets where id = 1 limit 1;
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, app_errors.NewNotFoundError("Budget not found")
		}
		return nil, errors.New("failed to find budget")
	}
	return &budget, nil
}

func (b *BudgetService) Update(budget *models.BudgetModel, payload *requests.UpdateBudgetRequest, id uint) (*models.BudgetModel, error) {
	if payload.Date != "" {
		timeParsed, err := time.Parse(time.DateOnly, payload.Date)
		if err != nil {
			return nil, errors.New("invalid date passed")
		}
		budget.Date = timeParsed
	}
	if payload.Amount > 0 {
		budget.Amount = payload.Amount
	}

	if payload.Description != nil {
		budget.Description = payload.Description
	}

	if payload.Title != "" {
		budget.Title = payload.Title
		slug := strings.ToLower(payload.Title)
		slug = strings.Replace(slug, " ", "_", -1)
		budget.Slug = slug
	}

	// user_id_slug_year_month // must always be unique
	count := b.countForYearAndMonthAndSlugAndUserIDExcludeBudgetID(budget.UserID, budget.Month, budget.Year, budget.Slug, id)
	log.Info("Count of budgets", count)
	if count > 0 {
		return nil, errors.New("budget with selected month, year and title already exist")
	}
	b.DB.Model(&budget).Updates(budget)
	return budget, nil
}
func (b *BudgetService) GetBudgetsByCategoryID(db *gorm.DB, categoryID uint) ([]*models.BudgetModel, error) {
	if db == nil {
		db = b.DB
	}
	var budgets []*models.BudgetModel
	result := db.Model(models.BudgetModel{}).Joins("JOIN budget_categories ON budget_categories.budget_model_id = budgets.id").
		Joins("JOIN categories ON categories.id = budget_categories.category_model_id").
		Where("categories.id = ?", categoryID).
		Find(&budgets)
	if result.Error != nil {
		fmt.Println(result.Error.Error())
		return nil, errors.New("could not get budgets")
	}
	return budgets, nil
}

func (b *BudgetService) DecrementBudgetBalance(db *gorm.DB, categoryID *uint, amount float64, userID uint) {
	if categoryID == nil {
		return
	}
	budgets, _ := b.GetBudgetsByCategoryID(db, *categoryID)
	for _, budget := range budgets {
		db.Model(models.BudgetModel{}).Scopes(common.WhereUserIDScope(userID)).Where("id", budget.ID).Updates(models.BudgetModel{Amount: budget.Amount - amount})
	}

}

func (b *BudgetService) IncrementBudgetBalance(db *gorm.DB, categoryID *uint, amount float64, userID uint) {
	if categoryID == nil {
		return
	}
	budgets, _ := b.GetBudgetsByCategoryID(db, *categoryID)
	fmt.Println("Increment budget balance", amount)
	for _, budget := range budgets {
		db.Model(models.BudgetModel{}).Scopes(common.WhereUserIDScope(userID)).Updates(models.BudgetModel{Amount: budget.Amount + amount})
	}
}

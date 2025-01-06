package services

import (
	"errors"
	"fmt"
	"github.com/harmlessprince/bougette-backend/cmd/api/requests"
	"github.com/harmlessprince/bougette-backend/common"
	"github.com/harmlessprince/bougette-backend/internal/app_errors"
	"github.com/harmlessprince/bougette-backend/internal/models"
	"gorm.io/gorm"
	"time"
)

const INCOME = "income"
const EXPENSE = "expense"

type TransactionService struct {
	DB *gorm.DB
}

func NewTransactionService(db *gorm.DB) *TransactionService {
	return &TransactionService{DB: db}
}

func (ts *TransactionService) List(transactions []*models.TransactionModel, userID uint, pagination *common.Pagination) (*common.Pagination, error) {
	ts.DB.Scopes(pagination.Paginate(), common.WhereUserIDWithTableNameScope(userID, "transactions")).Joins("Category").Joins("Wallet").Find(&transactions)
	pagination.Items = transactions
	return pagination, nil
}

func (ts *TransactionService) Create(payload *requests.StoreTransactionRequest, userID uint, isReversal bool, formatedDate *time.Time) (*models.TransactionModel, error) {
	transaction := models.TransactionModel{
		CategoryID:  payload.CategoryID,
		Description: payload.Description,
		Amount:      payload.Amount,
		Title:       payload.Title,
		UserID:      userID,
		Type:        payload.Type,
		WalletID:    payload.WalletID,
		IsReversal:  isReversal,
		Date:        *formatedDate,
		Month:       uint(formatedDate.Month()),
		Year:        uint16(formatedDate.Year()),
	}

	if isReversal {
		transaction.ParentID = payload.ParentID
	}
	result := ts.DB.Create(&transaction)
	if result.Error != nil {
		return nil, result.Error
	}
	budgetService := NewBudgetService(ts.DB)
	if transaction.Type == EXPENSE {
		budgetService.DecrementBudgetBalance(ts.DB, transaction.CategoryID, transaction.Amount, transaction.UserID)
	}
	if transaction.Type == INCOME && isReversal {
		budgetService.IncrementBudgetBalance(ts.DB, transaction.CategoryID, transaction.Amount, transaction.UserID)
	}
	return &transaction, nil
}

func (ts *TransactionService) FormatDate(date string) (*time.Time, error) {
	currentTime := time.Now()
	if date == "" {
		return &currentTime, nil
	}
	suppliedDate, err := time.Parse(time.DateOnly, date)
	if err != nil {
		return nil, errors.New("Invalid date format")
	}
	suppliedDateTime := time.Date(
		suppliedDate.Year(), suppliedDate.Month(), suppliedDate.Day(), currentTime.Hour(), currentTime.Minute(), currentTime.Second(), currentTime.Nanosecond(), time.UTC,
	)
	return &suppliedDateTime, nil
}

func (ts *TransactionService) Reverse(transaction *models.TransactionModel) error {
	transactionDate := time.Now()
	description := "Reversed transaction"
	transactionRequest := requests.StoreTransactionRequest{
		WalletID:    transaction.WalletID,
		Amount:      transaction.Amount,
		Date:        transactionDate.Format(time.DateOnly),
		Description: &description,
		ParentID:    &transaction.ID,
	}
	err := ts.DB.Transaction(func(tx *gorm.DB) error {
		walletService := NewWalletService(tx)
		if transaction.Type == INCOME {
			transactionRequest.Type = EXPENSE
			err := walletService.DecrementWalletBalance(tx, transaction.Wallet, transactionRequest.Amount)
			if err != nil {
				fmt.Println(err)
				return errors.New("Transaction could not be reversed, try again later")
			}

		}
		if transaction.Type == EXPENSE {
			transactionRequest.Type = INCOME
			err := walletService.IncrementWalletBalance(tx, transaction.Wallet, transactionRequest.Amount)
			if err != nil {
				fmt.Println(err)
				return errors.New("Transaction could not be reversed, try again later")
			}
		}
		fmt.Println("Got here")
		ts.DB = tx
		_, err := ts.Create(&transactionRequest, transaction.UserID, true, &transaction.Date)
		if err != nil {
			fmt.Println(err)
			return errors.New("Transaction could not be reversed, try again later")
		}
		fmt.Println("Got here after reversal")
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}

func (ts *TransactionService) FindById(db *gorm.DB, id uint) (*models.TransactionModel, error) {
	var transaction models.TransactionModel
	result := db.Joins("Wallet").Joins("Category").First(&transaction, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, app_errors.NewNotFoundError("Transaction not found")
		}
		return nil, result.Error
	}
	return &transaction, nil
}

func (ts *TransactionService) ReversalExist(id uint, userId uint) (*models.TransactionModel, error) {
	var transaction models.TransactionModel
	result := ts.DB.Model(models.TransactionModel{}).Scopes(common.WhereUserIDScope(userId)).Where("parent_id = ?", id).First(&transaction)
	fmt.Println("Error: {}", id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, app_errors.NewNotFoundError("Transaction not found")
		}
		return nil, result.Error
	}
	return &transaction, nil
}

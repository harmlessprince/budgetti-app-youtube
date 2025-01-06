package services

import (
	"github.com/harmlessprince/bougette-backend/cmd/api/requests"
	"github.com/harmlessprince/bougette-backend/internal/models"
	"gorm.io/gorm"
	"time"
)

const transfer = "transfer"

type TransferService struct {
	DB *gorm.DB
}

func NewTransferService(db *gorm.DB) *TransferService {
	return &TransferService{DB: db}
}

func (tr *TransferService) Transfer(
	sourceWallet *models.WalletModel,
	destinationWallet *models.WalletModel,
	amount float64,
	userId uint,
) error {
	categoryService := NewCategoryService(tr.DB)
	walletService := NewWalletService(tr.DB)
	// type of category
	category, err := categoryService.findBySlug(transfer, false)
	if err != nil {
		return err
	}
	title := "Account transfer"
	sourceTransactionData := requests.StoreTransactionRequest{
		CategoryID:  &category.ID,
		WalletID:    sourceWallet.ID,
		Amount:      amount,
		Date:        time.Now().Format(time.DateOnly),
		Description: &title,
		Title:       &title,
		Type:        EXPENSE,
	}

	destinationTransactionData := requests.StoreTransactionRequest{
		CategoryID:  &category.ID,
		WalletID:    destinationWallet.ID,
		Amount:      amount,
		Date:        time.Now().Format(time.DateOnly),
		Description: &title,
		Title:       &title,
		Type:        INCOME,
	}

	err = tr.DB.Transaction(func(tx *gorm.DB) error {
		transactionService := NewTransactionService(tx)
		formattedDate, _ := transactionService.FormatDate(sourceTransactionData.Date)

		err = walletService.DecrementWalletBalance(tx, sourceWallet, amount)
		if err != nil {
			return err
		}
		_, err = transactionService.Create(&sourceTransactionData, userId, false, formattedDate)

		err = walletService.IncrementWalletBalance(tx, destinationWallet, amount)
		if err != nil {
			return err
		}
		_, err = transactionService.Create(&destinationTransactionData, userId, false, formattedDate)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}

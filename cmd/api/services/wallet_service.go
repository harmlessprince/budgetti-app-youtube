package services

import (
	"github.com/harmlessprince/bougette-backend/cmd/api/requests"
	"github.com/harmlessprince/bougette-backend/internal/models"
	"gorm.io/gorm"
)

type WalletService struct {
	DB *gorm.DB
}

func NewWalletService(db *gorm.DB) *WalletService {
	return &WalletService{DB: db}
}

func (w WalletService) Create(data requests.CreateWalletRequest, UserID uint) (*models.WalletModel, error) {
	wallet := models.WalletModel{
		UserID:  UserID,
		Balance: data.Amount,
		Name:    data.Name,
	}
	result := w.DB.Create(&wallet)
	if result.Error != nil {
		return nil, result.Error
	}
	return &wallet, nil
}

func (w WalletService) WalletExistForNameAndUserID(name string, userID uint) (*models.WalletModel, error) {
	wallet := models.WalletModel{}
	result := w.DB.Where("user_id = ? AND name = ?", userID, name).First(&wallet)
	if result.Error != nil {
		return nil, result.Error
	}
	return &wallet, nil
}

func (w WalletService) GenerateDefaultWallets(id uint) ([]*models.WalletModel, error) {
	wallets := []string{"Cash", "Bank"}
	var walletsCreated []*models.WalletModel
	for _, name := range wallets {
		walletExist, _ := w.WalletExistForNameAndUserID(name, id)
		if walletExist != nil {
			walletsCreated = append(walletsCreated, walletExist)
			continue
		}
		walletRequest := requests.CreateWalletRequest{
			Name:   name,
			Amount: 0,
		}
		wallet, err := w.Create(walletRequest, id)
		if err != nil {
			return nil, err
		}
		walletsCreated = append(walletsCreated, wallet)
	}
	return walletsCreated, nil
}

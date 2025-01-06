package services

import (
	"github.com/harmlessprince/bougette-backend/cmd/api/requests"
	"github.com/harmlessprince/bougette-backend/common"
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

func (w WalletService) List(userID uint) ([]*models.WalletModel, error) {
	var wallets []*models.WalletModel
	result := w.DB.Model(models.WalletModel{}).Where("user_id = ?", userID).Find(&wallets)
	if result.Error != nil {
		return nil, result.Error
	}
	return wallets, nil
}
func (w WalletService) FindById(walletID uint, userID uint) (*models.WalletModel, error) {
	var wallet models.WalletModel
	result := w.DB.Scopes(common.WhereUserIDScope(userID)).First(&wallet, walletID)
	if result.Error != nil {
		return nil, result.Error
	}
	return &wallet, nil
}

func (w WalletService) IncrementWalletBalance(tx *gorm.DB, wallet *models.WalletModel, amount float64) error {
	balance := wallet.Balance + amount
	wallet.Balance = balance
	//result := tx.Model(wallet).Updates(models.WalletModel{Balance: balance})
	result := tx.Model(&models.WalletModel{}).Where("id = ?", wallet.ID).Update("balance", balance)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (w WalletService) DecrementWalletBalance(tx *gorm.DB, wallet *models.WalletModel, amount float64) error {
	balance := wallet.Balance - amount // allow negative balance
	wallet.Balance = balance
	result := tx.Model(&models.WalletModel{}).Where("id = ?", wallet.ID).Update("balance", balance)
	//result := tx.Model(wallet).Updates(models.WalletModel{Balance: balance})
	if result.Error != nil {
		return result.Error
	}
	return nil
}

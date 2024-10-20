package services

import (
	"errors"
	"github.com/harmlessprince/bougette-backend/internal/models"
	"gorm.io/gorm"
	"math/rand"
	"strconv"
	"time"
)

type AppTokenService struct {
	db *gorm.DB
}

func NewAppTokenService(db *gorm.DB) *AppTokenService {
	return &AppTokenService{db: db}
}

func (appTokenService *AppTokenService) getToken() int {
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(99999-10000+1) + 10000
}

func (appTokenService *AppTokenService) GenerateResetPasswordToken(user models.UserModel) (*models.AppTokenModel, error) {
	tokenCreated := models.AppTokenModel{
		TargetId:  user.ID,
		Type:      "reset_password",
		Token:     strconv.Itoa(appTokenService.getToken()),
		Used:      false,
		ExpiresAt: time.Now().Add(time.Hour * 24),
	}
	result := appTokenService.db.Create(&tokenCreated)
	if result.Error != nil {
		return nil, result.Error
	}
	return &tokenCreated, nil
}

func (appTokenService *AppTokenService) ValidateResetPasswordToken(user models.UserModel, token string) (*models.AppTokenModel, error) {
	var retrievedToken models.AppTokenModel
	result := appTokenService.db.Where(&models.AppTokenModel{
		TargetId: user.ID,
		Type:     "reset_password",
		Token:    token,
	}).First(&retrievedToken)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, errors.New("invalid password reset token")
		}
		return nil, result.Error
	}
	if retrievedToken.Used {
		return nil, errors.New("invalid password reset token")
	}

	if retrievedToken.ExpiresAt.Before(time.Now()) {
		return nil, errors.New("password reset token expired, please re-initiate forgot password")
	}
	return &retrievedToken, nil
}

func (appTokenService *AppTokenService) InvalidateToken(userId uint, appToken models.AppTokenModel) {
	appTokenService.db.Model(&models.AppTokenModel{}).Where("target_id = ? AND token = ?", userId, appToken.Token).Update("used", true)
}

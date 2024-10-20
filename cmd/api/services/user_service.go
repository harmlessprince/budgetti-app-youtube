package services

import (
	"errors"
	"fmt"
	"github.com/harmlessprince/bougette-backend/cmd/api/requests"
	"github.com/harmlessprince/bougette-backend/common"
	"github.com/harmlessprince/bougette-backend/internal/models"
	"gorm.io/gorm"
)

type UserService struct {
	db *gorm.DB
}

func NewUserService(db *gorm.DB) *UserService {
	return &UserService{db: db}
}

func (userService *UserService) RegisterUser(userRequest *requests.RegisterUserRequest) (*models.UserModel, error) {
	// hash the password
	hashedPassword, err := common.HashPassword(userRequest.Password)
	if err != nil {
		fmt.Println(err)
		return nil, errors.New("registration failed")
	}
	createdUser := models.UserModel{
		FirstName: &userRequest.FirstName,
		LastName:  &userRequest.LastName,
		Email:     userRequest.Email,
		Password:  hashedPassword,
	}
	result := userService.db.Create(&createdUser)
	if result.Error != nil {
		fmt.Println(result.Error)
		return nil, errors.New("registration failed")
	}
	return &createdUser, nil
}

func (userService *UserService) GetUserByEmail(email string) (*models.UserModel, error) {
	var user models.UserModel
	result := userService.db.Where("email = ?", email).First(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}

func (userService *UserService) ChangeUserPassword(newPassword string, user models.UserModel) error {
	hashedPassword, err := common.HashPassword(newPassword)
	if err != nil {
		fmt.Println(err)
		return errors.New("password change failed")
	}

	result := userService.db.Model(user).Update("Password", hashedPassword)
	if result.Error != nil {
		fmt.Println(result.Error)
		return errors.New("password change failed")
	}
	return nil
}

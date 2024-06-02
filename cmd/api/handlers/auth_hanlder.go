package handlers

import (
	"errors"
	"github.com/harmlessprince/bougette-backend/cmd/api/requests"
	"github.com/harmlessprince/bougette-backend/cmd/api/services"
	"github.com/harmlessprince/bougette-backend/common"
	"github.com/harmlessprince/bougette-backend/internal/mailer"
	"github.com/harmlessprince/bougette-backend/internal/models"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
	"os"
)

func (h *Handler) RegisterHandler(c echo.Context) error {
	// bind request body
	payload := new(requests.RegisterUserRequest)
	if err := (&echo.DefaultBinder{}).BindBody(c, payload); err != nil {
		c.Logger().Error(err)
		return common.SendBadRequestResponse(c, err.Error())
	}

	// validation
	validationErrors := h.ValidateBodyRequest(c, *payload)

	if validationErrors != nil {
		return common.SendFailedValidationResponse(c, validationErrors)
	}
	userService := services.NewUserService(h.DB)
	// Check if email exist
	_, err := userService.GetUserByEmail(payload.Email)
	if errors.Is(err, gorm.ErrRecordNotFound) == false {
		return common.SendBadRequestResponse(c, "Email has already been  taken")
	}
	//Create the user
	registeredUser, err := userService.RegisterUser(payload)
	if err != nil {
		return common.SendInternalServerErrorResponse(c, err.Error())
	}
	mailData := mailer.EmailData{
		Subject: "Welcome To " + os.Getenv("APP_NAME") + " Signup",
		Meta: struct {
			FirstName string
			LoginLink string
		}{
			FirstName: *registeredUser.FirstName,
			LoginLink: "#",
		},
	}
	// Send a welcome message to the user
	err = h.Mailer.Send(payload.Email, "welcome.html", mailData)
	if err != nil {
		h.Logger.Error(err)
	}
	// Send response
	return common.SendSuccessResponse(c, "User registration successful", registeredUser)
}

func (h *Handler) LoginHandler(c echo.Context) error {
	userService := services.NewUserService(h.DB)
	// bind our data/ retrieving data sent by client
	payload := new(requests.LoginRequest)
	if err := (&echo.DefaultBinder{}).BindBody(c, payload); err != nil {
		c.Logger().Error(err)
		return common.SendBadRequestResponse(c, err.Error())
	}

	// validate the data sent by client
	validationErrors := h.ValidateBodyRequest(c, *payload)
	if validationErrors != nil {
		return common.SendFailedValidationResponse(c, validationErrors)
	}
	// if the user with supplied email exist
	userRetrieved, err := userService.GetUserByEmail(payload.Email)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return common.SendBadRequestResponse(c, "Invalid credentials")
	}
	// compare the client password with the hashed password
	if common.ComparePasswordHash(payload.Password, userRetrieved.Password) == false {
		return common.SendBadRequestResponse(c, "Invalid credentials")
	}
	// we send response with user token
	accessToken, refreshToken, err := common.GenerateJWT(*userRetrieved)
	if err != nil {
		return common.SendInternalServerErrorResponse(c, err.Error())
	}

	return common.SendSuccessResponse(c, "User logged in", map[string]interface{}{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
		"user":          userRetrieved,
	})
}

func (h *Handler) GetAuthenticatedUser(c echo.Context) error {
	user, ok := c.Get("user").(models.UserModel)
	if !ok {
		return common.SendInternalServerErrorResponse(c, "User authentication failed")
	}
	return common.SendSuccessResponse(c, "Authenticated user retrieved", user)
}

func (h *Handler) UpdateUserPassword(c echo.Context) error {
	return nil
}

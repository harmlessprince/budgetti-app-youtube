package handlers

import (
	"errors"
	"github.com/harmlessprince/bougette-backend/cmd/api/requests"
	"github.com/harmlessprince/bougette-backend/cmd/api/services"
	"github.com/harmlessprince/bougette-backend/common"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
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
	// Send a welcome message to the user

	// Send response
	return common.SendSuccessResponse(c, "User registration successful", registeredUser)
}

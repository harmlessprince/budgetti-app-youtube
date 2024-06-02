package handlers

import (
	"github.com/harmlessprince/bougette-backend/cmd/api/requests"
	"github.com/harmlessprince/bougette-backend/cmd/api/services"
	"github.com/harmlessprince/bougette-backend/common"
	"github.com/harmlessprince/bougette-backend/internal/models"
	"github.com/labstack/echo/v4"
)

func (h *Handler) ChangeUserPassword(c echo.Context) error {
	user, ok := c.Get("user").(models.UserModel)
	if !ok {
		return common.SendInternalServerErrorResponse(c, "User authentication failed")
	}
	// bind request body
	payload := new(requests.ChangePasswordRequest)
	if err := h.BindBodyRequest(c, payload); err != nil {
		return common.SendBadRequestResponse(c, err.Error())
	}

	// validation
	validationErrors := h.ValidateBodyRequest(c, *payload)

	if validationErrors != nil {
		return common.SendFailedValidationResponse(c, validationErrors)
	}

	// compare the client current  password with the hashed password
	if common.ComparePasswordHash(payload.CurrentPassword, user.Password) == false {
		return common.SendBadRequestResponse(c, "the supplied password does not match your current password")
	}
	userService := services.NewUserService(h.DB)
	err := userService.ChangeUserPassword(payload.Password, user)

	if err != nil {
		return common.SendBadRequestResponse(c, err.Error())
	}
	return common.SendSuccessResponse(c, "Password changed successfully", nil)
}

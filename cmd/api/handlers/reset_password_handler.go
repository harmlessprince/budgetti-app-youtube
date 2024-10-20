package handlers

import (
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/harmlessprince/bougette-backend/cmd/api/requests"
	"github.com/harmlessprince/bougette-backend/cmd/api/services"
	"github.com/harmlessprince/bougette-backend/common"
	"github.com/harmlessprince/bougette-backend/internal/mailer"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
	"net/url"
)

func (h *Handler) ForgotPasswordHandler(c echo.Context) error {
	// bind request body
	payload := new(requests.ForgotPasswordRequest)
	if err := h.BindBodyRequest(c, payload); err != nil {
		return common.SendBadRequestResponse(c, err.Error())
	}

	// validation
	validationErrors := h.ValidateBodyRequest(c, *payload)

	if validationErrors != nil {
		return common.SendFailedValidationResponse(c, validationErrors)
	}

	userService := services.NewUserService(h.DB)
	appTokenService := services.NewAppTokenService(h.DB)
	// Check if email exist
	retrievedUser, err := userService.GetUserByEmail(payload.Email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return common.SendNotFoundResponse(c, "Record not found, register with this email")
		}
		return common.SendInternalServerErrorResponse(c, "An error occurred, try again later")
	}

	token, err := appTokenService.GenerateResetPasswordToken(*retrievedUser)
	if err != nil {
		return common.SendInternalServerErrorResponse(c, "An error occurred, try again later")
	}

	econdedEmail := base64.RawURLEncoding.EncodeToString([]byte(retrievedUser.Email))

	frontendUrl, err := url.Parse(payload.FrontendURL)
	if err != nil {
		return common.SendBadRequestResponse(c, "Invalid frontend URL")
	}
	fmt.Println(econdedEmail)
	query := url.Values{}
	query.Set("email", econdedEmail)
	query.Set("token", token.Token)
	frontendUrl.RawQuery = query.Encode()
	mailData := mailer.EmailData{
		Subject: "Request Password Reset",
		Meta: struct {
			Token       string
			FrontendUrl string
		}{
			Token:       token.Token,
			FrontendUrl: frontendUrl.String(),
		},
	}
	// Send a welcome message to the user
	err = h.Mailer.Send(payload.Email, "forgot-password.html", mailData)
	if err != nil {
		h.Logger.Error(err)
	}
	return common.SendSuccessResponse(c, "Forgot password email sent", nil)
}

func (h *Handler) ResetPasswordHandler(c echo.Context) error {
	// bind request body
	payload := new(requests.ResetPasswordRequest)
	if err := h.BindBodyRequest(c, payload); err != nil {
		return common.SendBadRequestResponse(c, err.Error())
	}

	// validation
	validationErrors := h.ValidateBodyRequest(c, *payload)
	if validationErrors != nil {
		return common.SendFailedValidationResponse(c, validationErrors)
	}

	email, err := base64.RawURLEncoding.DecodeString(payload.Meta)
	if err != nil {
		fmt.Println(string(email))
		fmt.Println(err)
		return common.SendInternalServerErrorResponse(c, "An error occurred, try again later")
	}

	userService := services.NewUserService(h.DB)
	appTokenService := services.NewAppTokenService(h.DB)

	retrievedUser, err := userService.GetUserByEmail(string(email))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return common.SendNotFoundResponse(c, "Invalid password reset token")
		}
		return common.SendInternalServerErrorResponse(c, "An error occurred, try again later")
	}

	appToken, err := appTokenService.ValidateResetPasswordToken(*retrievedUser, payload.Token)

	if err != nil {
		return common.SendInternalServerErrorResponse(c, err.Error())
	}
	err = userService.ChangeUserPassword(payload.Password, *retrievedUser)
	if err != nil {
		return common.SendInternalServerErrorResponse(c, err.Error())
	}
	appTokenService.InvalidateToken(retrievedUser.ID, *appToken)
	return common.SendSuccessResponse(c, "Reset password successful", nil)
}

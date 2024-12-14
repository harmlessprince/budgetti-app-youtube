package handlers

import (
	"github.com/harmlessprince/bougette-backend/cmd/api/requests"
	"github.com/harmlessprince/bougette-backend/cmd/api/services"
	"github.com/harmlessprince/bougette-backend/common"
	"github.com/harmlessprince/bougette-backend/internal/models"
	"github.com/labstack/echo/v4"
)

func (h *Handler) CreateWallet(c echo.Context) error {
	user, _ := c.Get("user").(models.UserModel)

	// bind request body
	payload := new(requests.CreateWalletRequest)
	if err := h.BindBodyRequest(c, payload); err != nil {
		return common.SendBadRequestResponse(c, err.Error())
	}

	// validation
	validationErrors := h.ValidateBodyRequest(c, *payload)

	if validationErrors != nil {
		return common.SendFailedValidationResponse(c, validationErrors)
	}

	walletService := services.NewWalletService(h.DB)

	walletExist, _ := walletService.WalletExistForNameAndUserID(payload.Name, user.ID)

	if walletExist != nil {
		return common.SendBadRequestResponse(c, "A wallet with name: "+payload.Name+" already exists")
	}
	wallet, err := walletService.Create(*payload, user.ID)
	if err != nil {
		return common.SendInternalServerErrorResponse(c, "An error occurred try again later")
	}
	return common.SendSuccessResponse(c, "Wallet created successfully", wallet)
}

func (h *Handler) GenerateDefaultWallets(c echo.Context) error {
	user, _ := c.Get("user").(models.UserModel)

	walletService := services.NewWalletService(h.DB)

	wallets, err := walletService.GenerateDefaultWallets(user.ID)
	if err != nil {
		return common.SendInternalServerErrorResponse(c, "An error occurred try again later")
	}
	return common.SendSuccessResponse(c, "Default Wallet created successfully", wallets)
}

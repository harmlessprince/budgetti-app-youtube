package handlers

import (
	"github.com/harmlessprince/bougette-backend/cmd/api/requests"
	"github.com/harmlessprince/bougette-backend/cmd/api/services"
	"github.com/harmlessprince/bougette-backend/common"
	"github.com/harmlessprince/bougette-backend/internal/models"
	"github.com/labstack/echo/v4"
)

func (h *Handler) Transfer(c echo.Context) error {
	user, _ := c.Get("user").(models.UserModel)
	// bind request body
	payload := new(requests.TransferRequest)
	if err := h.BindBodyRequest(c, payload); err != nil {
		return common.SendBadRequestResponse(c, err.Error())
	}

	// validation
	validationErrors := h.ValidateBodyRequest(c, *payload)

	if validationErrors != nil {
		return common.SendFailedValidationResponse(c, validationErrors)
	}
	walletService := services.NewWalletService(h.DB)
	transferService := services.NewTransferService(h.DB)

	sourceWallet, err := walletService.FindById(payload.SourceWalletID, user.ID)
	if err != nil {
		return common.SendBadRequestResponse(c, "invalid source wallet ID")
	}

	destinationWallet, err := walletService.FindById(payload.DestinationWalletID, user.ID)
	if err != nil {
		return common.SendBadRequestResponse(c, "invalid destination wallet ID")
	}

	if payload.Amount >= sourceWallet.Balance {
		return common.SendBadRequestResponse(c, "Insufficient funds")
	}
	err = transferService.Transfer(sourceWallet, destinationWallet, payload.Amount, user.ID)
	if err != nil {
		return common.SendInternalServerErrorResponse(c, "error transferring from wallet, try again later")
	}
	return common.SendSuccessResponse(c, "Transfer successful", nil)
}

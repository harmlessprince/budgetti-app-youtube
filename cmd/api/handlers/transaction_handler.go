package handlers

import (
	"errors"
	"fmt"
	"github.com/harmlessprince/bougette-backend/cmd/api/filters"
	"github.com/harmlessprince/bougette-backend/cmd/api/requests"
	"github.com/harmlessprince/bougette-backend/cmd/api/services"
	"github.com/harmlessprince/bougette-backend/common"
	"github.com/harmlessprince/bougette-backend/internal/app_errors"
	"github.com/harmlessprince/bougette-backend/internal/models"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func (h *Handler) ListTransactions(c echo.Context) error {
	user, _ := c.Get("user").(models.UserModel)
	// bind request body
	filter := new(filters.TransactionFilter)
	err := (&echo.DefaultBinder{}).BindQueryParams(c, filter)
	if err != nil {
		return err
	}
	err = filter.ValidateDate()
	if err != nil {
		return common.SendBadRequestResponse(c, err.Error())
	}

	fmt.Println(filter)

	query := filter.ApplyFilters(h.DB)

	transactionService := services.NewTransactionService(query)

	var transactions []*models.TransactionModel

	paginator := common.NewPaginator(transactions, c.Request(), query)
	paginatedTransactions, err := transactionService.List(transactions, user.ID, paginator)
	if err != nil {
		return common.SendInternalServerErrorResponse(c, err.Error())
	}
	return common.SendSuccessResponse(c, "ok", paginatedTransactions)
}

func (h *Handler) StoreTransaction(c echo.Context) error {
	user, _ := c.Get("user").(models.UserModel)
	// bind request body
	payload := new(requests.StoreTransactionRequest)
	if err := h.BindBodyRequest(c, payload); err != nil {
		return common.SendBadRequestResponse(c, err.Error())
	}

	// validation
	validationErrors := h.ValidateBodyRequest(c, *payload)

	if validationErrors != nil {
		return common.SendFailedValidationResponse(c, validationErrors)
	}
	categoryService := services.NewCategoryService(h.DB)
	walletService := services.NewWalletService(h.DB)
	formatedDate, err := services.NewTransactionService(h.DB).FormatDate(payload.Date)
	if err != nil {
		return common.SendBadRequestResponse(c, "Invalid date format")
	}
	var category *models.CategoryModel
	// retrieve the wallet
	wallet, _ := walletService.FindById(payload.WalletID, user.ID)
	if wallet == nil {
		return common.SendNotFoundResponse(c, "Invalid wallet")
	}
	// retrieve the category if provided

	if payload.CategoryID != nil {
		retrievedCategory, err := categoryService.GetById(*payload.CategoryID)
		if err != nil {
			return err
		}
		category = retrievedCategory
	}

	var transaction *models.TransactionModel

	err = h.DB.Transaction(func(tx *gorm.DB) error {
		//register income
		if payload.Type == services.INCOME {
			// increment the wallet balance
			err := walletService.IncrementWalletBalance(tx, wallet, payload.Amount)
			if err != nil {
				return errors.New("Transaction failed, try again later")
			}
		}
		if payload.Type == services.EXPENSE {
			// register an expense
			// decrement the wallet balance
			err := walletService.DecrementWalletBalance(tx, wallet, payload.Amount)
			if err != nil {
				return errors.New("Transaction failed, try again later")
			}
		}
		transactionService := services.NewTransactionService(tx)
		// create a transaction of type income or expense
		createdTransaction, err := transactionService.Create(payload, user.ID, false, formatedDate)
		if err != nil {
			return errors.New("Transaction failed, try again later")
		}
		transaction = createdTransaction
		transaction.Category = category
		transaction.Wallet = wallet
		return nil
	})
	if err != nil {
		return common.SendInternalServerErrorResponse(c, err.Error())
	}

	return common.SendSuccessResponse(c, "Transaction created successfully", transaction)
}

func (h *Handler) ReverseTransaction(c echo.Context) error {
	var transactionID requests.IDParamRequest
	err := (&echo.DefaultBinder{}).BindPathParams(c, &transactionID)
	if err != nil {
		return common.SendBadRequestResponse(c, err.Error())
	}
	user, _ := c.Get("user").(models.UserModel)
	transactionService := services.NewTransactionService(h.DB)
	query := transactionService.DB.Scopes(common.WhereUserIDWithTableNameScope(user.ID, "transactions"))
	transaction, err := transactionService.FindById(query, transactionID.ID)

	if err != nil {
		if errors.Is(err, app_errors.NewNotFoundError(err.Error())) {
			return common.SendNotFoundResponse(c, err.Error())
		}
		return common.SendInternalServerErrorResponse(c, err.Error())
	}

	if transaction.IsReversal {
		return common.SendBadRequestResponse(c, "You can not reverse a transaction of type reversal")
	}

	transactionReversalExist, _ := transactionService.ReversalExist(transactionID.ID, user.ID)

	if transactionReversalExist != nil {
		return common.SendBadRequestResponse(c, "The reversal has already been processed")
	}

	err = transactionService.Reverse(transaction)
	if err != nil {
		return common.SendInternalServerErrorResponse(c, err.Error())
	}
	return common.SendSuccessResponse(c, "Transaction reversed", nil)
}

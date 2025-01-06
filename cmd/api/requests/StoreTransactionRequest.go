package requests

type StoreTransactionRequest struct {
	CategoryID  *uint   `json:"category_id" validate:"omitnil,number"`
	WalletID    uint    `json:"wallet_id" validate:"required,number"`
	Amount      float64 `json:"amount" validate:"required,numeric,min=1"`
	Date        string  `json:"date" validate:"omitempty,datetime=2006-01-02"`
	Description *string `json:"description" validate:"omitnil,max=490,min=2"`
	Title       *string `json:"title" validate:"omitnil,max=100,min=2"`
	Type        string  `json:"type" validate:"required,oneof=income expense"`
	ParentID    *uint
}

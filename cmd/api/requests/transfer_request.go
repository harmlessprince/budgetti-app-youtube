package requests

type TransferRequest struct {
	SourceWalletID      uint    `json:"source_wallet_id" validation:"required,number"`
	DestinationWalletID uint    `json:"destination_wallet_id" validation:"required,number"`
	Amount              float64 `json:"amount" validation:"required,numeric"`
}

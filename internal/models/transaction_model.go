package models

import "time"

type TransactionModel struct {
	BaseModel
	ParentID    *uint          `gorm:"column:parent_id" json:"-"`
	Title       *string        `gorm:"type:varchar(200)" json:"title"`
	Description *string        `gorm:"type:varchar(500)" json:"description"`
	UserID      uint           `gorm:"not null;column:user_id" json:"user_id"`
	CategoryID  *uint          `gorm:"column:category_id" json:"category_id"`
	WalletID    uint           `gorm:"not null;column:wallet_id" json:"wallet_id"`
	Amount      float64        `gorm:"not null;column:amount" json:"amount"`
	Date        time.Time      `gorm:"type:datetime;not null" json:"date"`
	Month       uint           `gorm:"type:TINYINT UNSIGNED;not null;index:idx_month_year;" json:"month"`
	Year        uint16         `gorm:"type:INT UNSIGNED;not null;index:idx_month_year;" json:"year"`
	Type        string         `gorm:"not null;type:varchar(100);index" json:"type"`
	IsReversal  bool           `gorm:"not null;type:boolean;default:false;" json:"is_reversal"`
	Category    *CategoryModel `gorm:"foreignkey:CategoryID;constraint:OnDelete:CASCADE;" json:"category"`
	Wallet      *WalletModel   `gorm:"foreignkey:WalletID;constraint:OnDelete:CASCADE;" json:"wallet"`
	User        *UserModel     `gorm:"foreignkey:UserID;constraint:OnDelete:CASCADE;" json:"-"`
}

func (TransactionModel) TableName() string {
	return "transactions"
}

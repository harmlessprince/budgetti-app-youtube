package models

type WalletModel struct {
	BaseModel
	UserID  uint      `gorm:"not null;column:user_id;uniqueIndex:unique_userid_name" json:"user_id"`
	Balance float64   `gorm:"not null;default:0.0;type:double precision;" json:"balance"`
	Name    string    `gorm:"not null;size:100;index;uniqueIndex:unique_userid_name" json:"name"`
	Owner   UserModel `gorm:"foreignkey:UserID" json:"-"`
}

// Wallet Or Account (Cash and Bank)

func (w WalletModel) TableName() string {
	return "wallets"
}

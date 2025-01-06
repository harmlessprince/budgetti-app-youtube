package main

import (
	"github.com/harmlessprince/bougette-backend/common"
	"github.com/harmlessprince/bougette-backend/internal/models"
	"log"
)

func main() {
	db, err := common.NewMysql()
	if err != nil {
		panic(err)
	}
	err = db.AutoMigrate(
		&models.UserModel{},
		&models.AppTokenModel{},
		&models.CategoryModel{},
		&models.UserCategoryModel{},
		&models.BudgetModel{},
		&models.WalletModel{},
		&models.TransactionModel{},
	)
	if err != nil {
		panic(err)
	}
	log.Println("Migration completed")
}

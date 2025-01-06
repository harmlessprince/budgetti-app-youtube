package main

import (
	"github.com/harmlessprince/bougette-backend/cmd/api/requests"
	"github.com/harmlessprince/bougette-backend/cmd/api/services"
	"github.com/harmlessprince/bougette-backend/common"
)

func main() {
	db, err := common.NewMysql()
	if err != nil {
		panic(err)
	}
	categroyService := services.CategoryService{
		DB: db,
	}

	categories := []string{
		"Food", "Gifts", "Health", "Eat out",
		"Medical", "Home", "Transportation",
		"Personnel", "Pets", "Utilities", "Travel",
		"Debt", "Other", "Savings", "Paycheck", "Bonus",
		"Interest", "Internet", "Calls", "Laundry", "Charity",
		"Lendings", "Family Obligation", "Loans", "Grocery", "Chocolate",
		"Transfer",
	}
	for _, category := range categories {
		_, err = categroyService.Create(&requests.CreateCategoryRequest{
			Name:     category,
			IsCustom: false,
		})
		if err != nil {
			panic(err)
		}
		println("Category " + category + " created")
	}
}

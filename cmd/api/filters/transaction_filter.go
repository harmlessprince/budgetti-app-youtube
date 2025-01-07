package filters

import (
	"errors"
	"fmt"
	"github.com/jinzhu/now"
	"gorm.io/gorm"
	"time"
)

type TransactionFilter struct {
	FromDate    string `query:"from_date"`
	EnDate      string `query:"end_date"`
	CategoryID  uint   `query:"category_id"`
	WalletID    uint   `query:"wallet_id"`
	Type        string `query:"type"`
	Month       uint   `query:"month"`
	Year        uint   `query:"year"`
	Title       string `query:"title"`
	Description string `query:"description"`
}

func (f TransactionFilter) ApplyFilters(query *gorm.DB) *gorm.DB {
	query = f.filterByWalletID(query)
	query = f.filterByCategoryID(query)
	query = f.filterByMonth(query)
	query = f.filterByType(query)
	query = f.filterByYear(query)
	query = f.filterByDate(query)
	query = f.filterByDescription(query)
	query = f.filterByTitle(query)
	return query
}

func (f TransactionFilter) filterByWalletID(query *gorm.DB) *gorm.DB {
	if f.WalletID == 0 {
		return query
	}
	return query.Where("wallet_id = ?", f.WalletID)
}

func (f TransactionFilter) filterByCategoryID(query *gorm.DB) *gorm.DB {
	if f.CategoryID == 0 {
		return query
	}
	return query.Where("category_id = ?", f.CategoryID)
}

func (f TransactionFilter) filterByType(query *gorm.DB) *gorm.DB {
	if f.Type == "" {
		return query
	}
	return query.Where("type = ?", f.Type)
}

func (f TransactionFilter) filterByMonth(query *gorm.DB) *gorm.DB {
	if f.Month == 0 {
		return query
	}
	return query.Where("month = ?", f.Month)
}

func (f TransactionFilter) filterByYear(query *gorm.DB) *gorm.DB {
	if f.Year == 0 {
		return query
	}
	return query.Where("year = ?", f.Year)
}

func (f TransactionFilter) filterByDate(query *gorm.DB) *gorm.DB {
	startDate := now.BeginningOfMonth()
	endDate := now.EndOfMonth()
	if f.FromDate == "" && f.EnDate == "" {
		f.FromDate = startDate.Format("2006-01-02")
		f.EnDate = endDate.Format("2006-01-02")
	} else {
		parsedFromDate, err := time.Parse(time.DateOnly, f.FromDate)
		if err != nil {
			return query
		}
		parsedEnDate, err := time.Parse(time.DateOnly, f.EnDate)
		if err != nil {
			return query
		}
		f.FromDate = parsedFromDate.Format("2006-01-02")
		f.EnDate = parsedEnDate.Format("2006-01-02")
	}
	fmt.Printf("Start Date: %s", f.FromDate)
	fmt.Printf("End Date: %s", f.EnDate)
	return query.Where("date BETWEEN ? AND ?", f.FromDate, f.EnDate)
}

func (f TransactionFilter) ValidateDate() error {
	// check if either date is suppplied
	fmt.Println("From Date:" + f.FromDate)
	fmt.Println("End Date: " + f.EnDate)
	if f.FromDate != "" || f.EnDate != "" {
		if f.FromDate == "" || f.EnDate == "" {
			return errors.New("both from_date and end_date are required")
		}
		parsedFromDate, err := time.Parse(time.DateOnly, f.FromDate)
		if err != nil {
			return errors.New("invalid from_date, expected format is 2006-01-02 (YYYY-MM-DD)")
		}
		parsedEnDate, err := time.Parse(time.DateOnly, f.EnDate)
		if err != nil {
			return errors.New("invalid end_date, expected format is 2006-01-02 (YYYY-MM-DD)")
		}

		if parsedFromDate.After(parsedEnDate) {
			return errors.New("from_date must be after end_date")
		}
	}
	return nil
}

func (f TransactionFilter) filterByDescription(query *gorm.DB) *gorm.DB {
	if f.Description == "" {
		return query
	}
	return query.Where("description LIKE ?", "%"+f.Description+"%")
}

func (f TransactionFilter) filterByTitle(query *gorm.DB) *gorm.DB {
	if f.Title == "" {
		return query
	}
	return query.Where("title LIKE ?", "%"+f.Title+"%")
}

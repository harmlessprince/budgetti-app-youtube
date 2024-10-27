package common

import (
	"gorm.io/gorm"
	"math"
	"net/http"
	"strconv"
)

type Pagination struct {
	Limit     int         `query:"limit" json:"limit"`
	Page      int         `query:"page" json:"page"`
	Sort      string      `query:"sort" json:"sort"`
	TotalRows int64       `json:"total_rows"`
	TotalPage int         `json:"total_page"`
	Items     interface{} `json:"items"`
}

func (p *Pagination) GetPage() int {
	if p.Page <= 0 {
		p.Page = 1
	}
	return p.Page
}

func (p *Pagination) GetLimit() int {
	if p.Limit > 100 {
		p.Limit = 100
	} else if p.Limit <= 0 {
		p.Limit = 10
	}
	return p.Limit
}

func (p *Pagination) GetOffset() int {
	return (p.GetPage() - 1) * p.GetLimit()
}

func NewPaginator(model interface{}, r *http.Request, db *gorm.DB) *Pagination {
	var pagination Pagination
	q := r.URL.Query()
	page, _ := strconv.Atoi(q.Get("page"))
	limit, _ := strconv.Atoi(q.Get("limit"))

	var totalRows int64
	db.Model(model).Count(&totalRows)
	pagination.Limit = limit
	pagination.Page = page
	pagination.TotalRows = totalRows
	totalPage := int(math.Ceil(float64(totalRows) / float64(pagination.GetLimit())))
	pagination.TotalPage = totalPage
	return &pagination
}
func (p *Pagination) Paginate() func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Offset(p.GetOffset()).Limit(p.GetLimit())
	}
}

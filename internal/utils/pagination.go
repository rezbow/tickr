package utils

import "gorm.io/gorm"

type Pagination struct {
	Page     int `json:"page" form:"page"`
	PageSize int `json:"page_size" form:"page_size"`
}

func (p *Pagination) Paginate(db *gorm.DB) *gorm.DB {
	if p.Page <= 0 {
		p.Page = 1
	}
	if p.PageSize <= 0 {
		p.PageSize = 10
	}
	offset := (p.Page - 1) * p.PageSize
	return db.Offset(offset).Limit(p.PageSize)
}

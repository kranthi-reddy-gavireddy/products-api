package dtos

import (
	"math"
)

type SortClause struct {
	Field     string
	Direction string
}

type ListProductsResponse struct {
	Products   []ProductResponse `json:"products"`
	TotalCount int               `json:"total_count"`
	Page       int               `json:"page"`
	Limit      int               `json:"limit"`
}

type PageNation struct {
	// set default limit
	Limit int `query:"limit"  validate:"omitempty,gte=0"`
	Page  int `query:"page" validate:"omitempty,gte=1"`
}

func (p *PageNation) ApplyDefaults() {
	if p.Limit == 0 {
		p.Limit = 20
	}
}

type FilterParams struct {
	MinPrice float64 `query:"minPrice" validate:"gte=0"`
	MaxPrice float64 `query:"maxPrice" validate:"gte=0"`
	Category string  `query:"category" validate:"omitempty,min=3,max=100,alphanumspace"`
	Sort     string  `query:"sortField" validate:"omitempty"`
	PageNation
}

func (f *FilterParams) ApplyDefaults() {
	f.PageNation.ApplyDefaults()
	if f.MaxPrice == 0 {
		f.MaxPrice = math.MaxFloat64
	}
}

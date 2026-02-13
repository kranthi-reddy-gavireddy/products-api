package dtos

import (
	"math"
)

type SortClause struct {
	Field     string
	Direction string
}

type PageNation struct {
	Limit  int `query:"limit" validate:"gte=0"`
	Offset int `query:"offset" validate:"gte=0"`
}

type FilterParams struct {
	MinPrice float64 `query:"minPrice" validate:"gte=0"`
	MaxPrice float64 `query:"maxPrice" validate:"gte=0"`
	Category string  `query:"category" validate:"omitempty,min=3,max=100,alphanumspace"`
	Sort     string  `query:"sortField" validate:"omitempty"`
	PageNation
}

func (p *PageNation) ApplyDefaults() {
	if p.Limit == 0 {
		p.Limit = 20
	}
}

func (f *FilterParams) ApplyDefaults() {
	f.PageNation.ApplyDefaults()
	if f.MaxPrice == 0 {
		f.MaxPrice = math.MaxFloat64
	}
}

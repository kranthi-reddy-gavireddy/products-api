package models

type PageNation struct {
	Limit  int `query:"limit" validate:"gte=0"`
	Offset int `query:"offset" validate:"gte=0"`
}

type FilterParams struct {
	MinPrice float64 `query:"min_price" validate:"gte=0"`
	MaxPrice float64 `query:"max_price" validate:"gte=0"`
	Category string  `query:"category"`
	PageNation
}

func (p *PageNation) ApplyDefaults() {
	if p.Limit == 0 {
		p.Limit = 20
	}
}

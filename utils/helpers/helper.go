package helpers

import (
	"github.com/kranthi-reddy-gavireddy/products-api.git/dtos"

	"github.com/go-playground/validator/v10"
)

const (
	SortAsc  = "ASC"
	SortDesc = "DESC"
)

var (
	SortMapper = map[string]dtos.SortClause{
		"price": {
			Field:     "price",
			Direction: SortAsc,
		},
		"price:desc": {
			Field:     "price",
			Direction: SortDesc,
		},
		"date": {
			Field:     "created_at",
			Direction: SortAsc,
		},
		"date:desc": {
			Field:     "created_at",
			Direction: SortDesc,
		},
	}
)

func DataValidator[T *dtos.ProductRequest | *dtos.FilterParams | *dtos.OrderCreationListener](data T) error {

	return validator.New().Struct(data)
}

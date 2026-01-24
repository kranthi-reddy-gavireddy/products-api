package utils

import (
	"products-api/internal/models"

	"github.com/go-playground/validator/v10"
)

const (
	SortAsc  = "ASC"
	SortDesc = "DESC"
)

var (
	SortMapper = map[string]models.SortClause{
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

func DataValidator[T interface{}](data T) error {
	return validator.New().Struct(data)
}

package utils

import "github.com/go-playground/validator/v10"

func DataValidator[T interface{}](data T) error {
	return validator.New().Struct(data)
}

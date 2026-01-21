package listeners

import (
	"context"
	"encoding/json"
	"log"
	apperrors "products-api/app-errors"
	"products-api/internal/services"
	"products-api/utils"

	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
	"github.com/go-playground/validator/v10"
)

type OrderCreationListener struct {
	// Implementation details would go here
	ProductId string `json:"product_id" validate:"required"`
	Quantity  int    `json:"quantity" validate:"required,min=1,gt=0"`
}

func HandleOrderCreation(msg *types.Message, service *services.ProductService) error {
	var orderCreationListener OrderCreationListener
	log.Printf("OrderCreated Received message: %s", *msg.Body)
	err := json.Unmarshal([]byte(*msg.Body), &orderCreationListener)
	if err != nil {
		return err
	}
	if err := utils.DataValidator(&orderCreationListener); err != nil {
		log.Printf("Validation error: %v", err.(validator.ValidationErrors))
		return apperrors.ValidationError(err.(validator.ValidationErrors))
	}
	_, err = service.UpdateProductCount(context.Background(), orderCreationListener.ProductId, orderCreationListener.Quantity)
	if err != nil {
		return err
	}
	return nil
}

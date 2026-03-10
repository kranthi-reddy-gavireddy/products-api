package listeners

import (
	"context"
	"encoding/json"

	apperrors "github.com/kranthi-reddy-gavireddy/products-api.git/app-errors"
	"github.com/kranthi-reddy-gavireddy/products-api.git/dtos"
	"github.com/kranthi-reddy-gavireddy/products-api.git/internal/services"
	"github.com/kranthi-reddy-gavireddy/products-api.git/logger"

	"github.com/kranthi-reddy-gavireddy/products-api.git/utils/helpers"

	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
	"github.com/go-playground/validator/v10"
)

func HandleOrderCreation(msg *types.Message, service services.IProductService) error {
	var orderCreationListener dtos.OrderCreationListener

	logger := logger.New()

	logger.Infof("Received message for Order Created Topic %s", *msg.Body)

	err := json.Unmarshal([]byte(*msg.Body), &orderCreationListener)
	if err != nil {
		return err
	}

	if err := helpers.DataValidator(&orderCreationListener); err != nil {
		logger.Errorf("Validation error: %v", err.(validator.ValidationErrors))
		return apperrors.ValidationError(err.(validator.ValidationErrors))
	}

	_, err = service.UpdateCount(context.Background(), orderCreationListener.ProductId, orderCreationListener.Quantity)
	if err != nil {
		return err
	}
	return nil
}

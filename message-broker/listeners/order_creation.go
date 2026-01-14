package listeners

import (
	"context"
	"encoding/json"
	"log"
	"products-api/internal/services"

	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
)

type OrderCreationListener struct {
	// Implementation details would go here
	ProductId string `json:"product_id"`
	Quantity  int    `json:"quantity"`
}

func HandleOrderCreation(msg *types.Message, service *services.ProductService) error {
	var orderCreationListener OrderCreationListener
	log.Printf("OrderCreated Received message: %s", *msg.Body)
	err := json.Unmarshal([]byte(*msg.Body), &orderCreationListener)
	if err != nil {
		return err
	}
	_, err = service.UpdateProductCount(context.Background(), orderCreationListener.ProductId, orderCreationListener.Quantity)
	if err != nil {
		return err
	}
	return nil
}

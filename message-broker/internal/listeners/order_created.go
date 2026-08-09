package listeners

import "encoding/json"

func init() {
	Register("order_created", handleOrderCreated)
}

type OrderCreatedMessage struct {
	OrderID    string  `json:"order_id"`
	UserID     string  `json:"user_id"`
	TotalPrice float64 `json:"total_price"`
}

func handleOrderCreated(message []byte) error {
	var err error
	var orderCreatedMsg OrderCreatedMessage

	err = json.Unmarshal(message, &orderCreatedMsg)
	if err != nil {
		return err
	}

	// Process the order created message
	// ...

	return nil
}

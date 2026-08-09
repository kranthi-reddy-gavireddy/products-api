package listeners

import (
	"encoding/json"
	"log"
)

type NofitySellerMessage struct {
	ProductID string `json:"product_id"`
	SellerID  string `json:"seller_id"`
	Msg       string `json:"msg"`
}

func init() {
	Register("notify_seller", handleNotifySeller)
}

func handleNotifySeller(message []byte) error {

	var notifyMsg NofitySellerMessage
	err := json.Unmarshal(message, &notifyMsg)
	if err != nil {
		return err
	}

	log.Printf("Notify Seller: ProductID=%s, SellerID=%s, Msg=%s", notifyMsg.ProductID, notifyMsg.SellerID, notifyMsg.Msg)

	return nil
}

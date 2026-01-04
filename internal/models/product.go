package models

type Product struct {
	BaseModel
	Name     string  `json:"name"`
	Price    float64 `json:"price"`
	SellerID string  `json:"seller_id"`
	Quantity int     `json:"quantity"`
}

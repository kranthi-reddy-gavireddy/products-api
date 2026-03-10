package models

type Product struct {
	BaseModel
	//Name     string  `json:"name" validate:"required,min=3,max=100,matches=^[a-zA-Z0-9\\s]+$"`
	Name     string  `json:"name"`
	Price    float64 `json:"price"`
	SellerID string  `json:"seller_id"`
	Quantity int     `json:"quantity"`
	Category string  `json:"category"`
}

type ListProducts struct {
	Products   []Product `json:"products"`
	TotalCount int       `json:"total_count"`
}

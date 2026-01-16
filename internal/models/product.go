package models

type Product struct {
	BaseModel
	//Name     string  `json:"name" validate:"required,min=3,max=100,matches=^[a-zA-Z0-9\\s]+$"`
	Name     string  `json:"name" validate:"required,min=3,max=100,alphanumspace"`
	Price    float64 `json:"price" validate:"required,gt=0"`
	SellerID string  `json:"seller_id" validate:"required"`
	Quantity int     `json:"quantity" validate:"required,gt=0,lte=1000"`
}

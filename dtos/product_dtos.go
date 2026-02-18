package dtos

type ProductBase struct {
	Name     string  `json:"name" validate:"required,min=3,max=100,alphanumspace"`
	Price    float64 `json:"price" validate:"required,gt=0"`
	SellerID string  `json:"seller_id" validate:"required"`
	Quantity int     `json:"quantity" validate:"required,gt=0,lte=1000"`
	Category string  `json:"category" validate:"required,min=3,max=100,alphanumspace"`
}

type ProductRequest struct {
	ProductBase
}

type ProductResponse struct {
	ID string `json:"id"`
	ProductBase
}

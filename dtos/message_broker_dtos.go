package dtos

type OrderCreationListener struct {
	// Implementation details would go here
	ProductId string `json:"product_id" validate:"required"`
	Quantity  int    `json:"quantity" validate:"required,min=1,gt=0"`
}

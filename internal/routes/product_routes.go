package routes

import (
	"products-api/internal/handlers"
	"products-api/internal/server"
)

type ProductRoutes struct {
	hander handlers.ProductHandler
}

func NewProductRoutes(handler handlers.ProductHandler) *ProductRoutes {
	return &ProductRoutes{hander: handler}
}

func (r *ProductRoutes) RegisterRoutes(server *server.FiberServer) {

	server.App.Post("/products", r.hander.CreateProduct)
	server.App.Get("/products", r.hander.GetProducts)
}

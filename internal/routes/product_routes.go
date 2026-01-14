package routes

import (
	"products-api/internal/handlers"
	"products-api/internal/server"
)

type ProductRoutes struct {
	handler handlers.ProductHandler
}

func NewProductRoutes(handler handlers.ProductHandler) *ProductRoutes {
	return &ProductRoutes{handler: handler}
}

func (r *ProductRoutes) RegisterRoutes(server *server.FiberServer) {

	server.App.Post("/products", r.handler.CreateProduct)
	server.App.Get("/products", r.handler.GetProducts)
	server.App.Get("/products/:id", r.handler.GetByID)
	server.App.Delete("/products/:id", r.handler.Delete)
}

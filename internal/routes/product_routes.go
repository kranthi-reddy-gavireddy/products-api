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
	//I want to set /products as the base route for product-related endpoints
	group := server.App.Group("api/products")
	group.Post("", r.handler.CreateProduct)
	group.Get("", r.handler.GetProducts)
	group.Get("/id/:id", r.handler.GetByID)
	group.Delete("/id/:id", r.handler.Delete)
}

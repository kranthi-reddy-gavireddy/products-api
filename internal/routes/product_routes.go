package routes

import (
	"products-api/internal/handlers"
	"products-api/internal/server"
)

type ProductRoutes struct {
	handler handlers.IProductHandler
}

func (r *ProductRoutes) RegisterRoutes(server *server.FiberServer) {
	//I want to set /products as the base route for product-related endpoints
	group := server.App.Group("/api/products")
	group.Post("", r.handler.Create)
	group.Get("", r.handler.Filter)
	group.Get("/id/:id", r.handler.GetByID)
	group.Delete("/id/:id", r.handler.Delete)
	//group.Get("/search", r.handler.FilterProducts)

}

func New(handler handlers.IProductHandler) *ProductRoutes {
	return &ProductRoutes{handler: handler}
}

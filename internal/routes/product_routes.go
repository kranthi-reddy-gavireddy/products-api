package routes

import (
	"products-api/internal/handlers"

	"github.com/gofiber/fiber/v2"
)

type IProductRoutes interface {
	RegisterRoutes(engine fiber.Router)
}

type ProductRoutes struct {
	handler handlers.IProductHandler
}

func (r *ProductRoutes) RegisterRoutes(engine fiber.Router) {
	//I want to set /products as the base route for product-related endpoints
	group := engine.Group("/api/products")
	group.Post("", r.handler.Create)
	group.Put("/id/:id", r.handler.Update)
	group.Get("", r.handler.Filter)
	group.Get("/id/:id", r.handler.GetByID)
	group.Delete("/id/:id", r.handler.Delete)
	//group.Get("/search", r.handler.FilterProducts)

}

func New(handler handlers.IProductHandler) IProductRoutes {
	return &ProductRoutes{handler: handler}
}

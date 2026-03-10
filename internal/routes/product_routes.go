package routes

import (
	"github.com/kranthi-reddy-gavireddy/products-api.git/internal/handlers"
	"github.com/kranthi-reddy-gavireddy/products-api.git/middeware"

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
	group.Post("", middeware.CorrelationMiddleWare(middeware.RequestLogger(middeware.ErrorMiddleware(r.handler.Create))))
	group.Put("/id/:id", middeware.CorrelationMiddleWare(middeware.RequestLogger(middeware.ErrorMiddleware(r.handler.Update))))
	group.Get("", middeware.CorrelationMiddleWare(middeware.RequestLogger(middeware.ErrorMiddleware(r.handler.Filter))))
	group.Get("/id/:id", middeware.CorrelationMiddleWare(middeware.RequestLogger(middeware.ErrorMiddleware(r.handler.GetByID))))
	group.Delete("/id/:id", middeware.CorrelationMiddleWare(middeware.RequestLogger(middeware.ErrorMiddleware(r.handler.Delete))))
	//group.Get("/search", r.handler.FilterProducts)

}

func New(handler handlers.IProductHandler) IProductRoutes {
	return &ProductRoutes{handler: handler}
}

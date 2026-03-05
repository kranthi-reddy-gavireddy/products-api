package routes

import (
	"products-api/internal/database"
	"products-api/internal/handlers"
	"products-api/internal/repository"
	"products-api/internal/services"
	"products-api/middeware"
)

var handler handlers.IProductHandler

func RegisterRoutes() {

	dbInstance := database.New().GetDB()

	productRepo := repository.New(dbInstance)
	prodcutService := services.New(productRepo)
	handler = handlers.New(prodcutService)

}

var (
	Create  = middeware.CorrelationMiddleWare(middeware.RequestLogger(middeware.ErrorMiddleware(handler.Create)))
	Update  = middeware.CorrelationMiddleWare(middeware.RequestLogger(middeware.ErrorMiddleware(handler.Update)))
	GetAll  = middeware.CorrelationMiddleWare(middeware.RequestLogger(middeware.ErrorMiddleware(handler.Get)))
	GetByID = middeware.CorrelationMiddleWare(middeware.RequestLogger(middeware.ErrorMiddleware(handler.GetByID)))
	Delete  = middeware.CorrelationMiddleWare(middeware.RequestLogger(middeware.ErrorMiddleware(handler.Delete)))
)

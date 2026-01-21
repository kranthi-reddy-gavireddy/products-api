package middeware

import (
	"log"
	apperrors "products-api/app-errors"

	"github.com/gofiber/fiber/v2"
)

func ErrorMiddleware() func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		// Middleware logic to handle errors can be implemented here
		log.Printf("Processing request: Method %s OriginalURL %s Path %s Body %s IP %s", c.Method(), c.OriginalURL(), c.Path(), c.Body(), c.IP())
		//how to set correlation id in fiber context
		err := c.Next()
		if err == nil {
			return nil
		}
		if appError, ok := err.(*apperrors.AppError); ok {
			log.Printf("App error occurred: %v", appError)
			return c.Status(appError.HTTPStatus).JSON(appError)
		}
		log.Printf("Unhandled error occurred: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"code":    "INTERNAL_ERROR",
			"message": "Something went wrong",
		})
	}
}

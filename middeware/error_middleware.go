package middeware

import (
	apperrors "products-api/app-errors"
	"products-api/logger"

	"github.com/gofiber/fiber/v2"
)

func ErrorMiddleware() func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		err := c.Next()
		if err == nil {
			return nil
		}
		if appError, ok := err.(*apperrors.AppError); ok {
			logger.Errorf("App error occurred: %v", appError)
			return c.Status(appError.HTTPStatus).JSON(appError)
		}
		logger.Errorf("Unhandled error occurred: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"code":    "INTERNAL_ERROR",
			"message": "Something went wrong",
		})
	}
}

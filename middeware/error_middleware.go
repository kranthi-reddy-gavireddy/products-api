package middeware

import (
	apperrors "products-api/app-errors"
	"products-api/logger"
	ctx "products-api/utils/context"

	"github.com/gofiber/fiber/v2"
)

func ErrorMiddleware(next func(context *ctx.Context) error) func(c *ctx.Context) error {
	return func(c *ctx.Context) error {
		err := next(c)
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

package middeware

import (
	apperrors "github.com/kranthi-reddy-gavireddy/products-api.git/app-errors"
	ctx "github.com/kranthi-reddy-gavireddy/products-api.git/utils/context"

	"github.com/gofiber/fiber/v2"
)

func ErrorMiddleware(next func(context *ctx.Context) error) func(c *ctx.Context) error {
	return func(c *ctx.Context) error {
		err := next(c)
		if err == nil {
			return nil
		}
		if appError, ok := err.(*apperrors.AppError); ok {
			c.Logger.Errorf("App error occurred: %v", appError)
			return c.Status(appError.HTTPStatus).JSON(appError)
		}
		c.Logger.Errorf("Unhandled error occurred: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"code":    "INTERNAL_ERROR",
			"message": "Something went wrong",
		})
	}
}

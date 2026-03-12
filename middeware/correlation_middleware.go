package middeware

import (
	"github.com/kranthi-reddy-gavireddy/products-api.git/logger"
	ctx "github.com/kranthi-reddy-gavireddy/products-api.git/utils/context"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func CorrelationMiddleWare(next func(c *ctx.Context) error) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		correlationId := c.Get("X-Correlation-ID")
		if correlationId == "" {
			correlationId = uuid.NewString()
		}
		c.Locals("CorrelationID", correlationId)
		c.Set("X-Correlation-ID", correlationId)
		logger.WithCorrelationID(correlationId)
		context := ctx.Context{
			CorrelationID: correlationId,
			Logger:        logger.WithCorrelationID(correlationId),
			Ctx:           c,
		}
		c.Locals("Context", &context)
		return next(&context)
	}
}

package middeware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func CorrelationMiddleWare() func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		correlationId := c.Get("X-Correlation-ID")
		if correlationId == "" {
			correlationId = uuid.NewString()
		}
		c.Locals("CorrelationID", correlationId)
		c.Set("X-Correlation-ID", correlationId)
		return c.Next()
	}
}

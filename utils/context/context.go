package context

import (
	"products-api/logger"

	"github.com/gofiber/fiber/v2"
)

type Context struct {
	CorrelationID string
	Logger        logger.ZeroLogger
	*fiber.Ctx
}

func (c *Context) Copy() *Context {
	return &Context{
		CorrelationID: c.CorrelationID,
		Logger:        c.Logger,
		Ctx:           c.Ctx,
	}
}

package middeware

import (
	"products-api/logger"
	"time"

	"github.com/gofiber/fiber/v2"
)

type RequestLog struct {
	Method        string `json:"method"`
	Path          string `json:"path"`
	ReponseStatus int    `json:"status"`
	Latency       int    `json:"latency"`
}

func RequestLogger() func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		start := time.Now()
		err := c.Next()
		end := time.Since(start)
		requestLog := RequestLog{
			Method:        c.Method(),
			Path:          c.Path(),
			ReponseStatus: c.Response().StatusCode(),
			Latency:       int(end),
		}
		logger.Infof("Request Completed %v", &requestLog)
		return err
	}
}

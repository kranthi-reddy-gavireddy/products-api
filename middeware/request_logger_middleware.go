package middeware

import (
	"products-api/logger"
	ctx "products-api/utils/context"
	"time"
)

type RequestLog struct {
	Method        string `json:"method"`
	Path          string `json:"path"`
	ReponseStatus int    `json:"status"`
	Latency       int    `json:"latency"`
}

func RequestLogger(next func(app *ctx.Context) error) func(c *ctx.Context) error {
	return func(c *ctx.Context) error {
		start := time.Now()
		err := next(c)
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

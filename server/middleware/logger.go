package middleware

import (
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
	"ray-d-song.com/packer/utils"
)

func LoggerMiddleware() func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		start := time.Now()
		body := c.Body()
		queryParams := c.OriginalURL()

		contentType := c.Get("Content-Type")
		var requestBody string

		if strings.HasPrefix(contentType, "multipart/form-data") {
			requestBody = "multipart/form-data: [omitted]"
		} else {
			requestBody = string(body)
		}

		err := c.Next()

		utils.Logger.Info("Request",
			zap.String("method", c.Method()),
			zap.String("path", c.Path()),
			zap.Int("status", c.Response().StatusCode()),
			zap.Duration("latency", time.Since(start)),
			zap.String("ip", c.IP()),
			zap.String("user-agent", c.Get("User-Agent")),
			zap.String("request-body", requestBody),
			zap.String("query-params", queryParams),
		)

		utils.Logger.Info("Response",
			zap.String("body", string(c.Response().Body())),
		)

		return err
	}
}

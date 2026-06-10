package middleware

import (
	"time"

	"user-api/internal/logger"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

// RequestID injects a unique X-Request-ID header into every response.
// If the incoming request already carries this header, its value is reused.
func RequestID() fiber.Handler {
	return func(c *fiber.Ctx) error {
		reqID := c.Get("X-Request-ID")
		if reqID == "" {
			reqID = uuid.New().String()
		}
		c.Set("X-Request-ID", reqID)
		// Store in locals so downstream handlers can reference it.
		c.Locals("requestID", reqID)
		return c.Next()
	}
}

// RequestLogger logs the HTTP method, path, status code, and duration of
// every request using the global Uber Zap logger.
func RequestLogger() fiber.Handler {
	return func(c *fiber.Ctx) error {
		start := time.Now()
		err := c.Next()
		duration := time.Since(start)

		reqID, _ := c.Locals("requestID").(string)

		logger.Info("http request",
			zap.String("request_id", reqID),
			zap.String("method", c.Method()),
			zap.String("path", c.Path()),
			zap.Int("status", c.Response().StatusCode()),
			zap.Duration("duration", duration),
			zap.String("ip", c.IP()),
		)

		return err
	}
}

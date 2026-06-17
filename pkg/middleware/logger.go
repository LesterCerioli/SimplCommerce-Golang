package middleware

import (
	"log"
	"time"

	"github.com/gofiber/fiber/v3"
)

func LoggerMiddleware() fiber.Handler {
	return func(c fiber.Ctx) error {
		start := time.Now()
		path := c.Path()
		method := c.Method()

		err := c.Next()

		duration := time.Since(start)
		status := c.Response().StatusCode()

		log.Printf("[%s] %s %s %d %v",
			time.Now().Format(time.RFC3339),
			method,
			path,
			status,
			duration,
		)

		return err
	}
}

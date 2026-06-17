package middleware

import (
	"sync"
	"time"

	"github.com/gofiber/fiber/v3"
)

type rateLimiter struct {
	mu          sync.Mutex
	requests    map[string][]time.Time
	maxRequests int
	window      time.Duration
}

func RateLimit(maxRequests int, window time.Duration) fiber.Handler {
	rl := &rateLimiter{
		requests:    make(map[string][]time.Time),
		maxRequests: maxRequests,
		window:      window,
	}

	return func(c fiber.Ctx) error {
		ip := c.IP()

		rl.mu.Lock()
		now := time.Now()
		windowStart := now.Add(-rl.window)

		times := rl.requests[ip]
		var valid []time.Time
		for _, t := range times {
			if t.After(windowStart) {
				valid = append(valid, t)
			}
		}

		if len(valid) >= rl.maxRequests {
			rl.mu.Unlock()
			c.Status(fiber.StatusTooManyRequests)
			return c.JSON(fiber.Map{
				"error": "too many requests",
			})
		}

		rl.requests[ip] = append(valid, now)
		rl.mu.Unlock()

		return c.Next()
	}
}

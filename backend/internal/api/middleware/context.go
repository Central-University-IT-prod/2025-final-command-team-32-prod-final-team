package middleware

import (
	"context"
	"solution/pkg/logger"

	"github.com/gofiber/fiber/v2"
)

func CustomContext(ctx context.Context) fiber.Handler {
	return func(c *fiber.Ctx) error {
		c.SetUserContext(logger.CtxWithLogger(ctx))
		return c.Next()
	}
}

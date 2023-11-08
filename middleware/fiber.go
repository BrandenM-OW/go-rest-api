package middleware

import (
	"strings"
	"time"

	"github.com/form3tech-oss/jwt-go"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

func FiberMiddleware(a *fiber.App) {
	a.Use(
		cors.New(cors.Config{
			AllowOrigins: "http://127.0.0.1:8080, http://localhost:8080",
			AllowHeaders: "Origin, Content-Type, Accept",
			AllowMethods: "GET, POST, PATCH, DELETE",
		}),

		recover.New(),

		limiter.New(limiter.Config{
			Next: func(c *fiber.Ctx) bool {
				return strings.Contains(c.Path(), "/swagger") // don't limit swagger
			},
			Expiration: 10 * time.Second,
			Max:        300,
		}),

		logger.New(),
	)
}

func IsAdmin(c *fiber.Ctx) error {
	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)

	isAdmin, ok := claims["admin"]
	if !ok || !isAdmin.(bool) {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"msg": "Forbidden",
		})
	}

	return c.Next()
}

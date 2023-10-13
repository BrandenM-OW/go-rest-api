package initializers

import (
	"errors"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func CustomErrorHandler(c *fiber.Ctx, err error) error {
	code := fiber.StatusInternalServerError

	// Set appropriate status codes for gORM error types
	if errors.Is(err, gorm.ErrRecordNotFound) {
		code = fiber.StatusNotFound
	}

	if errors.Is(err, gorm.ErrInvalidTransaction) {
		code = fiber.StatusBadRequest
	}

	// Retrieve the custom status code if it's a *fiber.Error
	var e *fiber.Error
	if errors.As(err, &e) {
		code = e.Code
	}

	c.Set(fiber.HeaderContentType, fiber.MIMETextPlainCharsetUTF8)
	return c.Status(code).JSON(&fiber.Map{"error": err.Error()})
}

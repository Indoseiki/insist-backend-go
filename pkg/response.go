package pkg

import (
	"github.com/gofiber/fiber/v2"
)

func Response(c *fiber.Ctx, status int, message string, data interface{}) error {
	return c.Status(status).JSON(fiber.Map{
		"status":  status,
		"message": message,
		"data":    data,
	})
}

func ErrorResponse(c *fiber.Ctx, err error) error {
	var statusCode int
	var errorMessage string

	if fiberErr, ok := err.(*fiber.Error); ok {
		statusCode = fiberErr.Code
		errorMessage = fiberErr.Message
	}

	return Response(c, statusCode, errorMessage, nil)
}

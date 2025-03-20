package middleware

import (
	"insist-backend-golang/pkg"

	"github.com/gofiber/fiber/v2"
)

func VerifyToken(c *fiber.Ctx) error {
	token := c.Get("Authorization")
	if token == "" {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusForbidden, "Missing access token"))
	}

	if len(token) > 7 && token[:7] == "Bearer " {
		token = token[7:]
	} else {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusForbidden, "Invalid access token format"))
	}

	userID, err := pkg.VerifyAccessToken(token)
	if err != nil {
		if err.Error() == "token expired" {
			return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusForbidden, "Access token expired"))
		}
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusUnauthorized, "Invalid access token"))
	}

	c.Locals("userID", userID)

	return c.Next()
}

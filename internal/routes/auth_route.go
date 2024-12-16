package routes

import (
	"insist-backend-golang/internal/handler"
	"insist-backend-golang/internal/middleware"
	"insist-backend-golang/internal/service"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func AuthRoutes(api fiber.Router, db *gorm.DB) {
	authService := service.NewAuthService(db)
	passwordResetService := service.NewPasswordResetService(db)
	authHandler := handler.NewAuthHandler(authService, passwordResetService)

	api.Post("/login", authHandler.Login)
	api.Post("/two-fa", authHandler.TwoFactorAuth)
	api.Delete("/logout", authHandler.Logout)
	api.Get("/token", authHandler.RefreshToken)
	api.Get("/user-info", middleware.VerifyToken, authHandler.GetUserInfo)
	api.Put("/change-password", middleware.VerifyToken, authHandler.ChangePassword)
	api.Put("/:id/two-fa", middleware.VerifyToken, authHandler.SetTwoFactorAuth)
	api.Post("/:id/send-password-reset", middleware.VerifyToken, authHandler.SendPasswordReset)
	api.Post("/password-reset", authHandler.PasswordReset)
}

package routes

import (
	"insist-backend-golang/internal/handler"
	"insist-backend-golang/internal/service"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func UserRoutes(api fiber.Router, db *gorm.DB) {
	user := api.Group("master/users")

	userService := service.NewUserService(db)
	userHandler := handler.NewUserHandler(userService)

	user.Get("/", userHandler.GetUsers)
	user.Get("/:id", userHandler.GetUser)
	user.Post("/", userHandler.CreateUser)
	user.Put("/:id", userHandler.UpdateUser)
	user.Delete("/:id", userHandler.DeleteUser)
	user.Put("/:id/change-password", userHandler.ChangePassword)
}

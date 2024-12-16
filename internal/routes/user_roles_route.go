package routes

import (
	"insist-backend-golang/internal/handler"
	"insist-backend-golang/internal/service"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func UserRoleRoutes(api fiber.Router, db *gorm.DB) {
	userRole := api.Group("master/user-role")

	userRoleService := service.NewUserRoleService(db)
	userRoleHandler := handler.NewUserRoleHandler(userRoleService)

	userRole.Get("/", userRoleHandler.GetUserRoles)
	userRole.Get("/:id", userRoleHandler.GetUserRole)
	userRole.Put("/:id", userRoleHandler.UpdateUserRole)
}

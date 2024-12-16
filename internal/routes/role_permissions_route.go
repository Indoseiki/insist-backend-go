package routes

import (
	"insist-backend-golang/internal/handler"
	"insist-backend-golang/internal/service"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func RolePermissionRoutes(api fiber.Router, db *gorm.DB) {
	rolePermission := api.Group("role-permission")

	rolePermissionService := service.NewRolePermissionService(db)
	rolePermissionHandler := handler.NewRolePermissionHandler(rolePermissionService)

	rolePermission.Get("/:id", rolePermissionHandler.GetMenuTreeByRole)
	// rolePermission.Get("/:id", rolePermissionHandler.GetRolePermission)
	rolePermission.Post("/", rolePermissionHandler.UpdateOrCreateRolePermission)
	// rolePermission.Delete("/:id", rolePermissionHandler.DeleteRolePermission)
}

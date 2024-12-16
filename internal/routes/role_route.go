package routes

import (
	"insist-backend-golang/internal/handler"
	"insist-backend-golang/internal/service"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func RoleRoutes(api fiber.Router, db *gorm.DB) {
	role := api.Group("master/role")

	roleService := service.NewRoleService(db)
	roleHandler := handler.NewRoleHandler(roleService)

	role.Get("/", roleHandler.GetRoles)
	role.Get("/:id", roleHandler.GetRole)
	role.Post("/", roleHandler.CreateRole)
	role.Put("/:id", roleHandler.UpdateRole)
	role.Delete("/:id", roleHandler.DeleteRole)
}

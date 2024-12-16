package routes

import (
	"insist-backend-golang/internal/handler"
	"insist-backend-golang/internal/service"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func RoleMenuRoutes(api fiber.Router, db *gorm.DB) {
	roleMenu := api.Group("master/role-menu")

	roleMenuService := service.NewRoleMenuService(db)
	roleMenuHandler := handler.NewRoleMenuHandler(roleMenuService)

	roleMenu.Get("/", roleMenuHandler.GetRoleMenus)
	roleMenu.Get("/:id", roleMenuHandler.GetRoleMenu)
	roleMenu.Put("/:id", roleMenuHandler.UpdateRoleMenu)
}

package routes

import (
	"insist-backend-golang/internal/handler"
	"insist-backend-golang/internal/service"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func MenuRoutes(api fiber.Router, db *gorm.DB) {
	menu := api.Group("master/menu")

	menuService := service.NewMenuService(db)
	menuHandler := handler.NewMenuHandler(menuService)

	menu.Get("/", menuHandler.GetMenus)
	menu.Get("/:id", menuHandler.GetMenu)
	menu.Post("/", menuHandler.CreateMenu)
	menu.Put("/:id", menuHandler.UpdateMenu)
	menu.Delete("/:id", menuHandler.DeleteMenu)

	treeMenu := api.Group("master/tree-menu")
	treeMenu.Get("/", menuHandler.GetMenuTree)
	treeMenu.Get("/user", menuHandler.GetMenuTreeByUser)
}

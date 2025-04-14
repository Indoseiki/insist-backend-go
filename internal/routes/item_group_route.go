package routes

import (
	"insist-backend-golang/internal/handler"
	"insist-backend-golang/internal/service"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func ItemGroupRoutes(api fiber.Router, db *gorm.DB) {
	itemGroup := api.Group("master/item/group")

	itemGroupService := service.NewItemGroupService(db)
	itemGroupHandler := handler.NewItemGroupHandler(itemGroupService)

	itemGroup.Get("/", itemGroupHandler.GetItemGroups)
	itemGroup.Get("/:id", itemGroupHandler.GetItemGroup)
	itemGroup.Post("/", itemGroupHandler.CreateItemGroup)
	itemGroup.Put("/:id", itemGroupHandler.UpdateItemGroup)
	itemGroup.Delete("/:id", itemGroupHandler.DeleteItemGroup)
}

package routes

import (
	"insist-backend-golang/internal/handler"
	"insist-backend-golang/internal/service"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func ItemRoutes(api fiber.Router, db *gorm.DB) {
	item := api.Group("master/item/generate")

	itemService := service.NewItemService(db)
	itemHandler := handler.NewItemHandler(itemService)

	item.Get("/", itemHandler.GetItems)
	item.Get("/:id", itemHandler.GetItem)
	item.Post("/", itemHandler.CreateItem)
	item.Put("/:id", itemHandler.UpdateItem)
	item.Delete("/:id", itemHandler.DeleteItem)
}

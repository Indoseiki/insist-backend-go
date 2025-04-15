package routes

import (
	"insist-backend-golang/internal/handler"
	"insist-backend-golang/internal/service"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func ItemSourceRoutes(api fiber.Router, db *gorm.DB) {
	itemSource := api.Group("master/item/source")

	itemSourceService := service.NewItemSourceService(db)
	itemSourceHandler := handler.NewItemSourceHandler(itemSourceService)

	itemSource.Get("/", itemSourceHandler.GetItemSources)
	itemSource.Get("/:id", itemSourceHandler.GetItemSource)
	itemSource.Post("/", itemSourceHandler.CreateItemSource)
	itemSource.Put("/:id", itemSourceHandler.UpdateItemSource)
	itemSource.Delete("/:id", itemSourceHandler.DeleteItemSource)
}

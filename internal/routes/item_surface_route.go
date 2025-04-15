package routes

import (
	"insist-backend-golang/internal/handler"
	"insist-backend-golang/internal/service"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func ItemSurfaceRoutes(api fiber.Router, db *gorm.DB) {
	itemSurface := api.Group("master/item/surface")

	itemSurfaceService := service.NewItemSurfaceService(db)
	itemSurfaceHandler := handler.NewItemSurfaceHandler(itemSurfaceService)

	itemSurface.Get("/", itemSurfaceHandler.GetItemSurfaces)
	itemSurface.Get("/:id", itemSurfaceHandler.GetItemSurface)
	itemSurface.Post("/", itemSurfaceHandler.CreateItemSurface)
	itemSurface.Put("/:id", itemSurfaceHandler.UpdateItemSurface)
	itemSurface.Delete("/:id", itemSurfaceHandler.DeleteItemSurface)
}

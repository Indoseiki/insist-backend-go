package routes

import (
	"insist-backend-golang/internal/handler"
	"insist-backend-golang/internal/service"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func ItemProductTypeRoutes(api fiber.Router, db *gorm.DB) {
	itemProductType := api.Group("master/item/product-type")

	itemProductTypeService := service.NewItemProductTypeService(db)
	itemProductTypeHandler := handler.NewItemProductTypeHandler(itemProductTypeService)

	itemProductType.Get("/", itemProductTypeHandler.GetItemProductTypes)
	itemProductType.Get("/:id", itemProductTypeHandler.GetItemProductType)
	itemProductType.Post("/", itemProductTypeHandler.CreateItemProductType)
	itemProductType.Put("/:id", itemProductTypeHandler.UpdateItemProductType)
	itemProductType.Delete("/:id", itemProductTypeHandler.DeleteItemProductType)
}

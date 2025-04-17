package routes

import (
	"insist-backend-golang/internal/handler"
	"insist-backend-golang/internal/service"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func ItemRawMaterialRoutes(api fiber.Router, db *gorm.DB) {
	itemRawMaterial := api.Group("master/item/generate-raw-material")

	itemRawMaterialService := service.NewItemRawMaterialService(db)
	itemRawMaterialHandler := handler.NewItemRawMaterialHandler(itemRawMaterialService)

	itemRawMaterial.Get("/", itemRawMaterialHandler.GetItemRawMaterials)
	itemRawMaterial.Get("/:id", itemRawMaterialHandler.GetItemRawMaterial)
	itemRawMaterial.Post("/", itemRawMaterialHandler.CreateItemRawMaterial)
	itemRawMaterial.Put("/:id", itemRawMaterialHandler.UpdateItemRawMaterial)
	itemRawMaterial.Delete("/:id", itemRawMaterialHandler.DeleteItemRawMaterial)
}

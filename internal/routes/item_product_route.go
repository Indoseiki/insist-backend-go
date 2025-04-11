package routes

import (
	"insist-backend-golang/internal/handler"
	"insist-backend-golang/internal/service"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func ItemProductRoutes(api fiber.Router, db *gorm.DB) {
	itemProduct := api.Group("master/item/product")

	itemProductService := service.NewItemProductService(db)
	itemProductHandler := handler.NewItemProductHandler(itemProductService)

	itemProduct.Get("/", itemProductHandler.GetItemProducts)
	itemProduct.Get("/:id", itemProductHandler.GetItemProduct)
	itemProduct.Post("/", itemProductHandler.CreateItemProduct)
	itemProduct.Put("/:id", itemProductHandler.UpdateItemProduct)
	itemProduct.Delete("/:id", itemProductHandler.DeleteItemProduct)
}

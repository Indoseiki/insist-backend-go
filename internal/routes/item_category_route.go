package routes

import (
	"insist-backend-golang/internal/handler"
	"insist-backend-golang/internal/service"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func ItemCategoryRoutes(api fiber.Router, db *gorm.DB) {
	itemCategory := api.Group("master/item-category")

	itemCategoryService := service.NewItemCategoryService(db)
	itemCategoryHandler := handler.NewItemCategoryHandler(itemCategoryService)

	itemCategory.Get("/", itemCategoryHandler.GetItemCategories)
	itemCategory.Get("/:id", itemCategoryHandler.GetItemCategory)
	itemCategory.Post("/", itemCategoryHandler.CreateItemCategory)
	itemCategory.Put("/:id", itemCategoryHandler.UpdateItemCategory)
	itemCategory.Delete("/:id", itemCategoryHandler.DeleteItemCategory)
}

package routes

import (
	"insist-backend-golang/internal/handler"
	"insist-backend-golang/internal/service"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func ItemSubCategoryRoutes(api fiber.Router, db *gorm.DB) {
	itemSubCategory := api.Group("master/item/sub-category")

	itemSubCategoryService := service.NewItemSubCategoryService(db)
	itemSubCategoryHandler := handler.NewItemSubCategoryHandler(itemSubCategoryService)

	itemSubCategory.Get("/", itemSubCategoryHandler.GetItemSubCategories)
	itemSubCategory.Get("/:id", itemSubCategoryHandler.GetItemSubCategory)
	itemSubCategory.Post("/", itemSubCategoryHandler.CreateItemSubCategory)
	itemSubCategory.Put("/:id", itemSubCategoryHandler.UpdateItemSubCategory)
	itemSubCategory.Delete("/:id", itemSubCategoryHandler.DeleteItemSubCategory)
}

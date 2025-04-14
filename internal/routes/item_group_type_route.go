package routes

import (
	"insist-backend-golang/internal/handler"
	"insist-backend-golang/internal/service"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func ItemGroupTypeRoutes(api fiber.Router, db *gorm.DB) {
	itemGroupType := api.Group("master/item/group-type")

	itemGroupTypeService := service.NewItemGroupTypeService(db)
	itemGroupTypeHandler := handler.NewItemGroupTypeHandler(itemGroupTypeService)

	itemGroupType.Get("/", itemGroupTypeHandler.GetItemGroupTypes)
	itemGroupType.Get("/:id", itemGroupTypeHandler.GetItemGroupType)
	itemGroupType.Post("/", itemGroupTypeHandler.CreateItemGroupType)
	itemGroupType.Put("/:id", itemGroupTypeHandler.UpdateItemGroupType)
	itemGroupType.Delete("/:id", itemGroupTypeHandler.DeleteItemGroupType)
}

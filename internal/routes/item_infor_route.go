package routes

import (
	"insist-backend-golang/internal/handler"
	"insist-backend-golang/internal/service"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func ItemInforRoutes(api fiber.Router, db *gorm.DB) {
	itemInfor := api.Group("master/item")

	itemInforService := service.NewItemInforService(db)
	itemInforHandler := handler.NewItemInforHandler(itemInforService)

	itemInfor.Get("/", itemInforHandler.GetItemInfors)
}

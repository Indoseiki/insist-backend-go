package routes

import (
	"insist-backend-golang/internal/handler"
	"insist-backend-golang/internal/service"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func ItemProcessRoutes(api fiber.Router, db *gorm.DB) {
	itemProcess := api.Group("master/item/process")

	itemProcessService := service.NewItemProcessService(db)
	itemProcessHandler := handler.NewItemProcessHandler(itemProcessService)

	itemProcess.Get("/", itemProcessHandler.GetItemProcesses)
	itemProcess.Get("/:id", itemProcessHandler.GetItemProcess)
	itemProcess.Post("/", itemProcessHandler.CreateItemProcess)
	itemProcess.Put("/:id", itemProcessHandler.UpdateItemProcess)
	itemProcess.Delete("/:id", itemProcessHandler.DeleteItemProcess)
}

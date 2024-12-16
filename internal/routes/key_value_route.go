package routes

import (
	"insist-backend-golang/internal/handler"
	"insist-backend-golang/internal/service"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func KeyValueRoutes(api fiber.Router, db *gorm.DB) {
	keyValue := api.Group("master/key-value")

	keyValueService := service.NewKeyValueService(db)
	keyValueHandler := handler.NewKeyValueHandler(keyValueService)

	keyValue.Get("/", keyValueHandler.GetKeyValues)
	keyValue.Get("/:id", keyValueHandler.GetKeyValue)
	keyValue.Post("/", keyValueHandler.CreateKeyValue)
	keyValue.Put("/:id", keyValueHandler.UpdateKeyValue)
	keyValue.Delete("/:id", keyValueHandler.DeleteKeyValue)
}

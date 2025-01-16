package routes

import (
	"insist-backend-golang/internal/handler"
	"insist-backend-golang/internal/service"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func LocationRoutes(api fiber.Router, db *gorm.DB) {
	location := api.Group("master/location")

	locationService := service.NewLocationService(db)
	locationHandler := handler.NewLocationHandler(locationService)

	location.Get("/", locationHandler.GetLocations)
	location.Get("/:id", locationHandler.GetLocation)
	location.Post("/", locationHandler.CreateLocation)
	location.Put("/:id", locationHandler.UpdateLocation)
	location.Delete("/:id", locationHandler.DeleteLocation)
}

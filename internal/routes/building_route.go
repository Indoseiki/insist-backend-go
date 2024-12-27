package routes

import (
	"insist-backend-golang/internal/handler"
	"insist-backend-golang/internal/service"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func BuildingRoutes(api fiber.Router, db *gorm.DB) {
	building := api.Group("master/building")

	buildingService := service.NewBuildingService(db)
	buildingHandler := handler.NewBuildingHandler(buildingService)

	building.Get("/", buildingHandler.GetBuildings)
	building.Get("/:id", buildingHandler.GetBuilding)
	building.Post("/", buildingHandler.CreateBuilding)
	building.Put("/:id", buildingHandler.UpdateBuilding)
	building.Delete("/:id", buildingHandler.DeleteBuilding)
}

package routes

import (
	"insist-backend-golang/internal/handler"
	"insist-backend-golang/internal/service"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func FCSBuildingRoutes(api fiber.Router, db *gorm.DB) {
	fcsBuilding := api.Group("master/fcs-building")

	fcsBuildingService := service.NewFCSBuildingService(db)
	fcsBuildingHandler := handler.NewFCSBuildingHandler(fcsBuildingService)

	fcsBuilding.Get("/", fcsBuildingHandler.GetFCSBuildings)
	fcsBuilding.Get("/:id", fcsBuildingHandler.GetFCSBuilding)
	fcsBuilding.Put("/:id", fcsBuildingHandler.UpdateFCSBuilding)
}

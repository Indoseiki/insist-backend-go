package routes

import (
	"insist-backend-golang/internal/handler"
	"insist-backend-golang/internal/service"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func FCSRoutes(api fiber.Router, db *gorm.DB) {
	fcs := api.Group("master/fcs")

	fcsService := service.NewFCSService(db)
	fcsHandler := handler.NewFCSHandler(fcsService)

	fcs.Get("/", fcsHandler.GetFCSs)
	fcs.Get("/:id", fcsHandler.GetFCS)
	fcs.Post("/", fcsHandler.CreateFCS)
	fcs.Put("/:id", fcsHandler.UpdateFCS)
	fcs.Delete("/:id", fcsHandler.DeleteFCS)
}

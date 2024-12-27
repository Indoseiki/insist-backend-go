package routes

import (
	"insist-backend-golang/internal/handler"
	"insist-backend-golang/internal/service"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func SubSectionRoutes(api fiber.Router, db *gorm.DB) {
	subSection := api.Group("master/sub-section")

	subSectionService := service.NewSubSectionService(db)
	subSectionHandler := handler.NewSubSectionHandler(subSectionService)

	subSection.Get("/", subSectionHandler.GetSubSections)
	subSection.Get("/:id", subSectionHandler.GetSubSection)
	subSection.Post("/", subSectionHandler.CreateSubSection)
	subSection.Put("/:id", subSectionHandler.UpdateSubSection)
	subSection.Delete("/:id", subSectionHandler.DeleteSubSection)
}

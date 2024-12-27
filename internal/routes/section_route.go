package routes

import (
	"insist-backend-golang/internal/handler"
	"insist-backend-golang/internal/service"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func SectionRoutes(api fiber.Router, db *gorm.DB) {
	section := api.Group("master/section")

	sectionService := service.NewSectionService(db)
	sectionHandler := handler.NewSectionHandler(sectionService)

	section.Get("/", sectionHandler.GetSections)
	section.Get("/:id", sectionHandler.GetSection)
	section.Post("/", sectionHandler.CreateSection)
	section.Put("/:id", sectionHandler.UpdateSection)
	section.Delete("/:id", sectionHandler.DeleteSection)
}

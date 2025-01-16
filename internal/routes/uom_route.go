package routes

import (
	"insist-backend-golang/internal/handler"
	"insist-backend-golang/internal/service"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func UoMRoutes(api fiber.Router, db *gorm.DB) {
	uom := api.Group("master/uom")

	uomService := service.NewUoMService(db)
	uomHandler := handler.NewUoMHandler(uomService)

	uom.Get("/", uomHandler.GetUoMs)
	uom.Get("/:id", uomHandler.GetUoM)
	uom.Post("/", uomHandler.CreateUoM)
	uom.Put("/:id", uomHandler.UpdateUoM)
	uom.Delete("/:id", uomHandler.DeleteUoM)
}

package routes

import (
	"insist-backend-golang/internal/handler"
	"insist-backend-golang/internal/service"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func TaxCodeRoutes(api fiber.Router, db *gorm.DB) {
	taxCode := api.Group("master/tax-code")

	taxCodeService := service.NewTaxCodeService(db)
	taxCodeHandler := handler.NewTaxCodeHandler(taxCodeService)

	taxCode.Get("/", taxCodeHandler.GetTaxCodes)
	taxCode.Get("/:id", taxCodeHandler.GetTaxCode)
	taxCode.Post("/", taxCodeHandler.CreateTaxCode)
	taxCode.Put("/:id", taxCodeHandler.UpdateTaxCode)
	taxCode.Delete("/:id", taxCodeHandler.DeleteTaxCode)
}

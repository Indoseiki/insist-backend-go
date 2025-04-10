package routes

import (
	"insist-backend-golang/internal/handler"
	"insist-backend-golang/internal/service"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func BillingTermRoutes(api fiber.Router, db *gorm.DB) {
	billingTerm := api.Group("master/billing-term")

	billingTermService := service.NewBillingTermService(db)
	billingTermHandler := handler.NewBillingTermHandler(billingTermService)

	billingTerm.Get("/", billingTermHandler.GetBillingTerms)
	billingTerm.Get("/:id", billingTermHandler.GetBillingTerm)
	billingTerm.Post("/", billingTermHandler.CreateBillingTerm)
	billingTerm.Put("/:id", billingTermHandler.UpdateBillingTerm)
	billingTerm.Delete("/:id", billingTermHandler.DeleteBillingTerm)
}

package routes

import (
	"insist-backend-golang/internal/handler"
	"insist-backend-golang/internal/service"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func CurrencyRateRoutes(api fiber.Router, db *gorm.DB) {
	currencyRate := api.Group("master/currency-rate")

	currencyRateService := service.NewCurrencyRateService(db)
	currencyRateHandler := handler.NewCurrencyRateHandler(currencyRateService)

	currencyRate.Get("/", currencyRateHandler.GetCurrencyRates)
	currencyRate.Get("/:id", currencyRateHandler.GetCurrencyRate)
	currencyRate.Post("/", currencyRateHandler.CreateCurrencyRate)
	currencyRate.Put("/:id", currencyRateHandler.UpdateCurrencyRate)
	currencyRate.Delete("/:id", currencyRateHandler.DeleteCurrencyRate)
}

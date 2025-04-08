package routes

import (
	"insist-backend-golang/internal/handler"
	"insist-backend-golang/internal/service"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func CurrencyRoutes(api fiber.Router, db *gorm.DB) {
	currency := api.Group("master/currency")
	generate := api.Group("master/currency-generate")

	currencyService := service.NewCurrencyService(db)
	currencyHandler := handler.NewCurrencyHandler(currencyService)

	currency.Get("/", currencyHandler.GetCurrencies)
	currency.Get("/:id", currencyHandler.GetCurrency)
	currency.Post("/", currencyHandler.CreateCurrency)
	currency.Put("/:id", currencyHandler.UpdateCurrency)
	currency.Delete("/:id", currencyHandler.DeleteCurrency)
	generate.Post("/", currencyHandler.GenerateCurrency)
}

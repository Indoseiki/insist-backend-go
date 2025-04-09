package routes

import (
	"insist-backend-golang/internal/handler"
	"insist-backend-golang/internal/service"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func BankRoutes(api fiber.Router, db *gorm.DB) {
	bank := api.Group("master/bank")

	bankService := service.NewBankService(db)
	bankHandler := handler.NewBankHandler(bankService)

	bank.Get("/", bankHandler.GetBanks)
	bank.Get("/:id", bankHandler.GetBank)
	bank.Post("/", bankHandler.CreateBank)
	bank.Put("/:id", bankHandler.UpdateBank)
	bank.Delete("/:id", bankHandler.DeleteBank)
}

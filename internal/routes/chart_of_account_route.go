package routes

import (
	"insist-backend-golang/internal/handler"
	"insist-backend-golang/internal/service"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func ChartOfAccountRoutes(api fiber.Router, db *gorm.DB) {
	chartOfAccount := api.Group("master/chart-of-account")

	chartOfAccountService := service.NewChartOfAccountService(db)
	chartOfAccountHandler := handler.NewChartOfAccountHandler(chartOfAccountService)

	chartOfAccount.Get("/", chartOfAccountHandler.GetChartOfAccounts)
	chartOfAccount.Get("/:id", chartOfAccountHandler.GetChartOfAccount)
	chartOfAccount.Post("/", chartOfAccountHandler.CreateChartOfAccount)
	chartOfAccount.Put("/:id", chartOfAccountHandler.UpdateChartOfAccount)
	chartOfAccount.Delete("/:id", chartOfAccountHandler.DeleteChartOfAccount)
}

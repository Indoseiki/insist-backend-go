package routes

import (
	"insist-backend-golang/internal/handler"
	"insist-backend-golang/internal/service"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func ProcessRoutes(api fiber.Router, db *gorm.DB) {
	process := api.Group("master/process")

	processService := service.NewProcessService(db)
	processHandler := handler.NewProcessHandler(processService)

	process.Get("/", processHandler.GetProcesses)
	process.Get("/:id", processHandler.GetProcess)
	process.Post("/", processHandler.CreateProcess)
	process.Put("/:id", processHandler.UpdateProcess)
	process.Delete("/:id", processHandler.DeleteProcess)
}

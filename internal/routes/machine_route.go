package routes

import (
	"insist-backend-golang/internal/handler"
	"insist-backend-golang/internal/service"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func MachineRoutes(api fiber.Router, db *gorm.DB) {
	machine := api.Group("master/machine")

	machineService := service.NewMachineService(db)
	machineHandler := handler.NewMachineHandler(machineService)

	machine.Get("/", machineHandler.GetMachines)
	machine.Get("/:id", machineHandler.GetMachine)
	machine.Post("/", machineHandler.CreateMachine)
	machine.Put("/:id", machineHandler.UpdateMachine)
	machine.Delete("/:id", machineHandler.DeleteMachine)
	machine.Put("/:id/revision", machineHandler.RevisionMachine)
	machine.Get("/:id/detail", machineHandler.GetMachineDetails)
	machine.Get("/:id/status", machineHandler.GetMachineStatus)
	machine.Put("/:id/status", machineHandler.CreateStatusMachine)
}

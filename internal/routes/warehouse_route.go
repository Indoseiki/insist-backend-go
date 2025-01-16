package routes

import (
	"insist-backend-golang/internal/handler"
	"insist-backend-golang/internal/service"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func WarehouseRoutes(api fiber.Router, db *gorm.DB) {
	warehouse := api.Group("master/warehouse")

	warehouseService := service.NewWarehouseService(db)
	warehouseHandler := handler.NewWarehouseHandler(warehouseService)

	warehouse.Get("/", warehouseHandler.GetWarehouses)
	warehouse.Get("/:id", warehouseHandler.GetWarehouse)
	warehouse.Post("/", warehouseHandler.CreateWarehouse)
	warehouse.Put("/:id", warehouseHandler.UpdateWarehouse)
	warehouse.Delete("/:id", warehouseHandler.DeleteWarehouse)
}

package routes

import (
	"insist-backend-golang/internal/handler"
	"insist-backend-golang/internal/service"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func MaterialRoutes(api fiber.Router, db *gorm.DB) {
	material := api.Group("master/material")

	materialService := service.NewMaterialService(db)
	materialHandler := handler.NewMaterialHandler(materialService)

	material.Get("/", materialHandler.GetMaterials)
	material.Get("/:id", materialHandler.GetMaterial)
	material.Post("/", materialHandler.CreateMaterial)
	material.Put("/:id", materialHandler.UpdateMaterial)
	material.Delete("/:id", materialHandler.DeleteMaterial)
}

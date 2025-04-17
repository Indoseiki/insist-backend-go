package routes

import (
	"insist-backend-golang/internal/handler"
	"insist-backend-golang/internal/service"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func MaterialDetailRoutes(api fiber.Router, db *gorm.DB) {
	materialDetail := api.Group("master/material-detail")

	materialDetailService := service.NewMaterialDetailService(db)
	materialDetailHandler := handler.NewMaterialDetailHandler(materialDetailService)

	materialDetail.Get("/", materialDetailHandler.GetMaterialDetails)
	materialDetail.Get("/:id", materialDetailHandler.GetMaterialDetail)
	materialDetail.Post("/", materialDetailHandler.CreateMaterialDetail)
	materialDetail.Put("/:id", materialDetailHandler.UpdateMaterialDetail)
	materialDetail.Delete("/:id", materialDetailHandler.DeleteMaterialDetail)
}

package routes

import (
	"insist-backend-golang/internal/handler"
	"insist-backend-golang/internal/service"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func ReasonRoutes(api fiber.Router, db *gorm.DB) {
	reason := api.Group("master/reason")

	reasonService := service.NewReasonService(db)
	reasonHandler := handler.NewReasonHandler(reasonService)

	reason.Get("/", reasonHandler.GetReasons)
	reason.Get("/:id", reasonHandler.GetReason)
	reason.Post("/", reasonHandler.CreateReason)
	reason.Put("/:id", reasonHandler.UpdateReason)
	reason.Delete("/:id", reasonHandler.DeleteReason)
}

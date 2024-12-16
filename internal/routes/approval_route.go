package routes

import (
	"insist-backend-golang/internal/handler"
	"insist-backend-golang/internal/service"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func ApprovalRoutes(api fiber.Router, db *gorm.DB) {
	approval := api.Group("approval")

	approvalService := service.NewApprovalService(db)
	approvalHandler := handler.NewApprovalHandler(approvalService)

	approval.Get("/", approvalHandler.GetApprovals)
	approval.Get("/:id", approvalHandler.GetApproval)
	approval.Post("/", approvalHandler.CreateApproval)
	approval.Put("/", approvalHandler.UpdateApproval)
	approval.Delete("/:id", approvalHandler.DeleteApproval)
}

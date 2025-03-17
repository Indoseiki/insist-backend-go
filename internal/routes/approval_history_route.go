package routes

import (
	"insist-backend-golang/internal/handler"
	"insist-backend-golang/internal/service"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func ApprovalHistoryRoutes(api fiber.Router, db *gorm.DB) {
	approvalHistory := api.Group("approval-history")

	approvalHistoryService := service.NewApprovalHistoryService(db)
	approvalHistoryHandler := handler.NewApprovalHistoryHandler(approvalHistoryService)

	approvalHistory.Get("/", approvalHistoryHandler.GetApprovalHistories)
	approvalHistory.Get("/:id", approvalHistoryHandler.GetApprovalHistory)
	approvalHistory.Post("/", approvalHistoryHandler.CreateApprovalHistory)
	approvalHistory.Put("/:id", approvalHistoryHandler.UpdateApprovalHistory)
	approvalHistory.Delete("/:id", approvalHistoryHandler.DeleteApprovalHistory)
}

package routes

import (
	"insist-backend-golang/internal/handler"
	"insist-backend-golang/internal/service"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func ApprovalUserRoutes(api fiber.Router, db *gorm.DB) {
	approvalUser := api.Group("approval-user")

	approvalUserService := service.NewApprovalUserService(db)
	approvalUserHandler := handler.NewApprovalUserHandler(approvalUserService)

	approvalUser.Get("/", approvalUserHandler.GetApprovalUsers)
	approvalUser.Get("/:id/approval", approvalUserHandler.GetApprovalUsersByIdApproval)
	approvalUser.Put("/:id", approvalUserHandler.UpdateApprovalUser)
	approvalUser.Delete("/:id", approvalUserHandler.DeleteApprovalUser)
}

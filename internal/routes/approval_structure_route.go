package routes

import (
	"insist-backend-golang/internal/handler"
	"insist-backend-golang/internal/service"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func ApprovalStructureRoutes(api fiber.Router, db *gorm.DB) {
	approvalStructure := api.Group("approval-structure")

	approvalStructureService := service.NewApprovalStructureService(db)
	approvalStructureHandler := handler.NewApprovalStructureHandler(approvalStructureService)

	approvalStructure.Get("/", approvalStructureHandler.GetApprovalStructures)
	approvalStructure.Get("/:id", approvalStructureHandler.GetApprovalStructure)
}

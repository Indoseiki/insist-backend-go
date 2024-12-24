package routes

import (
	"insist-backend-golang/internal/handler"
	"insist-backend-golang/internal/service"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func ActivityLogRoutes(api fiber.Router, db *gorm.DB) {
	activityLog := api.Group("log")

	activityLogService := service.NewActivityLogService(db)
	activityLogHandler := handler.NewActivityLogHandler(activityLogService)

	activityLog.Get("/", activityLogHandler.GetActivityLogs)
	activityLog.Get("/:id", activityLogHandler.GetActivityLog)
	activityLog.Post("/", activityLogHandler.CreateActivityLog)
}

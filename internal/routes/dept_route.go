package routes

import (
	"insist-backend-golang/internal/handler"
	"insist-backend-golang/internal/service"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func DeptRoutes(api fiber.Router, db *gorm.DB) {
	dept := api.Group("master/department")

	deptService := service.NewDeptService(db)
	deptHandler := handler.NewDeptHandler(deptService)

	dept.Get("/", deptHandler.GetDepts)
	dept.Get("/:id", deptHandler.GetDept)
	dept.Post("/", deptHandler.CreateDept)
	dept.Put("/:id", deptHandler.UpdateDept)
	dept.Delete("/:id", deptHandler.DeleteDept)
}

package routes

import (
	"insist-backend-golang/internal/cron"
	"insist-backend-golang/internal/handler"
	"insist-backend-golang/internal/service"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/valyala/fasthttp"
	"gorm.io/gorm"
)

func EmployeeRoutes(api fiber.Router, db *gorm.DB) {
	employee := api.Group("master/employee")

	employeeService := service.NewEmployeeService(db)
	employeeHandler := handler.NewEmployeeHandler(employeeService)

	employee.Get("/", employeeHandler.GetEmployees)
	employee.Get("/:number", employeeHandler.GetEmployee)
	employee.Post("/sync", employeeHandler.SyncEmployee)

	cron.SetupCron(func() {
		ctx := fiber.New().AcquireCtx(&fasthttp.RequestCtx{})
		defer fiber.New().ReleaseCtx(ctx)

		if err := employeeHandler.SyncEmployee(ctx); err != nil {
			log.Println("Error executing SyncEmployee:", err)
		}
	}, "10 17 * * *")
}

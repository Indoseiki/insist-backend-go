// @title           INSIST REST API Documentation
// @version         1.0
// @description     INSIST adalah sistem manajemen manufaktur milik PT. Indoseiki Metalutama. REST API Documentation untuk INSIST menjelaskan bagaimana pengembang dapat menggunakan API ini untuk mengelola dan mengakses data manufaktur perusahaan dengan aman dan efisien, menggunakan protokol HTTP dan metode REST seperti GET, POST, PUT, & DELETE.
// @host            localhost:5050/api
// @BasePath        /
package main

import (
	_ "insist-backend-golang/docs"
	"insist-backend-golang/internal/config"
	"insist-backend-golang/internal/middleware"
	"insist-backend-golang/internal/routes"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/joho/godotenv"
	fiberSwagger "github.com/swaggo/fiber-swagger"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	config.ConnectDBINSIST()
	config.ConnectDBINFOR()

	app := fiber.New()
	app.Use(cors.New(cors.Config{
		AllowOrigins:     os.Getenv("CORS"),
		AllowMethods:     "GET,POST,PUT,DELETE",
		AllowHeaders:     "Origin, Content-Type, Accept, Authorization",
		AllowCredentials: true,
	}))

	app.Get("/swagger/*", fiberSwagger.WrapHandler)

	api := app.Group("/api")

	// Auth Routes
	apiAuth := api.Group("/auth")
	routes.AuthRoutes(apiAuth, config.DBINSIST)

	// ADM Routes
	apiADM := api.Group("/admin", middleware.VerifyToken)
	routes.UserRoutes(apiADM, config.DBINSIST)
	routes.DeptRoutes(apiADM, config.DBINSIST)
	routes.RoleRoutes(apiADM, config.DBINSIST)
	routes.MenuRoutes(apiADM, config.DBINSIST)
	routes.EmployeeRoutes(apiADM, config.DBINSIST)
	routes.ReasonRoutes(apiADM, config.DBINSIST)
	routes.KeyValueRoutes(apiADM, config.DBINSIST)
	routes.UserRoleRoutes(apiADM, config.DBINSIST)
	routes.RoleMenuRoutes(apiADM, config.DBINSIST)
	routes.RolePermissionRoutes(apiADM, config.DBINSIST)
	routes.ApprovalRoutes(apiADM, config.DBINSIST)
	routes.ApprovalUserRoutes(apiADM, config.DBINSIST)
	routes.ApprovalStructureRoutes(apiADM, config.DBINSIST)

	// EGD Routes
	apiEGD := api.Group("/egd", middleware.VerifyToken)
	routes.ProcessRoutes(apiEGD, config.DBINSIST)
	routes.UoMRoutes(apiEGD, config.DBINSIST)

	// MNT Routes
	apiMNT := api.Group("/mnt")
	routes.MachineRoutes(apiMNT, config.DBINSIST)

	// PID Routes
	apiPID := api.Group("/pid", middleware.VerifyToken)
	routes.WarehouseRoutes(apiPID, config.DBINSIST)
	routes.LocationRoutes(apiPID, config.DBINSIST)

	// PRD Routes
	apiPRD := api.Group("/prd", middleware.VerifyToken)
	routes.BuildingRoutes(apiPRD, config.DBINSIST)
	routes.FCSRoutes(apiPRD, config.DBINSIST)
	routes.SectionRoutes(apiPRD, config.DBINSIST)
	routes.SubSectionRoutes(apiPRD, config.DBINSIST)
	routes.FCSBuildingRoutes(apiPRD, config.DBINSIST)

	// Log Route
	routes.ActivityLogRoutes(api, config.DBINSIST)

	println("Starting app with port " + os.Getenv("PORT"))
	log.Fatal(app.Listen(":" + os.Getenv("PORT")))
}

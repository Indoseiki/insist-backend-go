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
		AllowOrigins:     "http://localhost:5173",
		AllowMethods:     "GET,POST,PUT,DELETE",
		AllowHeaders:     "Origin, Content-Type, Accept, Authorization",
		AllowCredentials: true,
	}))

	app.Get("/swagger/*", fiberSwagger.WrapHandler)

	api := app.Group("/api")

	// Auth Routes
	apiAuth := api.Group("/auth")
	routes.AuthRoutes(apiAuth, config.DBINSIST)

	// Admin Routes
	apiAdmin := api.Group("/admin", middleware.VerifyToken)
	routes.UserRoutes(apiAdmin, config.DBINSIST)
	routes.DeptRoutes(apiAdmin, config.DBINSIST)
	routes.RoleRoutes(apiAdmin, config.DBINSIST)
	routes.MenuRoutes(apiAdmin, config.DBINSIST)
	routes.EmployeeRoutes(apiAdmin, config.DBINSIST)
	routes.ReasonRoutes(apiAdmin, config.DBINSIST)
	routes.KeyValueRoutes(apiAdmin, config.DBINSIST)
	routes.UserRoleRoutes(apiAdmin, config.DBINSIST)
	routes.RoleMenuRoutes(apiAdmin, config.DBINSIST)
	routes.RolePermissionRoutes(apiAdmin, config.DBINSIST)
	routes.ApprovalRoutes(apiAdmin, config.DBINSIST)
	routes.ApprovalUserRoutes(apiAdmin, config.DBINSIST)
	routes.ApprovalStructureRoutes(apiAdmin, config.DBINSIST)

	println("Starting app with port " + os.Getenv("PORT"))
	log.Fatal(app.Listen(":" + os.Getenv("PORT")))
}

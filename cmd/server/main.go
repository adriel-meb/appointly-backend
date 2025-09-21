package main

import (
	"github.com/adriel-meb/appointly-backend/internal/config"
	"github.com/adriel-meb/appointly-backend/internal/controllers"
	"github.com/adriel-meb/appointly-backend/internal/db"
	"github.com/adriel-meb/appointly-backend/internal/middleware"
	"github.com/gin-gonic/gin"
	"log"
)

func init() {
	config.LoadEnvVariables()
	db.DbConnect()
	db.DbMigration()

}

func main() {
	router := gin.Default()
	gin.DebugPrintRouteFunc = func(httpMethod, absolutePath, handlerName string, nuHandlers int) {
		log.Printf("endpoint %v %v %v %v", httpMethod, absolutePath, handlerName, absolutePath)
	}

	router.GET("/", controllers.GetWelcome)
	router.POST("/auth/register", controllers.Signup)
	router.POST("/auth/login", controllers.Login)

	router.GET("/users", controllers.GetAllUsers)
	router.DELETE("/users/:email", controllers.DeleteUser)

	router.GET("/validate", middleware.RequireAuthMiddleware(), controllers.Validate)

	// Provider routes
	providers := router.Group("/providers")
	{
		providers.GET("/", controllers.GetAllProviders)
		providers.GET("/:id", controllers.GetProviderByID)
		providers.POST("/", controllers.CreateProvider)
		providers.PUT("/:id", middleware.RequireAuthMiddleware(), controllers.UpdateProvider)
		providers.DELETE("/:id", middleware.RequireAuthMiddleware(), controllers.DeleteProvider)
	}

	// Specialization routes - ADD AUTHORIZATION LATER
	specializations := router.Group("/specializations")
	{
		specializations.GET("/", controllers.GetAllSpecializations)
		specializations.POST("/", controllers.CreateSpecialization)
		specializations.PUT("/:id", controllers.UpdateSpecialization)
		specializations.DELETE("/:id", controllers.DeleteSpecialization)
	}

	// Service route
	services := router.Group("/services")
	{
		services.POST("/", controllers.CreateService)
		services.GET("/", controllers.GetAllServices)
		services.GET("/:id", controllers.GetServiceByID)
		services.PUT("/", controllers.UpdateServices)
		services.DELETE("/", controllers.DeleteServices)
	}

	// Availability routes
	availabilities := router.Group("/availabilities")
	{
		// 1️⃣ Create a new availability
		// POST /availabilities/
		availabilities.POST("/", controllers.CreateAvailability)

		// 2️⃣ Get all availabilities with optional filters
		// GET /availabilities/?provider_id=3&date=2025-09-20&start_date=2025-09-20&end_date=2025-09-30
		availabilities.GET("/", controllers.GetAllAvailability)

		// 3️⃣ Get a specific availability by ID
		// GET /availabilities/:id
		availabilities.GET("/:id", controllers.GetAvailabilityByID)

		// 4️⃣ Update an existing availability
		// PUT /availabilities/:id
		availabilities.PUT("/:id", controllers.UpdateAvailability)

		// 5️⃣ Delete an availability
		// DELETE /availabilities/:id
		availabilities.DELETE("/:id", controllers.DeleteAvailability)
	}

	bookings := router.Group("/bookings")
	{
		bookings.POST("/", controllers.CreateBooking)
		bookings.GET("/", controllers.GetAllBooking)
		bookings.POST("/confirm", controllers.ConfirmBooking)
	}

	// Start the server
	router.Run()
}

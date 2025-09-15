package main

import (
	"github.com/adriel-meb/appointly-backend/internal/config"
	"github.com/adriel-meb/appointly-backend/internal/controllers"
	"github.com/adriel-meb/appointly-backend/internal/db"
	"github.com/adriel-meb/appointly-backend/internal/middleware"
	"github.com/gin-gonic/gin"
)

func init() {
	config.LoadEnvVariables()
	db.DbConnect()
	db.DbMigration()

}

func main() {
	router := gin.Default()

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
		services.PUT("/", controllers.UpdateServices)
		services.DELETE("/", controllers.DeleteServices)
	}

	// Availability route
	availabilities := router.Group("/availabilities")
	{
		availabilities.POST("/", controllers.CreateService)
		//availabilities.GET("/", controllers.GetAllServices)
		//availabilities.PUT("/", controllers.UpdateServices)
		//availabilities.DELETE("/", controllers.DeleteServices)
	}

	// Start the server
	router.Run()
}

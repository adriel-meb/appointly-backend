package controllers

import (
	"github.com/adriel-meb/appointly-backend/internal/db"
	"github.com/adriel-meb/appointly-backend/internal/models"
	"github.com/gin-gonic/gin"
	"net/http"
)

func CreateService(c *gin.Context) {
	// 1. Input validation struct
	type CreateServiceInput struct {
		Title           string  `json:"title" binding:"required,min=3,max=100"`
		Description     string  `json:"description" binding:"omitempty,max=500"`
		DurationMinutes int     `json:"duration_minutes" binding:"required,gt=0"`
		Price           float64 `json:"price" binding:"required,gt=0"`
		ProviderID      uint    `json:"provider_id" binding:"required"`
	}

	var input CreateServiceInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, APIResponse{
			Status:  "error",
			Message: "Invalid input",
			Error:   err.Error(),
		})
		return
	}

	// 2. Check if provider exists
	var provider models.Provider
	if err := db.DB.First(&provider, input.ProviderID).Error; err != nil {
		c.JSON(http.StatusNotFound, APIResponse{
			Status:  "error",
			Message: "Provider not found",
		})
		return
	}

	// 3. Check if provider already offers this service
	var existingService models.Service
	if err := db.DB.Where("provider_id = ? AND title = ?", input.ProviderID, input.Title).First(&existingService).Error; err == nil {
		c.JSON(http.StatusBadRequest, APIResponse{
			Status:  "error",
			Message: "This service already exists for the provider",
		})
		return
	}

	// 4. Create new service
	service := models.Service{
		Title:           input.Title,
		Description:     input.Description,
		DurationMinutes: uint(input.DurationMinutes),
		Price:           float32(input.Price),
		ProviderID:      input.ProviderID,
	}

	if err := db.DB.Create(&service).Error; err != nil {
		c.JSON(http.StatusInternalServerError, APIResponse{
			Status:  "error",
			Message: "Failed to create service",
			Error:   err.Error(),
		})
		return
	}

	// 5. Return success
	c.JSON(http.StatusCreated, APIResponse{
		Status:  "success",
		Message: "Service created successfully",
		Data:    service,
	})
}

func GetAllServices(c *gin.Context) {
	var services []models.Service

	if err := db.DB.Preload("Provider").Find(&services).Error; err != nil {
		c.JSON(http.StatusInternalServerError, APIResponse{
			Status:  "error",
			Message: "Failed to fetch services",
			Error:   err.Error(),
		})
		return
	}

	// Success response
	c.JSON(http.StatusOK, APIResponse{
		Status:  "success",
		Message: "Services fetched successfully",
		Data:    services,
	})
}

// DeleteProvider placeholder
func DeleteServices(c *gin.Context) {
	c.JSON(http.StatusOK, APIResponse{
		Status:  "success",
		Message: "DeleteProvider endpoint - to be implemented",
	})
}

// UpdateProvider placeholder
func UpdateServices(c *gin.Context) {
	c.JSON(http.StatusOK, APIResponse{
		Status:  "success",
		Message: "UpdateProvider endpoint - to be implemented",
	})
}

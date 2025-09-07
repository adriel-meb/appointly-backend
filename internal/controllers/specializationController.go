package controllers

import (
	"github.com/adriel-meb/appointly-backend/internal/db"
	"github.com/adriel-meb/appointly-backend/internal/models"
	"github.com/gin-gonic/gin"
	"net/http"
)

// Create a specialization
func CreateSpecialization(c *gin.Context) {
	type CreateSpecializationInput struct {
		Name        string `json:"name" binding:"required"`
		Description string `json:"description,omitempty"`
	}

	var input CreateSpecializationInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, APIResponse{
			Status:  "error",
			Message: "Invalid input",
			Error:   err.Error(),
		})
		return
	}

	// Check if specialization exists (case-insensitive)
	var existing models.Specialization
	if err := db.DB.Where("LOWER(name) = LOWER(?)", input.Name).First(&existing).Error; err == nil {
		c.JSON(http.StatusBadRequest, APIResponse{
			Status: "error",
			Error:  "Specialization already registered",
		})
		return
	}

	specialization := models.Specialization{
		Name:        input.Name,
		Description: input.Description,
	}

	if err := db.DB.Create(&specialization).Error; err != nil {
		c.JSON(http.StatusInternalServerError, APIResponse{
			Status:  "error",
			Message: "Failed to create specialization",
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, APIResponse{
		Status:  "success",
		Message: "Specialization created successfully",
		Data:    specialization,
	})
}

// Get all specializations
func GetAllSpecializations(c *gin.Context) {
	var specializations []models.Specialization

	if err := db.DB.Find(&specializations).Error; err != nil {
		c.JSON(http.StatusInternalServerError, APIResponse{
			Status:  "error",
			Message: "Failed to fetch specializations",
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, APIResponse{
		Status:  "success",
		Message: "Specializations fetched successfully",
		Data:    specializations,
	})
}

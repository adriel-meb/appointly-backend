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

// UpdateSpecialization handles PUT /specializations/:id
func UpdateSpecialization(c *gin.Context) {
	// Get specialization ID from URL
	id := c.Param("id")

	// Define input
	type UpdateSpecializationInput struct {
		Name        string `json:"name,omitempty"`
		Description string `json:"description,omitempty"`
	}

	var input UpdateSpecializationInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, APIResponse{
			Status:  "error",
			Message: "Invalid input",
			Error:   err.Error(),
		})
		return
	}

	// Find specialization
	var specialization models.Specialization
	if err := db.DB.First(&specialization, id).Error; err != nil {
		c.JSON(http.StatusNotFound, APIResponse{
			Status:  "error",
			Message: "Specialization not found",
			Error:   err.Error(),
		})
		return
	}

	// Update fields
	if input.Name != "" {
		specialization.Name = input.Name
	}
	if input.Description != "" {
		specialization.Description = input.Description
	}

	// Save
	if err := db.DB.Save(&specialization).Error; err != nil {
		c.JSON(http.StatusInternalServerError, APIResponse{
			Status:  "error",
			Message: "Failed to update specialization",
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, APIResponse{
		Status:  "success",
		Message: "Specialization updated successfully",
		Data:    specialization,
	})
}

// DeleteSpecialization handles DELETE /specializations/:id
func DeleteSpecialization(c *gin.Context) {
	id := c.Param("id")

	// Find specialization
	var specialization models.Specialization
	if err := db.DB.First(&specialization, id).Error; err != nil {
		c.JSON(http.StatusNotFound, APIResponse{
			Status:  "error",
			Message: "Specialization not found",
			Error:   err.Error(),
		})
		return
	}

	// Delete
	if err := db.DB.Delete(&specialization).Error; err != nil {
		c.JSON(http.StatusInternalServerError, APIResponse{
			Status:  "error",
			Message: "Failed to delete specialization",
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, APIResponse{
		Status:  "success",
		Message: "Specialization deleted successfully",
	})
}

package controllers

import (
	"net/http"

	"github.com/adriel-meb/appointly-backend/internal/db"
	"github.com/adriel-meb/appointly-backend/internal/models"
	"github.com/gin-gonic/gin"
)

// CreateProvider handles POST /providers (admin only)
func CreateProvider(c *gin.Context) {
	type CreateProviderInput struct {
		SpecializationID uint   `json:"specialization_id" binding:"required"` // FK to specializations
		Bio              string `json:"bio"`
		UserID           uint   `json:"user_id" binding:"required"`
	}

	var input CreateProviderInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, APIResponse{
			Status:  "error",
			Message: "Invalid input",
			Error:   err.Error(),
		})
		return
	}

	// Check if user exists
	var user models.User
	if result := db.DB.First(&user, input.UserID); result.RowsAffected == 0 {
		c.JSON(http.StatusBadRequest, APIResponse{
			Status:  "error",
			Message: "User not found",
		})
		return
	}

	// Check if already a provider
	var existingProvider models.Provider
	if db.DB.Where("user_id = ?", input.UserID).First(&existingProvider); db.DB.RowsAffected > 0 {
		c.JSON(http.StatusBadRequest, APIResponse{
			Status:  "error",
			Message: "User is already a provider",
		})
		return
	}

	// Check user role
	if user.Role != models.RoleProvider {
		c.JSON(http.StatusBadRequest, APIResponse{
			Status:  "error",
			Message: "User role must be 'provider' to create a provider profile",
		})
		return
	}

	// Optional: validate specialization exists
	var specialization models.Specialization
	if err := db.DB.First(&specialization, input.SpecializationID).Error; err != nil {
		c.JSON(http.StatusBadRequest, APIResponse{
			Status:  "error",
			Message: "Specialization not found",
		})
		return
	}

	// Create provider
	provider := models.Provider{
		UserID:           input.UserID,
		SpecializationID: input.SpecializationID,
		Bio:              input.Bio,
	}

	if err := db.DB.Create(&provider).Error; err != nil {
		c.JSON(http.StatusInternalServerError, APIResponse{
			Status:  "error",
			Message: "Failed to create provider",
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, APIResponse{
		Status:  "success",
		Message: "Provider created successfully",
		Data:    provider,
	})
}

// GetAllProviders handles GET /providers
// It retrieves all providers from the database including their associated users.
// Returns a standardized APIResponse with the list of providers or an error message.
func GetAllProviders(c *gin.Context) {
	var providers []models.Provider

	// Fetch providers with preloaded user and specialization data
	if err := db.DB.Preload("User").Preload("Specialization").Find(&providers).Error; err != nil {
		c.JSON(http.StatusInternalServerError, APIResponse{
			Status:  "error",
			Message: "Failed to fetch providers",
			Error:   err.Error(),
		})
		return
	}

	// Success response
	c.JSON(http.StatusOK, APIResponse{
		Status:  "success",
		Message: "Providers fetched successfully",
		Data:    providers,
	})
}

// GetProviderByID placeholder
func GetProviderByID(c *gin.Context) {
	c.JSON(http.StatusOK, APIResponse{
		Status:  "success",
		Message: "GetProviderByID endpoint - to be implemented",
	})
}

// DeleteProvider placeholder
func DeleteProvider(c *gin.Context) {
	c.JSON(http.StatusOK, APIResponse{
		Status:  "success",
		Message: "DeleteProvider endpoint - to be implemented",
	})
}

// UpdateProvider placeholder
func UpdateProvider(c *gin.Context) {
	c.JSON(http.StatusOK, APIResponse{
		Status:  "success",
		Message: "UpdateProvider endpoint - to be implemented",
	})
}

package controllers

import (
	"net/http"

	"github.com/adriel-meb/appointly-backend/internal/db"
	"github.com/adriel-meb/appointly-backend/internal/models"
	"github.com/gin-gonic/gin"
)

// Create a new provider (admin only)
func CreateProvider(c *gin.Context) {

	//Validate the input (specialization, bio, user_id)
	type CreateProviderInput struct {
		Specialization string `json:"specialization" binding:"required"`
		Bio            string `json:"bio"`
		UserID         uint   `json:"user_id" binding:"required"`
	}

	var input CreateProviderInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	//check if the user exists

	var user models.User
	result := db.DB.First(&user, input.UserID)
	if result.RowsAffected == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User not found"})
		return
	}

	// verify that the user is not already a provider
	var existingProvider models.Provider
	db.DB.Where("user_id = ?", input.UserID).First(&existingProvider)
	if db.DB.RowsAffected > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User is already a provider"})
		return
	}

	//verify that the user role is 'provider'
	if user.Role != models.RoleProvider {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User role must be 'provider' to create a provider profile"})
		return
	}

	// Create the provider
	provider := models.Provider{
		Specialization: input.Specialization,
		Bio:            input.Bio,
		UserID:         input.UserID,
	}

	if err := db.DB.Create(&provider).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Respond with the created provider (without sensitive info)
	c.JSON(http.StatusCreated, gin.H{"message": "provider created", "provider": provider})

}

// retrieve all providers
func GetAllProviders(c *gin.Context) {
	var providers []models.Provider
	result := db.DB.Preload("User").Find(&providers)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch users"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"providers": providers,
	})
}

// retrieve a specific provider by ID
func GetProviderByID(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "GetProviderByID endpoint - to be implemented",
	})
}

// delete a provider by ID (admin only)
func DeleteProvider(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "DeleteProvider endpoint - to be implemented",
	})
}

// update a provider by ID (admin only)
func UpdateProvider(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "UpdateProvider endpoint - to be implemented",
	})
}

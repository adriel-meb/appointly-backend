package controllers

import (
	"log"
	"net/http"

	"github.com/adriel-meb/appointly-backend/internal/db"
	"github.com/adriel-meb/appointly-backend/internal/models"
	"github.com/gin-gonic/gin"
)

// APIResponse is the standard response format for all endpoints
type APIResponse struct {
	Status  string      `json:"status"`            // success or error
	Message string      `json:"message,omitempty"` // human-readable message
	Data    interface{} `json:"data,omitempty"`    // returned data (if any)
	Error   string      `json:"error,omitempty"`   // error details (hidden from clients in real prod)
	Version string      `json:"version,omitempty"` // optional API version
}

// ---------------------- ROUTES ---------------------- //

// GetWelcome -> GET /
// A simple health check / welcome endpoint
func GetWelcome(c *gin.Context) {
	c.JSON(http.StatusOK, APIResponse{
		Status:  "success",
		Message: "Welcome to Appointly API",
		Version: "v1.0.0",
		Data: gin.H{
			"endpoints": gin.H{
				"auth":      []string{"/auth/register", "/auth/login"},
				"users":     []string{"/users"},
				"providers": []string{"/providers"},
			},
		},
	})
}

// GetAllUsers -> GET /users
// Fetch all users from the database
func GetAllUsers(c *gin.Context) {
	var users []models.User

	// Try to fetch all users
	if err := db.DB.Find(&users).Error; err != nil {
		log.Printf("DB error fetching users: %v", err) // log only on server side
		c.JSON(http.StatusInternalServerError, APIResponse{
			Status: "error",
			Error:  "Failed to fetch users", // generic error to client
		})
		return
	}

	// Return the list of users
	c.JSON(http.StatusOK, APIResponse{
		Status: "success",
		Data:   users,
	})
}

// DeleteUser -> DELETE /users/:email
// Delete a user by email
func DeleteUser(c *gin.Context) {
	email := c.Param("email")

	// Validate param
	if email == "" {
		c.JSON(http.StatusBadRequest, APIResponse{
			Status: "error",
			Error:  "Email parameter is required",
		})
		return
	}

	// Try to delete user
	result := db.DB.Where("email = ?", email).Delete(&models.User{})

	if result.Error != nil {
		log.Printf("DB error deleting user %s: %v", email, result.Error)
		c.JSON(http.StatusInternalServerError, APIResponse{
			Status: "error",
			Error:  "Failed to delete user",
		})
		return
	}

	// No user found
	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, APIResponse{
			Status: "error",
			Error:  "User not found",
		})
		return
	}

	// Success
	c.JSON(http.StatusOK, APIResponse{
		Status:  "success",
		Message: "User deleted successfully",
	})
}

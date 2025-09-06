package controllers

import (
	"github.com/adriel-meb/appointly-backend/internal/db"
	"github.com/adriel-meb/appointly-backend/internal/models"
	"github.com/gin-gonic/gin"
)

// handles GET /users request
func GetUsersfunc(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "WELCOME TO APPOINTLY API",
	})
}

// handles GET /users request to fetch all users
func GetAllUsers(c *gin.Context) {
	var users []models.User
	result := db.DB.Find(&users)

	if result.Error != nil {
		c.JSON(500, gin.H{"error": "Failed to fetch users"})
		return
	}

	c.JSON(200, gin.H{"users": users})
}

// handles PUT /users/:id request to update a user

// handles DELETE /users/:email request to delete a user
func DeleteUser(c *gin.Context) {

	email := c.Param("email")
	result := db.DB.Where("email = ?", email).Delete(&models.User{})

	if result.Error != nil {
		c.JSON(500, gin.H{"message": "Failed to delete user", "error": result.Error.Error()})
		return
	}

	if result.RowsAffected == 0 {
		c.JSON(404, gin.H{"error": "User not found"})
		return
	}

	c.JSON(200, gin.H{"message": "User deleted successfully"})
}

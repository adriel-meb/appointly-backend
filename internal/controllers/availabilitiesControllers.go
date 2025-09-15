package controllers

import (
	"github.com/adriel-meb/appointly-backend/internal/db"
	"github.com/adriel-meb/appointly-backend/internal/models"
	"github.com/adriel-meb/appointly-backend/scripts"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

// CreateAvailability handles POST /availability
// This function allows a provider to define their available time slots.
// Steps:
// 1. Bind and validate the incoming JSON request (provider ID, day/date, time range).
// 2. Check that the given provider exists in the database.
// 3. (Future implementation) Save the availability record (recurring or one-time) into the DB.
// 4. Return a success or error response.
func CreateAvailability(c *gin.Context) {

	// Define the expected request input with validation rules
	type AvailabilityInput struct {
		ProviderID uint `json:"provider_id" binding:"required"`

		// For recurring weekly slots (optional, e.g., "Monday")
		DayOfWeek   *string `json:"day_of_week"`
		IsRecurring bool    `json:"is_recurring" binding:"required"`

		// For one-time specific slots (optional, e.g., "2025-09-20")
		Date *time.Time `json:"date"`

		// Required time range
		StartTime string `json:"start_time" binding:"required"` // "09:00"
		EndTime   string `json:"end_time" binding:"required"`   // "17:00"
	}

	var input AvailabilityInput

	// Step 1: Bind and validate the JSON input
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, APIResponse{
			Status:  "error",
			Message: "Invalid input",
			Error:   err.Error(),
		})
		return
	}

	// Step 2: Check if the provider exists in the database
	var provider models.Provider
	if result := db.DB.First(&provider, input.ProviderID); result.RowsAffected == 0 {
		c.JSON(http.StatusBadRequest, APIResponse{
			Status:  "error",
			Message: "Provider not found",
		})
		return
	}

	// 🚧 Step 3: Not yet implemented
	// Here you will construct a `models.Availability` object based on whether
	// it's recurring (DayOfWeek + Time range) or one-time (Date + Time range).

	if err := scripts.ValidateAvailabilityInput(input); err != nil {
		c.JSON(http.StatusBadRequest, APIResponse{
			Status:  "error",
			Message: err.Error(),
		})
		return
	}

	// Then save it in the database with db.DB.Create(&availability).
	// After saving, return a success response with the created availability.

	// Example placeholder response:
	c.JSON(http.StatusOK, APIResponse{
		Status:  "success",
		Message: "Availability validated (saving logic not yet implemented)",
		Data:    input,
	})
}

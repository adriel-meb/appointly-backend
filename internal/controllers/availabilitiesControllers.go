package controllers

import (
	"github.com/adriel-meb/appointly-backend/internal/db"
	"github.com/adriel-meb/appointly-backend/internal/models"
	"github.com/adriel-meb/appointly-backend/scripts"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"time"
)

// CreateAvailability handles POST /availability
// Allows a provider to define their available time slots (recurring or one-time).
func CreateAvailability(c *gin.Context) {

	// Local input struct
	type AvailabilityInput struct {
		ProviderID uint `json:"provider_id" binding:"required"`

		DayOfWeek   *string    `json:"day_of_week"`
		IsRecurring bool       `json:"is_recurring" binding:"required"`
		Date        *time.Time `json:"date"`
		StartTime   string     `json:"start_time" binding:"required"`
		EndTime     string     `json:"end_time" binding:"required"`
	}

	var input AvailabilityInput

	// Step 1: Bind JSON input
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, APIResponse{
			Status:  "error",
			Message: "Invalid input",
			Error:   err.Error(),
		})
		return
	}

	// Step 2: Check provider exists
	var provider models.Provider
	if result := db.DB.First(&provider, input.ProviderID); result.RowsAffected == 0 {
		c.JSON(http.StatusBadRequest, APIResponse{
			Status:  "error",
			Message: "Provider not found",
		})
		return
	}

	// Step 3: Convert local input to models.Availability for validation
	availabilityForValidation := models.Availability{
		DayOfWeek:   input.DayOfWeek,
		IsRecurring: input.IsRecurring,
		Date:        input.Date,
	}

	// Step 4: Validate recurring vs one-time rules
	if err := scripts.ValidateAvailabilityInput(availabilityForValidation); err != nil {
		c.JSON(http.StatusBadRequest, APIResponse{
			Status:  "error",
			Message: err.Error(),
		})
		return
	}

	// Step 5: Validate StartTime < EndTime
	start, err1 := time.Parse("15:04", input.StartTime)
	end, err2 := time.Parse("15:04", input.EndTime)
	if err1 != nil || err2 != nil {
		c.JSON(http.StatusBadRequest, APIResponse{
			Status:  "error",
			Message: "StartTime or EndTime is not valid. Format must be HH:MM",
		})
		return
	}
	if !end.After(start) {
		c.JSON(http.StatusBadRequest, APIResponse{
			Status:  "error",
			Message: "EndTime must be after StartTime",
		})
		return
	}

	// Step 6: Create availability record
	availability := models.Availability{
		ProviderID:  input.ProviderID,
		DayOfWeek:   input.DayOfWeek,
		IsRecurring: input.IsRecurring,
		Date:        input.Date,
		StartTime:   input.StartTime,
		EndTime:     input.EndTime,
	}

	if err := db.DB.Create(&availability).Error; err != nil {
		c.JSON(http.StatusInternalServerError, APIResponse{
			Status:  "error",
			Message: "Failed to create availability",
			Error:   err.Error(),
		})
		return
	}

	// Step 7: Return success response
	c.JSON(http.StatusCreated, APIResponse{
		Status:  "success",
		Message: "Availability created successfully",
		Data:    availability,
	})
}

// GetAllAvailability handles GET /availability
//
// This endpoint fetches provider availabilities with optional filtering.
//
// Filters (query parameters):
// - provider_id: integer → fetch availabilities for a specific provider
// - date: YYYY-MM-DD → fetch availabilities for a specific date
// - start_date: YYYY-MM-DD → fetch availabilities on or after this date
// - end_date: YYYY-MM-DD → fetch availabilities on or before this date
//
// Examples:
// 1️⃣ Fetch all availabilities:
//
//	GET /availability
//
// 2️⃣ Fetch availabilities for provider 3:
//
//	GET /availability?provider_id=3
//
// 3️⃣ Fetch availabilities on 2025-09-20:
//
//	GET /availability?date=2025-09-20
//
// 4️⃣ Fetch availabilities from 2025-09-20 to 2025-09-30:
//
//	GET /availability?start_date=2025-09-20&end_date=2025-09-30
//
// 5️⃣ Combine provider and date filters:
//
//	GET /availability?provider_id=3&start_date=2025-09-20&end_date=2025-09-30
func GetAllAvailability(c *gin.Context) {
	var availabilities []models.Availability

	// Start building the DB query
	dbQuery := db.DB.Preload("Provider")

	// -----------------------------
	// 1️⃣ Optional filter: provider_id
	// -----------------------------
	if providerID := c.Query("provider_id"); providerID != "" {
		dbQuery = dbQuery.Where("provider_id = ?", providerID)
	}

	// -----------------------------
	// 2️⃣ Optional filter: exact date
	// -----------------------------
	if dateStr := c.Query("date"); dateStr != "" {
		date, err := time.Parse("2006-01-02", dateStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, APIResponse{
				Status:  "error",
				Message: "Invalid date format. Use YYYY-MM-DD.",
			})
			return
		}
		dbQuery = dbQuery.Where("date = ?", date)
	}

	// -----------------------------
	// 3️⃣ Optional filter: start_date (>=)
	// -----------------------------
	if startStr := c.Query("start_date"); startStr != "" {
		start, err := time.Parse("2006-01-02", startStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, APIResponse{
				Status:  "error",
				Message: "Invalid start_date format. Use YYYY-MM-DD.",
			})
			return
		}
		dbQuery = dbQuery.Where("date >= ?", start)
	}

	// -----------------------------
	// 4️⃣ Optional filter: end_date (<=)
	// -----------------------------
	if endStr := c.Query("end_date"); endStr != "" {
		end, err := time.Parse("2006-01-02", endStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, APIResponse{
				Status:  "error",
				Message: "Invalid end_date format. Use YYYY-MM-DD.",
			})
			return
		}
		dbQuery = dbQuery.Where("date <= ?", end)
	}

	// -----------------------------
	// 5️⃣ Execute the query
	// -----------------------------
	if err := dbQuery.Find(&availabilities).Error; err != nil {
		c.JSON(http.StatusInternalServerError, APIResponse{
			Status:  "error",
			Message: "Failed to fetch availabilities",
			Error:   err.Error(),
		})
		return
	}

	// -----------------------------
	// 6️⃣ Success response
	// -----------------------------
	c.JSON(http.StatusOK, APIResponse{
		Status:  "success",
		Message: "Availabilities fetched successfully",
		Data:    availabilities,
	})
}

// GetAvailabilityByID handles GET /availabilities/:id
//
// Fetch a single availability by its ID, including the related Provider info.
// Steps:
// 1️⃣ Parse the availability ID from the URL.
// 2️⃣ Fetch the availability from the database with the Provider preloaded.
// 3️⃣ Return a 404 error if not found.
// 4️⃣ Return the availability in a success response.
func GetAvailabilityByID(c *gin.Context) {
	// Step 1: Parse ID from URL
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, APIResponse{
			Status:  "error",
			Message: "Invalid availability ID",
			Error:   err.Error(),
		})
		return
	}

	// Step 2: Fetch availability
	var availability models.Availability
	if err := db.DB.Preload("Provider").First(&availability, id).Error; err != nil {
		c.JSON(http.StatusNotFound, APIResponse{
			Status:  "error",
			Message: "Availability not found",
			Error:   err.Error(),
		})
		return
	}

	// Step 3: Success response
	c.JSON(http.StatusOK, APIResponse{
		Status:  "success",
		Message: "Availability fetched successfully",
		Data:    availability,
	})
}

// UpdateAvailability handles PUT /availabilities/:id
//
// Update an existing availability.
// Steps:
// 1️⃣ Parse ID from URL and check existence.
// 2️⃣ Bind and validate the JSON input (same structure as CreateAvailability).
// 3️⃣ Run ValidateAvailabilityInput to ensure DayOfWeek/Date consistency.
// 4️⃣ Update fields in the database.
// 5️⃣ Return updated availability.
func UpdateAvailability(c *gin.Context) {
	// Parse ID from URL
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, APIResponse{
			Status:  "error",
			Message: "Invalid availability ID",
			Error:   err.Error(),
		})
		return
	}

	// Fetch existing availability
	var availability models.Availability
	if err := db.DB.First(&availability, id).Error; err != nil {
		c.JSON(http.StatusNotFound, APIResponse{
			Status:  "error",
			Message: "Availability not found",
			Error:   err.Error(),
		})
		return
	}

	// Bind input JSON
	type AvailabilityInput struct {
		ProviderID  uint       `json:"provider_id" binding:"required"`
		DayOfWeek   *string    `json:"day_of_week"`
		IsRecurring bool       `json:"is_recurring" binding:"required"`
		Date        *time.Time `json:"date"`
		StartTime   string     `json:"start_time" binding:"required"`
		EndTime     string     `json:"end_time" binding:"required"`
	}

	var input AvailabilityInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, APIResponse{
			Status:  "error",
			Message: "Invalid input",
			Error:   err.Error(),
		})
		return
	}

	// Optional: Check if Provider exists
	var provider models.Provider
	if err := db.DB.First(&provider, input.ProviderID).Error; err != nil {
		c.JSON(http.StatusBadRequest, APIResponse{
			Status:  "error",
			Message: "Provider not found",
		})
		return
	}

	// Validate DayOfWeek/Date consistency
	tempAvailability := models.Availability{
		ProviderID:  input.ProviderID,
		DayOfWeek:   input.DayOfWeek,
		IsRecurring: input.IsRecurring,
		Date:        input.Date,
		StartTime:   input.StartTime,
		EndTime:     input.EndTime,
	}
	if err := scripts.ValidateAvailabilityInput(tempAvailability); err != nil {
		c.JSON(http.StatusBadRequest, APIResponse{
			Status:  "error",
			Message: err.Error(),
		})
		return
	}

	// Update fields
	availability.ProviderID = input.ProviderID
	availability.DayOfWeek = input.DayOfWeek
	availability.IsRecurring = input.IsRecurring
	availability.Date = input.Date
	availability.StartTime = input.StartTime
	availability.EndTime = input.EndTime

	if err := db.DB.Save(&availability).Error; err != nil {
		c.JSON(http.StatusInternalServerError, APIResponse{
			Status:  "error",
			Message: "Failed to update availability",
			Error:   err.Error(),
		})
		return
	}

	// Success response
	c.JSON(http.StatusOK, APIResponse{
		Status:  "success",
		Message: "Availability updated successfully",
		Data:    availability,
	})
}

// DeleteAvailability handles DELETE /availabilities/:id
//
// Delete an existing availability.
// Steps:
// 1️⃣ Parse ID from URL.
// 2️⃣ Check if availability exists.
// 3️⃣ Delete from the database.
// 4️⃣ Return a success response.
func DeleteAvailability(c *gin.Context) {
	// Parse ID
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, APIResponse{
			Status:  "error",
			Message: "Invalid availability ID",
			Error:   err.Error(),
		})
		return
	}

	// Check existence
	var availability models.Availability
	if err := db.DB.First(&availability, id).Error; err != nil {
		c.JSON(http.StatusNotFound, APIResponse{
			Status:  "error",
			Message: "Availability not found",
			Error:   err.Error(),
		})
		return
	}

	// Delete
	if err := db.DB.Delete(&availability).Error; err != nil {
		c.JSON(http.StatusInternalServerError, APIResponse{
			Status:  "error",
			Message: "Failed to delete availability",
			Error:   err.Error(),
		})
		return
	}

	// Success response
	c.JSON(http.StatusOK, APIResponse{
		Status:  "success",
		Message: "Availability deleted successfully",
	})
}

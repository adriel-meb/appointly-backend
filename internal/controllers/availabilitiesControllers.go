package controllers

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/adriel-meb/appointly-backend/internal/db"
	"github.com/adriel-meb/appointly-backend/internal/models"
	"github.com/adriel-meb/appointly-backend/scripts"
	"github.com/gin-gonic/gin"
)

// -----------------------------
// 1️⃣ CREATE AVAILABILITY
// -----------------------------
func CreateAvailability(c *gin.Context) {
	type AvailabilityInput struct {
		ProviderID  uint                  `json:"provider_id" binding:"required"`
		DayOfWeek   *models.DayOfWeekEnum `json:"day_of_week"`
		IsRecurring bool                  `json:"is_recurring"`
		Date        *time.Time            `json:"date"`
		StartTime   string                `json:"start_time" binding:"required"`
		EndTime     string                `json:"end_time" binding:"required"`
	}

	var input AvailabilityInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Invalid input", "error": err.Error()})
		return
	}

	// Normalize DayOfWeek to uppercase
	if input.DayOfWeek != nil {
		day := strings.ToUpper(string(*input.DayOfWeek))
		enumDay := models.DayOfWeekEnum(day)
		input.DayOfWeek = &enumDay
	}

	// Check provider exists
	var provider models.Provider
	if err := db.DB.First(&provider, input.ProviderID).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Provider not found"})
		return
	}

	// Fetch existing slots for overlap validation
	var existing []models.Availability
	query := db.DB.Where("provider_id = ?", input.ProviderID)
	if input.IsRecurring {
		query = query.Where("day_of_week = ? AND is_recurring = true", input.DayOfWeek)
	} else {
		query = query.Where("date = ? AND is_recurring = false", input.Date)
	}
	query.Find(&existing)

	// Prepare availability struct
	newAvail := models.Availability{
		ProviderID:  input.ProviderID,
		DayOfWeek:   input.DayOfWeek,
		IsRecurring: input.IsRecurring,
		Date:        input.Date,
		StartTime:   input.StartTime,
		EndTime:     input.EndTime,
	}

	// Validate input
	if err := scripts.ValidateAvailabilityInput(newAvail, existing); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": err.Error()})
		return
	}

	// Save to DB
	tx := db.DB.Begin()
	if err := tx.Create(&newAvail).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Failed to create availability", "error": err.Error()})
		return
	}
	tx.Commit()

	c.JSON(http.StatusCreated, gin.H{"status": "success", "message": "Availability created successfully", "data": newAvail})
}

// -----------------------------
// 2️⃣ GET ALL AVAILABILITIES
// -----------------------------
func GetAllAvailability(c *gin.Context) {
	var availabilities []models.Availability
	dbQuery := db.DB.Preload("Provider")

	// Optional filters
	if providerID := c.Query("provider_id"); providerID != "" {
		dbQuery = dbQuery.Where("provider_id = ?", providerID)
	}

	if dateStr := c.Query("date"); dateStr != "" {
		date, err := time.Parse("2006-01-02", dateStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Invalid date format. Use YYYY-MM-DD."})
			return
		}
		dbQuery = dbQuery.Where("date = ?", date)
	}

	if startStr := c.Query("start_date"); startStr != "" {
		start, err := time.Parse("2006-01-02", startStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Invalid start_date format. Use YYYY-MM-DD."})
			return
		}
		dbQuery = dbQuery.Where("date >= ?", start)
	}

	if endStr := c.Query("end_date"); endStr != "" {
		end, err := time.Parse("2006-01-02", endStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Invalid end_date format. Use YYYY-MM-DD."})
			return
		}
		dbQuery = dbQuery.Where("date <= ?", end)
	}

	if err := dbQuery.Find(&availabilities).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Failed to fetch availabilities", "error": err.Error()})
		return
	}

	if len(availabilities) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"status": "error", "message": "the provider has no availabilities recorded"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "success", "message": "Availabilities fetched successfully", "data": availabilities})
}

// -----------------------------
// 3️⃣ GET AVAILABILITY BY ID
// -----------------------------
func GetAvailabilityByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Invalid availability ID", "error": err.Error()})
		return
	}

	var availability models.Availability
	if err := db.DB.Preload("Provider").First(&availability, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"status": "error", "message": "Availability not found", "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "message": "Availability fetched successfully", "data": availability})
}

// -----------------------------
// 4️⃣ UPDATE AVAILABILITY
// -----------------------------
func UpdateAvailability(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Invalid availability ID", "error": err.Error()})
		return
	}

	var availability models.Availability
	if err := db.DB.First(&availability, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"status": "error", "message": "Availability not found", "error": err.Error()})
		return
	}

	type AvailabilityInput struct {
		ProviderID  uint                  `json:"provider_id" binding:"required"`
		DayOfWeek   *models.DayOfWeekEnum `json:"day_of_week"`
		IsRecurring bool                  `json:"is_recurring"`
		Date        *time.Time            `json:"date"`
		StartTime   string                `json:"start_time" binding:"required"`
		EndTime     string                `json:"end_time" binding:"required"`
	}

	var input AvailabilityInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Invalid input", "error": err.Error()})
		return
	}

	if input.DayOfWeek != nil {
		day := strings.ToUpper(string(*input.DayOfWeek))
		enumDay := models.DayOfWeekEnum(day)
		input.DayOfWeek = &enumDay
	}

	// Fetch existing slots excluding the one being updated
	var existing []models.Availability
	query := db.DB.Where("provider_id = ? AND id != ?", input.ProviderID, availability.ID)
	if input.IsRecurring && input.DayOfWeek != nil {
		query = query.Where("day_of_week = ? AND is_recurring = true", input.DayOfWeek)
	} else if !input.IsRecurring && input.Date != nil {
		query = query.Where("date = ? AND is_recurring = false", input.Date)
	}
	query.Find(&existing)

	// Validate
	updated := models.Availability{
		ProviderID:  input.ProviderID,
		DayOfWeek:   input.DayOfWeek,
		IsRecurring: input.IsRecurring,
		Date:        input.Date,
		StartTime:   input.StartTime,
		EndTime:     input.EndTime,
	}

	if err := scripts.ValidateAvailabilityInput(updated, existing); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": err.Error()})
		return
	}

	// Apply updates
	availability.ProviderID = input.ProviderID
	availability.DayOfWeek = input.DayOfWeek
	availability.IsRecurring = input.IsRecurring
	availability.Date = input.Date
	availability.StartTime = input.StartTime
	availability.EndTime = input.EndTime

	if err := db.DB.Save(&availability).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Failed to update availability", "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "message": "Availability updated successfully", "data": availability})
}

// -----------------------------
// 5️⃣ DELETE AVAILABILITY
// -----------------------------
func DeleteAvailability(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Invalid availability ID", "error": err.Error()})
		return
	}

	var availability models.Availability
	if err := db.DB.First(&availability, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"status": "error", "message": "Availability not found", "error": err.Error()})
		return
	}

	if err := db.DB.Delete(&availability).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Failed to delete availability", "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "message": "Availability deleted successfully"})
}

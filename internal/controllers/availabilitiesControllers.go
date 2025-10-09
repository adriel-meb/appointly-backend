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

type AvailabilitySlotResponse struct {
	ID        uint   `json:"id"`
	StartTime string `json:"start_time"`
	EndTime   string `json:"end_time"`
	IsBooked  bool   `json:"is_booked"`
}

type AvailabilityResponse struct {
	ID          uint                       `json:"id"`
	ProviderID  uint                       `json:"provider_id"`
	DayOfWeek   *models.DayOfWeekEnum      `json:"day_of_week"`
	IsRecurring bool                       `json:"is_recurring"`
	Date        *time.Time                 `json:"date"`
	StartTime   string                     `json:"start_time"`
	EndTime     string                     `json:"end_time"`
	Slots       []AvailabilitySlotResponse `json:"slots"`
}

// -----------------------------
// 1️⃣ CREATE AVAILABILITY
// -----------------------------
func CreateAvailability(c *gin.Context) {
	type AvailabilityInput struct {
		ProviderID  uint                  `json:"provider_id" binding:"required"`
		DayOfWeek   *models.DayOfWeekEnum `json:"day_of_week"`
		IsRecurring bool                  `json:"is_recurring"`
		Date        string                `json:"date"` // now as string for flexible parsing
		StartTime   string                `json:"start_time" binding:"required"`
		EndTime     string                `json:"end_time" binding:"required"`
		SlotMinutes int                   `json:"slot_minutes"` // optional, default 30
	}

	var input AvailabilityInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Invalid input",
			"error":   err.Error(),
		})
		return
	}

	// Default slot length = 30 min
	if input.SlotMinutes <= 0 {
		input.SlotMinutes = 30
	}

	// Normalize DayOfWeek
	if input.DayOfWeek != nil {
		day := strings.ToUpper(string(*input.DayOfWeek))
		enumDay := models.DayOfWeekEnum(day)
		input.DayOfWeek = &enumDay
	}

	// ✅ Parse date if provided
	var parsedDate *time.Time
	if input.Date != "" {
		dateVal, err := scripts.ParseDateFlexible(input.Date)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  "error",
				"message": "Invalid date format (use YYYY-MM-DD or DD-MM-YYYY)",
				"error":   err.Error(),
			})
			return
		}
		parsedDate = dateVal
	}

	// ✅ Check provider exists
	var provider models.Provider
	if err := db.DB.First(&provider, input.ProviderID).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Provider not found",
		})
		return
	}

	// ✅ Check for existing overlaps
	var existing []models.Availability
	query := db.DB.Where("provider_id = ?", input.ProviderID)
	if input.IsRecurring {
		query = query.Where("day_of_week = ? AND is_recurring = true", input.DayOfWeek)
	} else {
		query = query.Where("date = ? AND is_recurring = false", parsedDate)
	}
	query.Find(&existing)

	newAvail := models.Availability{
		ProviderID:  input.ProviderID,
		DayOfWeek:   input.DayOfWeek,
		IsRecurring: input.IsRecurring,
		Date:        parsedDate,
		StartTime:   input.StartTime,
		EndTime:     input.EndTime,
	}

	if err := scripts.ValidateAvailabilityInput(newAvail, existing); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	// ✅ Begin transaction
	tx := db.DB.Begin()

	// Create availability record
	if err := tx.Create(&newAvail).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Failed to create availability",
			"error":   err.Error(),
		})
		return
	}

	// ✅ Generate slots
	slotTimes, err := scripts.GenerateTimeSlots(input.StartTime, input.EndTime, input.SlotMinutes, true, true)
	if err != nil {
		tx.Rollback()
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Failed to generate slots",
			"error":   err.Error(),
		})
		return
	}

	// Build slot models
	var slots []models.AvailabilitySlot
	for i := 0; i < len(slotTimes)-1; i++ {
		slots = append(slots, models.AvailabilitySlot{
			AvailabilityID: newAvail.ID,
			StartTime:      slotTimes[i],
			EndTime:        slotTimes[i+1],
			IsBooked:       false,
		})
	}

	// ✅ Save slots
	if len(slots) > 0 {
		if err := tx.Create(&slots).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  "error",
				"message": "Failed to save slots",
				"error":   err.Error(),
			})
			return
		}
	}

	tx.Commit()

	c.JSON(http.StatusCreated, gin.H{
		"status":  "success",
		"message": "Availability created successfully with slots",
		"data": gin.H{
			"availability": newAvail,
			"slots":        slots,
		},
	})
}

// -----------------------------
// 2️⃣ GET ALL AVAILABILITIES
// -----------------------------
func GetAllAvailability(c *gin.Context) {
	var availabilities []models.Availability
	dbQuery := db.DB.Preload("Provider").Preload("Slots")

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
		c.JSON(http.StatusNotFound, gin.H{"status": "error", "message": "No availabilities recorded"})
		return
	}

	// Map to clean response
	var response []AvailabilityResponse
	for _, a := range availabilities {
		var slots []AvailabilitySlotResponse
		for _, s := range a.Slots {
			slots = append(slots, AvailabilitySlotResponse{
				ID:        s.ID,
				StartTime: s.StartTime,
				EndTime:   s.EndTime,
				IsBooked:  s.IsBooked,
			})
		}
		response = append(response, AvailabilityResponse{
			ID:          a.ID,
			ProviderID:  a.ProviderID,
			DayOfWeek:   a.DayOfWeek,
			IsRecurring: a.IsRecurring,
			Date:        a.Date,
			StartTime:   a.StartTime,
			EndTime:     a.EndTime,
			Slots:       slots,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Availabilities fetched successfully",
		"data":    response,
	})
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
	if err := db.DB.Preload("Provider").Preload("Slots").First(&availability, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"status": "error", "message": "Availability not found", "error": err.Error()})
		return
	}

	var slots []AvailabilitySlotResponse
	for _, s := range availability.Slots {
		slots = append(slots, AvailabilitySlotResponse{
			ID:        s.ID,
			StartTime: s.StartTime,
			EndTime:   s.EndTime,
			IsBooked:  s.IsBooked,
		})
	}

	response := AvailabilityResponse{
		ID:          availability.ID,
		ProviderID:  availability.ProviderID,
		DayOfWeek:   availability.DayOfWeek,
		IsRecurring: availability.IsRecurring,
		Date:        availability.Date,
		StartTime:   availability.StartTime,
		EndTime:     availability.EndTime,
		Slots:       slots,
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Availability fetched successfully",
		"data":    response,
	})
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
	if err := db.DB.Preload("Slots").First(&availability, id).Error; err != nil {
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
		SlotMinutes int                   `json:"slot_minutes"` // optional, default 30
	}

	var input AvailabilityInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Invalid input", "error": err.Error()})
		return
	}

	if input.SlotMinutes <= 0 {
		input.SlotMinutes = 30
	}

	// Normalize DayOfWeek
	if input.DayOfWeek != nil {
		day := strings.ToUpper(string(*input.DayOfWeek))
		enumDay := models.DayOfWeekEnum(day)
		input.DayOfWeek = &enumDay
	}

	// Check existing availabilities excluding current one
	var existing []models.Availability
	query := db.DB.Where("provider_id = ? AND id != ?", input.ProviderID, availability.ID)
	if input.IsRecurring && input.DayOfWeek != nil {
		query = query.Where("day_of_week = ? AND is_recurring = true", input.DayOfWeek)
	} else if !input.IsRecurring && input.Date != nil {
		query = query.Where("date = ? AND is_recurring = false", input.Date)
	}
	query.Find(&existing)

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

	// Begin transaction
	tx := db.DB.Begin()

	// Update availability fields
	availability.ProviderID = input.ProviderID
	availability.DayOfWeek = input.DayOfWeek
	availability.IsRecurring = input.IsRecurring
	availability.Date = input.Date
	availability.StartTime = input.StartTime
	availability.EndTime = input.EndTime

	if err := tx.Save(&availability).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Failed to update availability", "error": err.Error()})
		return
	}

	// Delete old slots
	if len(availability.Slots) > 0 {
		if err := tx.Where("availability_id = ?", availability.ID).Delete(&models.AvailabilitySlot{}).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Failed to delete old slots", "error": err.Error()})
			return
		}
	}

	// Generate new slots
	slotTimes, err := scripts.GenerateTimeSlots(input.StartTime, input.EndTime, input.SlotMinutes, true, false)
	if err != nil {
		tx.Rollback()
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Failed to generate slots", "error": err.Error()})
		return
	}

	var newSlots []models.AvailabilitySlot
	for i := 0; i < len(slotTimes)-1; i++ {
		newSlots = append(newSlots, models.AvailabilitySlot{
			AvailabilityID: availability.ID,
			StartTime:      slotTimes[i],
			EndTime:        slotTimes[i+1],
			IsBooked:       false,
		})
	}

	if len(newSlots) > 0 {
		if err := tx.Create(&newSlots).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Failed to save new slots", "error": err.Error()})
			return
		}
	}

	tx.Commit()

	c.JSON(http.StatusOK, gin.H{"status": "success", "message": "Availability updated with slots successfully", "data": availability})
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

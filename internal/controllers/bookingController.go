package controllers

import (
	"github.com/adriel-meb/appointly-backend/internal/db"
	"github.com/adriel-meb/appointly-backend/internal/models"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"time"
)

// CreateBookingInput represents request body
type CreateBookingInput struct {
	PatientID  uint   `json:"patient_id" binding:"required"`
	ProviderID uint   `json:"provider_id" binding:"required"`
	ServiceID  uint   `json:"service_id" binding:"required"`
	StartTime  string `json:"start_time" binding:"required"` // ISO8601 e.g. "2025-09-02T15:00:00Z"
	Notes      string `json:"notes" `
}

// CreateBooking handles POST /bookings
func CreateBooking(c *gin.Context) {
	var input CreateBookingInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Invalid input", "error": err.Error()})
		return
	}

	// 1. Parse start_time
	startTime, err := time.Parse(time.RFC3339, input.StartTime)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Invalid time format. Use RFC3339 (e.g. 2025-09-02T15:00:00Z)"})
		return
	}

	// 2. Fetch service to get duration
	var service models.Service
	if err := db.DB.First(&service, input.ServiceID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"status": "error", "message": "Service not found"})
		return
	}

	// 3. Compute end time
	endTime := startTime.Add(time.Duration(service.DurationMinutes) * time.Minute)

	// 4. Check provider availability
	var availability models.Availability
	err = db.DB.Where("provider_id = ?", input.ProviderID).
		Where("? BETWEEN start_time AND end_time", startTime.Format("15:04")).
		Where("? BETWEEN start_time AND end_time", endTime.Format("15:04")).
		First(&availability).Error

	if err != nil {
		c.JSON(http.StatusConflict, gin.H{"status": "error", "message": "Provider not available at this time"})
		return
	}

	// 5. Check for overlapping bookings
	var overlap int64
	db.DB.Model(&models.Booking{}).
		Where("provider_id = ?", input.ProviderID).
		Where("status IN ?", []string{"pending", "confirmed"}).
		Where("start_time < ? AND end_time > ?", endTime, startTime).
		Count(&overlap)

	if overlap > 0 {
		c.JSON(http.StatusConflict, gin.H{"status": "error", "message": "This slot is already booked"})
		return
	}

	// 6. Create booking
	booking := models.Booking{
		PatientID:  input.PatientID,
		ProviderID: input.ProviderID,
		ServiceID:  input.ServiceID,
		StartTime:  startTime,
		EndTime:    endTime,
		Status:     models.Pending, // default
		Notes:      input.Notes,
		Amount:     service.Price,
	}

	if err := db.DB.Create(&booking).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Failed to create booking", "error": err.Error()})
		return
	}

	// 7. Success response
	c.JSON(http.StatusCreated, gin.H{
		"status":  "success",
		"message": "Booking created successfully",
		"data":    booking,
	})
}

func GetAllBooking(c *gin.Context) {

	var bookings []models.Booking

	if err := db.DB.Find(&bookings).Error; err != nil {
		c.JSON(http.StatusInternalServerError, APIResponse{
			Status:  "error",
			Message: "Failed to fetch bookings",
			Error:   err.Error(),
		})
		return
	}

	// Success response
	c.JSON(http.StatusOK, APIResponse{
		Status:  "success",
		Message: "bookings fetched successfully",
		Data:    bookings,
	})
}

func ConfirmBooking(c *gin.Context) {
	// booking input form
	type ConfirmInput struct {
		ID uint `json:"id" binding:"required"`
	}

	var inputID ConfirmInput
	if err := c.ShouldBindJSON(&inputID); err != nil {
		c.JSON(http.StatusBadRequest, APIResponse{
			Status:  "error",
			Message: "Invalid input",
			Error:   err.Error(),
		})
		return
	}

	// check if booking exists
	var booking models.Booking
	if err := db.DB.First(&booking, inputID.ID).Error; err != nil {
		c.JSON(http.StatusNotFound, APIResponse{
			Status:  "error",
			Message: "Booking not found",
		})
		return
	}

	// ensure booking is still pending
	if booking.Status != models.Pending {
		c.JSON(http.StatusBadRequest, APIResponse{
			Status:  "error",
			Message: "Only pending bookings can be confirmed",
		})
		return
	}

	// update the status
	booking.Status = models.Confirmed
	if err := db.DB.Save(&booking).Error; err != nil {
		c.JSON(http.StatusInternalServerError, APIResponse{
			Status:  "error",
			Message: "Failed to update booking",
			Error:   err.Error(),
		})
		return
	}

	// return success
	c.JSON(http.StatusOK, APIResponse{
		Status:  "success",
		Message: "Booking confirmed successfully",
		Data:    booking,
	})
}

func CancelBooking(c *gin.Context) {
	type CancelInput struct {
		ID uint `json:"id" binding:"required"`
	}

	var input CancelInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, APIResponse{
			Status:  "error",
			Message: "Invalid input",
			Error:   err.Error(),
		})
		return
	}

	var booking models.Booking
	if err := db.DB.First(&booking, input.ID).Error; err != nil {
		c.JSON(http.StatusNotFound, APIResponse{
			Status:  "error",
			Message: "Booking not found",
		})
		return
	}

	if booking.Status == models.Completed {
		c.JSON(http.StatusBadRequest, APIResponse{
			Status:  "error",
			Message: "Completed bookings cannot be cancelled",
		})
		return
	}
	if booking.Status == models.Cancelled {
		c.JSON(http.StatusBadRequest, APIResponse{
			Status:  "error",
			Message: "Booking is already cancelled",
		})
		return
	}

	booking.Status = models.Cancelled
	if err := db.DB.Save(&booking).Error; err != nil {
		c.JSON(http.StatusInternalServerError, APIResponse{
			Status:  "error",
			Message: "Failed to cancel booking",
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, APIResponse{
		Status:  "success",
		Message: "Booking cancelled successfully",
		Data:    booking,
	})
}

func CompleteBooking(c *gin.Context) {
	type CompleteInput struct {
		ID uint `json:"id" binding:"required"`
	}

	var input CompleteInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, APIResponse{
			Status:  "error",
			Message: "Invalid input",
			Error:   err.Error(),
		})
		return
	}

	var booking models.Booking
	if err := db.DB.First(&booking, input.ID).Error; err != nil {
		c.JSON(http.StatusNotFound, APIResponse{
			Status:  "error",
			Message: "Booking not found",
		})
		return
	}

	if booking.Status != models.Confirmed {
		c.JSON(http.StatusBadRequest, APIResponse{
			Status:  "error",
			Message: "Only confirmed bookings can be marked as completed",
		})
		return
	}

	booking.Status = models.Completed
	if err := db.DB.Save(&booking).Error; err != nil {
		c.JSON(http.StatusInternalServerError, APIResponse{
			Status:  "error",
			Message: "Failed to complete booking",
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, APIResponse{
		Status:  "success",
		Message: "Booking marked as completed successfully",
		Data:    booking,
	})
}

func GetBookingByID(c *gin.Context) {
	// Extract booking ID from URL params
	idParam := c.Param("id")
	if idParam == "" {
		c.JSON(http.StatusBadRequest, APIResponse{
			Status:  "error",
			Message: "Booking ID is required",
		})
		return
	}

	// Convert to uint
	id, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, APIResponse{
			Status:  "error",
			Message: "Invalid booking ID format",
			Error:   err.Error(),
		})
		return
	}

	// Find booking
	var booking models.Booking
	if err := db.DB.First(&booking, id).Error; err != nil {
		c.JSON(http.StatusNotFound, APIResponse{
			Status:  "error",
			Message: "Booking not found",
		})
		return
	}

	// Success response
	c.JSON(http.StatusOK, APIResponse{
		Status:  "success",
		Message: "Booking fetched successfully",
		Data:    booking,
	})
}

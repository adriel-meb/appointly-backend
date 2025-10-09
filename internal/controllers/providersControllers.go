package controllers

import (
	"github.com/adriel-meb/appointly-backend/internal/db"
	"github.com/adriel-meb/appointly-backend/internal/models"
	"github.com/adriel-meb/appointly-backend/scripts"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

type InsuranceResponse2 struct {
	Name string `json:"name"`
}

type AvailabilitiesResponse struct {
	IsRecurring bool                  `json:"is_recurring"`
	DayOfWeek   *models.DayOfWeekEnum `json:"day_of_week"`
	Date        *string               `json:"date"`
	StartTime   string                `json:"start_time"`
	EndTime     string                `json:"end_time"`
}
type ProviderResponse struct {
	ID             uint                     `json:"id"`
	Bio            string                   `json:"bio"`
	Rating         float32                  `json:"rating"`
	Price          string                   `json:"price"`
	Address        string                   `json:"address"`
	Lat            float64                  `json:"lat"`
	Lng            float64                  `json:"lng"`
	ImageURL       string                   `json:"image_url"`
	UserName       string                   `json:"user_name"`
	UserEmail      string                   `json:"user_email"`
	UserPhone      *string                  `json:"user_phone"`
	Specialization string                   `json:"specialization"`
	City           string                   `json:"city"`
	Insurances     []InsuranceResponse2     `json:"insurances"`
	Availabilities []AvailabilitiesResponse `json:"availabilities"`
}

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
// Supports optional filters: ?search=dr+name&location=Libreville&insurance=CNAMGS
func GetAllProviders(c *gin.Context) {
	var providers []models.Provider

	query := c.Query("search")
	city := c.Query("location")
	insurance := c.Query("insurance")

	tx := db.DB.Preload("User").
		Preload("Specialization").
		Preload("City").
		Preload("Insurances").
		Preload("Availabilities")

	// ------------------ Filters ------------------

	// Search in bio, user name, specialization
	if query != "" {
		query = strings.ToLower(query)
		tx = tx.Joins("JOIN users u ON u.id = providers.user_id").
			Joins("JOIN specializations s ON s.id = providers.specialization_id").
			Where("LOWER(providers.bio) LIKE ? OR LOWER(u.name) LIKE ? OR LOWER(s.name) LIKE ?",
				"%"+query+"%", "%"+query+"%", "%"+query+"%")
	}

	// Filter by city name
	if city != "" {
		tx = tx.Joins("JOIN cities c ON c.id = providers.city_id").
			Where("LOWER(c.name) LIKE ?", "%"+city+"%")
	}

	// Filter by insurance name
	if insurance != "" {
		tx = tx.Joins("JOIN provider_insurances pi ON pi.provider_id = providers.id").
			Joins("JOIN insurances i ON i.id = pi.insurance_id").
			Where("LOWER(i.name) LIKE ?", "%"+insurance+"%").
			Distinct("providers.id")
	}

	// ------------------ Execute Query ------------------

	if err := tx.Find(&providers).Error; err != nil {
		c.JSON(http.StatusInternalServerError, APIResponse{
			Status:  "error",
			Message: "Failed to fetch providers",
			Error:   err.Error(),
		})
		return
	}

	// ------------------ Map to Custom Response ------------------

	var response []ProviderResponse

	for _, p := range providers {
		// Format insurances
		var insurances []InsuranceResponse2
		for _, ins := range p.Insurances {
			insurances = append(insurances, InsuranceResponse2{
				Name: ins.Name,
			})
		}

		// Format availabilities
		var availabilities []AvailabilitiesResponse
		for _, availability := range p.Availabilities {
			availabilities = append(availabilities, AvailabilitiesResponse{
				Date:        scripts.FormatDate(availability.Date),
				IsRecurring: availability.IsRecurring,
				StartTime:   availability.StartTime,
				EndTime:     availability.EndTime,
				DayOfWeek:   availability.DayOfWeek,
			})
		}

		// Append provider
		response = append(response, ProviderResponse{
			ID:             p.ID,
			UserName:       p.User.Name,
			Bio:            p.Bio,
			Specialization: p.Specialization.Name,
			City:           p.City.Name,
			Insurances:     insurances,
			Rating:         p.Rating,
			Price:          p.Price,
			Address:        p.Address,
			Lat:            p.Lat,
			Lng:            p.Lng,
			ImageURL:       p.ImageURL,
			UserEmail:      p.User.Email,
			UserPhone:      p.User.PhoneNumber,
			Availabilities: availabilities,
		})
	}

	// ------------------ Response ------------------

	c.JSON(http.StatusOK, APIResponse{
		Status:  "success",
		Length:  len(response),
		Message: "Providers fetched successfully",
		Data:    response,
	})
}

// GetProviderByID placeholder
// GetProviderByID handles GET /providers/:id
func GetProviderByID(c *gin.Context) {
	id := c.Param("id")
	var provider models.Provider

	if err := db.DB.Preload("User").
		Preload("Specialization").
		Preload("City").
		Preload("Insurances").
		Preload("Availabilities").
		First(&provider, id).Error; err != nil {
		c.JSON(http.StatusNotFound, APIResponse{
			Status:  "error",
			Message: "Provider not found",
			Error:   err.Error(),
		})
		return
	}

	var availabilities []AvailabilitiesResponse
	for _, availability := range provider.Availabilities {
		availabilities = append(availabilities, AvailabilitiesResponse{
			Date:        scripts.FormatDate(availability.Date),
			IsRecurring: availability.IsRecurring,
			StartTime:   availability.StartTime,
			EndTime:     availability.EndTime,
			DayOfWeek:   availability.DayOfWeek,
		})

	}

	var insurances []InsuranceResponse2
	for _, ins := range provider.Insurances {
		insurances = append(insurances, InsuranceResponse2{
			Name: ins.Name,
		})
	}

	//customs models
	response := ProviderResponse{
		ID:             provider.ID,
		UserName:       provider.User.Name,
		Bio:            provider.Bio,
		Specialization: provider.Specialization.Name,
		City:           provider.City.Name,
		Insurances:     insurances,
		Rating:         provider.Rating,
		Price:          provider.Price,
		Address:        provider.Address,
		Lat:            provider.Lat,
		Lng:            provider.Lng,
		ImageURL:       provider.ImageURL,
		UserEmail:      provider.User.Email,
		UserPhone:      provider.User.PhoneNumber,
		Availabilities: availabilities,
	}

	//instead of returning in data provider i return response
	c.JSON(http.StatusOK, APIResponse{
		Status:  "success",
		Message: "Provider fetched successfully",
		Data:    response,
	})
}

// UpdateProvider handles PUT /providers/:id
func UpdateProvider(c *gin.Context) {
	id := c.Param("id")
	var provider models.Provider

	if err := db.DB.First(&provider, id).Error; err != nil {
		c.JSON(http.StatusNotFound, APIResponse{
			Status:  "error",
			Message: "Provider not found",
			Error:   err.Error(),
		})
		return
	}

	type UpdateProviderInput struct {
		SpecializationID *uint    `json:"specialization_id"`
		Bio              *string  `json:"bio"`
		CityID           *uint    `json:"city_id"`
		ImageURL         *string  `json:"image_url"`
		Price            *string  `json:"price"`
		Address          *string  `json:"address"`
		Lat              *float64 `json:"lat"`
		Lng              *float64 `json:"lng"`
	}

	var input UpdateProviderInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, APIResponse{
			Status:  "error",
			Message: "Invalid input",
			Error:   err.Error(),
		})
		return
	}

	// Apply updates if values are not nil
	if input.SpecializationID != nil {
		provider.SpecializationID = *input.SpecializationID
	}
	if input.Bio != nil {
		provider.Bio = *input.Bio
	}
	if input.CityID != nil {
		provider.CityID = *input.CityID
	}
	if input.ImageURL != nil {
		provider.ImageURL = *input.ImageURL
	}
	if input.Price != nil {
		provider.Price = *input.Price
	}
	if input.Address != nil {
		provider.Address = *input.Address
	}
	if input.Lat != nil {
		provider.Lat = *input.Lat
	}
	if input.Lng != nil {
		provider.Lng = *input.Lng
	}

	if err := db.DB.Save(&provider).Error; err != nil {
		c.JSON(http.StatusInternalServerError, APIResponse{
			Status:  "error",
			Message: "Failed to update provider",
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, APIResponse{
		Status:  "success",
		Message: "Provider updated successfully",
		Data:    provider,
	})
}

// DeleteProvider handles DELETE /providers/:id
func DeleteProvider(c *gin.Context) {
	id := c.Param("id")
	var provider models.Provider

	if err := db.DB.First(&provider, id).Error; err != nil {
		c.JSON(http.StatusNotFound, APIResponse{
			Status:  "error",
			Message: "Provider not found",
			Error:   err.Error(),
		})
		return
	}

	if err := db.DB.Delete(&provider).Error; err != nil {
		c.JSON(http.StatusInternalServerError, APIResponse{
			Status:  "error",
			Message: "Failed to delete provider",
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, APIResponse{
		Status:  "success",
		Message: "Provider deleted successfully",
	})
}

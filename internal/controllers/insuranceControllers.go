package controllers

import (
	"net/http"
	"strconv"

	"github.com/adriel-meb/appointly-backend/internal/db"
	"github.com/adriel-meb/appointly-backend/internal/models"
	"github.com/gin-gonic/gin"
)

// Type pour ne renvoyer que le nom des villes
type CityNameResponse struct {
	Name string `json:"name"`
}

type InsuranceResponse struct {
	ID          uint               `json:"id"`
	Name        string             `json:"name"`
	Description string             `json:"description"`
	Coverage    string             `json:"coverage"`
	Phone       string             `json:"phone"`
	Email       string             `json:"email"`
	Website     string             `json:"website"`
	LogoURL     string             `json:"logo_url"`
	Cities      []CityNameResponse `json:"cities"`
}

// ------------------ CREATE INSURANCE ------------------

func CreateInsurance(c *gin.Context) {
	var input struct {
		Name        string `json:"name"`
		Description string `json:"description"`
		Coverage    string `json:"coverage"`
		Phone       string `json:"phone"`
		Email       string `json:"email"`
		Website     string `json:"website"`
		LogoURL     string `json:"logo_url"`
		CityIDs     []uint `json:"city_ids"` // IDs des villes associées
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, APIResponse{
			Status:  "error",
			Message: "Invalid input",
			Error:   err.Error(),
		})
		return
	}

	insurance := models.Insurance{
		Name:        input.Name,
		Description: input.Description,
		Coverage:    input.Coverage,
		Phone:       input.Phone,
		Email:       input.Email,
		Website:     input.Website,
		LogoURL:     input.LogoURL,
	}

	// Lier les villes si fournies
	if len(input.CityIDs) > 0 {
		var cities []models.City
		if err := db.DB.Find(&cities, input.CityIDs).Error; err != nil {
			c.JSON(http.StatusInternalServerError, APIResponse{
				Status:  "error",
				Message: "Failed to fetch cities",
				Error:   err.Error(),
			})
			return
		}
		insurance.Cities = cities
	}

	if err := db.DB.Create(&insurance).Error; err != nil {
		c.JSON(http.StatusInternalServerError, APIResponse{
			Status:  "error",
			Message: "Failed to create insurance",
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, APIResponse{
		Status:  "success",
		Message: "Insurance created successfully",
		Data:    insurance,
	})
}

// ------------------ GET ALL INSURANCES ------------------

func GetAllInsurances(c *gin.Context) {
	var insurances []models.Insurance

	name := c.Query("name")
	coverage := c.Query("coverage")
	email := c.Query("email")

	query := db.DB.Preload("Cities").Model(&models.Insurance{})
	if name != "" {
		query = query.Where("LOWER(name) LIKE ?", "%"+name+"%")
	}
	if coverage != "" {
		query = query.Where("LOWER(coverage) LIKE ?", "%"+coverage+"%")
	}
	if email != "" {
		query = query.Where("LOWER(email) LIKE ?", "%"+email+"%")
	}

	if err := query.Find(&insurances).Error; err != nil {
		c.JSON(http.StatusInternalServerError, APIResponse{
			Status:  "error",
			Message: "Failed to fetch insurances",
			Error:   err.Error(),
		})
		return
	}

	// Transformer pour ne garder que le nom des villes
	var response []InsuranceResponse
	for _, ins := range insurances {
		var cities []CityNameResponse
		for _, city := range ins.Cities {
			cities = append(cities, CityNameResponse{Name: city.Name})
		}

		response = append(response, InsuranceResponse{
			ID:          ins.ID,
			Name:        ins.Name,
			Description: ins.Description,
			Coverage:    ins.Coverage,
			Phone:       ins.Phone,
			Email:       ins.Email,
			Website:     ins.Website,
			LogoURL:     ins.LogoURL,
			Cities:      cities,
		})
	}

	c.JSON(http.StatusOK, APIResponse{
		Status:  "success",
		Message: "Insurances fetched successfully",
		Data:    response,
	})
}

// ------------------ GET INSURANCE BY ID ------------------

func GetInsuranceByID(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, APIResponse{
			Status:  "error",
			Message: "Invalid insurance ID format",
			Error:   err.Error(),
		})
		return
	}

	var insurance models.Insurance
	if err := db.DB.Preload("Cities").First(&insurance, id).Error; err != nil {
		c.JSON(http.StatusNotFound, APIResponse{
			Status:  "error",
			Message: "Insurance not found",
		})
		return
	}

	var cities []CityNameResponse
	for _, city := range insurance.Cities {
		cities = append(cities, CityNameResponse{Name: city.Name})
	}

	response := InsuranceResponse{
		ID:          insurance.ID,
		Name:        insurance.Name,
		Description: insurance.Description,
		Coverage:    insurance.Coverage,
		Phone:       insurance.Phone,
		Email:       insurance.Email,
		Website:     insurance.Website,
		LogoURL:     insurance.LogoURL,
		Cities:      cities,
	}

	c.JSON(http.StatusOK, APIResponse{
		Status:  "success",
		Message: "Insurance fetched successfully",
		Data:    response,
	})
}

// ------------------ UPDATE INSURANCE ------------------

func UpdateInsurance(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, APIResponse{
			Status:  "error",
			Message: "Invalid insurance ID format",
			Error:   err.Error(),
		})
		return
	}

	var insurance models.Insurance
	if err := db.DB.Preload("Cities").First(&insurance, id).Error; err != nil {
		c.JSON(http.StatusNotFound, APIResponse{
			Status:  "error",
			Message: "Insurance not found",
		})
		return
	}

	var input struct {
		Name        string `json:"name"`
		Description string `json:"description"`
		Coverage    string `json:"coverage"`
		Phone       string `json:"phone"`
		Email       string `json:"email"`
		Website     string `json:"website"`
		LogoURL     string `json:"logo_url"`
		CityIDs     []uint `json:"city_ids"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, APIResponse{
			Status:  "error",
			Message: "Invalid input",
			Error:   err.Error(),
		})
		return
	}

	insurance.Name = input.Name
	insurance.Description = input.Description
	insurance.Coverage = input.Coverage
	insurance.Phone = input.Phone
	insurance.Email = input.Email
	insurance.Website = input.Website
	insurance.LogoURL = input.LogoURL

	// Mettre à jour les villes
	if input.CityIDs != nil {
		var cities []models.City
		if err := db.DB.Find(&cities, input.CityIDs).Error; err != nil {
			c.JSON(http.StatusInternalServerError, APIResponse{
				Status:  "error",
				Message: "Failed to fetch cities",
				Error:   err.Error(),
			})
			return
		}
		db.DB.Model(&insurance).Association("Cities").Replace(cities)
	}

	if err := db.DB.Save(&insurance).Error; err != nil {
		c.JSON(http.StatusInternalServerError, APIResponse{
			Status:  "error",
			Message: "Failed to update insurance",
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, APIResponse{
		Status:  "success",
		Message: "Insurance updated successfully",
		Data:    insurance,
	})
}

// ------------------ DELETE INSURANCE ------------------

func DeleteInsurance(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, APIResponse{
			Status:  "error",
			Message: "Invalid insurance ID format",
			Error:   err.Error(),
		})
		return
	}

	var insurance models.Insurance
	if err := db.DB.First(&insurance, id).Error; err != nil {
		c.JSON(http.StatusNotFound, APIResponse{
			Status:  "error",
			Message: "Insurance not found",
		})
		return
	}

	// Supprimer les associations villes avant suppression
	if err := db.DB.Model(&insurance).Association("Cities").Clear(); err != nil {
		c.JSON(http.StatusInternalServerError, APIResponse{
			Status:  "error",
			Message: "Failed to clear city associations",
			Error:   err.Error(),
		})
		return
	}

	if err := db.DB.Delete(&insurance).Error; err != nil {
		c.JSON(http.StatusInternalServerError, APIResponse{
			Status:  "error",
			Message: "Failed to delete insurance",
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, APIResponse{
		Status:  "success",
		Message: "Insurance deleted successfully",
	})
}

package controllers

import (
	"net/http"
	"strconv"

	"github.com/adriel-meb/appointly-backend/internal/db"
	"github.com/adriel-meb/appointly-backend/internal/models"
	"github.com/gin-gonic/gin"
)

// ------------------ CREATE CITY ------------------

func CreateCity(c *gin.Context) {
	var input models.City

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, APIResponse{
			Status:  "error",
			Message: "Invalid input",
			Error:   err.Error(),
		})
		return
	}

	if err := db.DB.Create(&input).Error; err != nil {
		c.JSON(http.StatusInternalServerError, APIResponse{
			Status:  "error",
			Message: "Failed to create city",
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, APIResponse{
		Status:  "success",
		Message: "City created successfully",
		Data:    input,
	})
}

// ------------------ GET ALL CITIES ------------------

func GetAllCities(c *gin.Context) {
	var cities []models.City

	if err := db.DB.Preload("Insurances").Find(&cities).Error; err != nil {
		c.JSON(http.StatusInternalServerError, APIResponse{
			Status:  "error",
			Message: "Failed to fetch cities",
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, APIResponse{
		Status:  "success",
		Message: "Cities fetched successfully",
		Data:    cities,
	})
}

// ------------------ GET CITY BY ID ------------------

func GetCityByID(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, APIResponse{
			Status:  "error",
			Message: "Invalid city ID format",
			Error:   err.Error(),
		})
		return
	}

	var city models.City
	if err := db.DB.Preload("Insurances").First(&city, id).Error; err != nil {
		c.JSON(http.StatusNotFound, APIResponse{
			Status:  "error",
			Message: "City not found",
		})
		return
	}

	c.JSON(http.StatusOK, APIResponse{
		Status:  "success",
		Message: "City fetched successfully",
		Data:    city,
	})
}

// ------------------ UPDATE CITY ------------------

func UpdateCity(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, APIResponse{
			Status:  "error",
			Message: "Invalid city ID format",
			Error:   err.Error(),
		})
		return
	}

	var city models.City
	if err := db.DB.First(&city, id).Error; err != nil {
		c.JSON(http.StatusNotFound, APIResponse{
			Status:  "error",
			Message: "City not found",
		})
		return
	}

	var input models.City
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, APIResponse{
			Status:  "error",
			Message: "Invalid input",
			Error:   err.Error(),
		})
		return
	}

	city.Name = input.Name

	if err := db.DB.Save(&city).Error; err != nil {
		c.JSON(http.StatusInternalServerError, APIResponse{
			Status:  "error",
			Message: "Failed to update city",
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, APIResponse{
		Status:  "success",
		Message: "City updated successfully",
		Data:    city,
	})
}

// ------------------ DELETE CITY ------------------

func DeleteCity(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, APIResponse{
			Status:  "error",
			Message: "Invalid city ID format",
			Error:   err.Error(),
		})
		return
	}

	var city models.City
	if err := db.DB.First(&city, id).Error; err != nil {
		c.JSON(http.StatusNotFound, APIResponse{
			Status:  "error",
			Message: "City not found",
		})
		return
	}

	if err := db.DB.Delete(&city).Error; err != nil {
		c.JSON(http.StatusInternalServerError, APIResponse{
			Status:  "error",
			Message: "Failed to delete city",
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, APIResponse{
		Status:  "success",
		Message: "City deleted successfully",
	})
}

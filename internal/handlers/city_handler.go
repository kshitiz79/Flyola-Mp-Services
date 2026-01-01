package handlers

import (
	"flyola-services/internal/models"
	"flyola-services/internal/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type CityHandler struct {
	cityService *services.CityService
}

func NewCityHandler(cityService *services.CityService) *CityHandler {
	return &CityHandler{cityService: cityService}
}

func (h *CityHandler) GetCities(c *gin.Context) {
	cities, err := h.cityService.GetAllCities()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to fetch cities",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Cities retrieved successfully",
		"data":    cities,
	})
}

func (h *CityHandler) GetCityByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid city ID",
		})
		return
	}

	city, err := h.cityService.GetCityByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "City not found",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "City retrieved successfully",
		"data":    city,
	})
}

func (h *CityHandler) CreateCity(c *gin.Context) {
	var city models.City
	if err := c.ShouldBindJSON(&city); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request data",
		})
		return
	}

	if err := h.cityService.CreateCity(&city); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to create city",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "City created successfully",
		"data":    city,
	})
}

func (h *CityHandler) UpdateCity(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid city ID",
		})
		return
	}

	var updates models.City
	if err := c.ShouldBindJSON(&updates); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request data",
		})
		return
	}

	city, err := h.cityService.UpdateCity(uint(id), &updates)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to update city",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "City updated successfully",
		"data":    city,
	})
}

func (h *CityHandler) DeleteCity(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid city ID",
		})
		return
	}

	if err := h.cityService.DeleteCity(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to delete city",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "City deleted successfully",
	})
}
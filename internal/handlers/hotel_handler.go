package handlers

import (
	"encoding/json"
	"flyola-services/internal/models"
	"flyola-services/internal/services"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type HotelHandler struct {
	hotelService *services.HotelService
}

func NewHotelHandler(hotelService *services.HotelService) *HotelHandler {
	return &HotelHandler{hotelService: hotelService}
}

func (h *HotelHandler) GetHotels(c *gin.Context) {
	hotels, err := h.hotelService.GetAllHotels()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to fetch hotels",
		})
		return
	}

	// Transform each hotel to frontend format
	transformedHotels := make([]map[string]interface{}, len(hotels))
	for i, hotel := range hotels {
		transformedHotels[i] = transformHotelResponse(&hotel)
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Hotels retrieved successfully",
		"data":    transformedHotels,
	})
}

// Helper function to transform hotel data for frontend
func transformHotelResponse(hotel *models.Hotel) map[string]interface{} {
	response := map[string]interface{}{
		"id":           hotel.ID,
		"name":         hotel.Name,
		"cityId":       hotel.CityID,
		"city":         hotel.City,
		"address":      hotel.Address,
		"description":  hotel.Description,
		"starRating":   hotel.StarRating,
		"imageUrl":     hotel.ImageURL,
		"contactPhone": hotel.ContactPhone,
		"contactEmail": hotel.ContactEmail,
		"checkInTime":  hotel.CheckInTime,
		"checkOutTime": hotel.CheckOutTime,
		"createdAt":    hotel.CreatedAt,
		"updatedAt":    hotel.UpdatedAt,
	}

	// Convert status integer to string
	if hotel.Status == 0 {
		response["status"] = "Active"
	} else {
		response["status"] = "Inactive"
	}

	// Parse images JSON string to array
	var images []interface{}
	if hotel.Images != "" {
		// Handle double-encoded JSON
		var imagesStr string
		if err := json.Unmarshal([]byte(hotel.Images), &imagesStr); err == nil {
			// It was double-encoded, now parse the actual JSON
			json.Unmarshal([]byte(imagesStr), &images)
		} else {
			// Try parsing directly
			json.Unmarshal([]byte(hotel.Images), &images)
		}
	}
	response["images"] = images

	// Parse amenities JSON string to array
	var amenities []interface{}
	if hotel.Amenities != "" {
		// Handle double-encoded JSON
		var amenitiesStr string
		if err := json.Unmarshal([]byte(hotel.Amenities), &amenitiesStr); err == nil {
			// It was double-encoded, now parse the actual JSON
			json.Unmarshal([]byte(amenitiesStr), &amenities)
		} else {
			// Try parsing directly
			json.Unmarshal([]byte(hotel.Amenities), &amenities)
		}
	}
	response["amenities"] = amenities

	return response
}

func (h *HotelHandler) GetHotelByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid hotel ID",
		})
		return
	}

	hotel, err := h.hotelService.GetHotelByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Hotel not found",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Hotel retrieved successfully",
		"data":    transformHotelResponse(hotel),
	})
}

func (h *HotelHandler) CreateHotel(c *gin.Context) {
	var requestData map[string]interface{}
	if err := c.ShouldBindJSON(&requestData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request data",
		})
		return
	}

	var hotel models.Hotel

	// Map the data to hotel model (same logic as before)
	if name, ok := requestData["name"].(string); ok {
		hotel.Name = name
	}
	if cityID, ok := requestData["cityId"].(float64); ok {
		hotel.CityID = uint(cityID)
	} else if cityIDStr, ok := requestData["cityId"].(string); ok {
		var cityIDInt int
		if _, err := fmt.Sscanf(cityIDStr, "%d", &cityIDInt); err == nil {
			hotel.CityID = uint(cityIDInt)
		}
	}
	if address, ok := requestData["address"].(string); ok {
		hotel.Address = address
	}
	if description, ok := requestData["description"].(string); ok {
		hotel.Description = description
	}
	if starRating, ok := requestData["starRating"].(float64); ok {
		hotel.StarRating = int(starRating)
	}
	if contactPhone, ok := requestData["contactPhone"].(string); ok {
		hotel.ContactPhone = contactPhone
	}
	if contactEmail, ok := requestData["contactEmail"].(string); ok {
		hotel.ContactEmail = contactEmail
	}
	if status, ok := requestData["status"].(string); ok {
		if status == "Active" {
			hotel.Status = 0
		} else {
			hotel.Status = 1
		}
	}

	// Handle amenities and images JSON
	if amenities, ok := requestData["amenities"]; ok {
		if amenitiesJSON, err := json.Marshal(amenities); err == nil {
			hotel.Amenities = string(amenitiesJSON)
		}
	}
	if images, ok := requestData["images"]; ok {
		if imagesJSON, err := json.Marshal(images); err == nil {
			hotel.Images = string(imagesJSON)
		}
	}

	if err := h.hotelService.CreateHotel(&hotel); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to create hotel",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Hotel created successfully",
		"data":    hotel,
	})
}

func (h *HotelHandler) UpdateHotel(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid hotel ID",
		})
		return
	}

	var requestData map[string]interface{}
	if err := c.ShouldBindJSON(&requestData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request data",
		})
		return
	}

	var updates models.Hotel

	// Map the data to hotel model (same logic as CreateHotel)
	if name, ok := requestData["name"].(string); ok {
		updates.Name = name
	}
	if cityID, ok := requestData["cityId"].(float64); ok {
		updates.CityID = uint(cityID)
	} else if cityIDStr, ok := requestData["cityId"].(string); ok {
		var cityIDInt int
		if _, err := fmt.Sscanf(cityIDStr, "%d", &cityIDInt); err == nil {
			updates.CityID = uint(cityIDInt)
		}
	}
	if address, ok := requestData["address"].(string); ok {
		updates.Address = address
	}
	if description, ok := requestData["description"].(string); ok {
		updates.Description = description
	}
	if starRating, ok := requestData["starRating"].(float64); ok {
		updates.StarRating = int(starRating)
	}
	if contactPhone, ok := requestData["contactPhone"].(string); ok {
		updates.ContactPhone = contactPhone
	}
	if contactEmail, ok := requestData["contactEmail"].(string); ok {
		updates.ContactEmail = contactEmail
	}
	if status, ok := requestData["status"].(string); ok {
		if status == "Active" {
			updates.Status = 0
		} else {
			updates.Status = 1
		}
	}

	// Handle amenities and images JSON
	if amenities, ok := requestData["amenities"]; ok {
		if amenitiesJSON, err := json.Marshal(amenities); err == nil {
			updates.Amenities = string(amenitiesJSON)
		}
	}
	if images, ok := requestData["images"]; ok {
		if imagesJSON, err := json.Marshal(images); err == nil {
			updates.Images = string(imagesJSON)
		}
	}

	hotel, err := h.hotelService.UpdateHotel(uint(id), &updates)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to update hotel",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Hotel updated successfully",
		"data":    hotel,
	})
}

func (h *HotelHandler) DeleteHotel(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid hotel ID",
		})
		return
	}

	if err := h.hotelService.DeleteHotel(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to delete hotel",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Hotel deleted successfully",
	})
}

func (h *HotelHandler) GetHotelsByCity(c *gin.Context) {
	cityID, err := strconv.ParseUint(c.Param("cityId"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid city ID",
		})
		return
	}

	hotels, err := h.hotelService.GetHotelsByCity(uint(cityID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to fetch hotels",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Hotels retrieved successfully",
		"data":    hotels,
	})
}

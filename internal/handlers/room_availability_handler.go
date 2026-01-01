package handlers

import (
	"flyola-services/internal/models"
	"flyola-services/internal/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type RoomAvailabilityHandler struct {
	roomAvailabilityService *services.RoomAvailabilityService
}

func NewRoomAvailabilityHandler(roomAvailabilityService *services.RoomAvailabilityService) *RoomAvailabilityHandler {
	return &RoomAvailabilityHandler{roomAvailabilityService: roomAvailabilityService}
}

func (h *RoomAvailabilityHandler) GetRoomAvailability(c *gin.Context) {
	roomID := c.Query("room_id")
	date := c.Query("date")

	availabilities, err := h.roomAvailabilityService.GetAllRoomAvailability(roomID, date)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to fetch room availability",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": availabilities,
	})
}

func (h *RoomAvailabilityHandler) GetRoomAvailabilityByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid availability ID",
		})
		return
	}

	availability, err := h.roomAvailabilityService.GetRoomAvailabilityByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Room availability not found",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": availability,
	})
}

func (h *RoomAvailabilityHandler) CreateRoomAvailability(c *gin.Context) {
	var availability models.RoomAvailability
	if err := c.ShouldBindJSON(&availability); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request data",
		})
		return
	}

	if err := h.roomAvailabilityService.CreateRoomAvailability(&availability); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to create room availability",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"data": availability,
	})
}

func (h *RoomAvailabilityHandler) UpdateRoomAvailability(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid availability ID",
		})
		return
	}

	var updates models.RoomAvailability
	if err := c.ShouldBindJSON(&updates); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request data",
		})
		return
	}

	availability, err := h.roomAvailabilityService.UpdateRoomAvailability(uint(id), &updates)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to update room availability",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": availability,
	})
}

func (h *RoomAvailabilityHandler) DeleteRoomAvailability(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid availability ID",
		})
		return
	}

	if err := h.roomAvailabilityService.DeleteRoomAvailability(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to delete room availability",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Room availability deleted successfully",
	})
}
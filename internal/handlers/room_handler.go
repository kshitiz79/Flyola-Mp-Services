package handlers

import (
	"flyola-services/internal/models"
	"flyola-services/internal/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type RoomHandler struct {
	roomService *services.RoomService
}

func NewRoomHandler(roomService *services.RoomService) *RoomHandler {
	return &RoomHandler{roomService: roomService}
}

func (h *RoomHandler) GetRooms(c *gin.Context) {
	rooms, err := h.roomService.GetAllRooms()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to fetch rooms",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Rooms retrieved successfully",
		"data":    rooms,
	})
}

func (h *RoomHandler) GetRoomByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid room ID",
		})
		return
	}

	room, err := h.roomService.GetRoomByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Room not found",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Room retrieved successfully",
		"data":    room,
	})
}

func (h *RoomHandler) CreateRoom(c *gin.Context) {
	var requestData map[string]interface{}
	if err := c.ShouldBindJSON(&requestData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request data",
		})
		return
	}

	var room models.Room

	// Map request data to room model
	if hotelID, ok := requestData["hotel_id"].(float64); ok {
		room.HotelID = uint(hotelID)
	}
	if categoryID, ok := requestData["room_category_id"].(float64); ok {
		room.RoomCategoryID = uint(categoryID)
	}
	if roomNumber, ok := requestData["room_number"].(string); ok {
		room.RoomNumber = roomNumber
	}
	if floor, ok := requestData["floor"].(float64); ok {
		room.Floor = int(floor)
	}
	if basePrice, ok := requestData["base_price"].(float64); ok {
		room.BasePrice = basePrice
	}
	if singlePrice, ok := requestData["single_price"].(float64); ok {
		room.SinglePrice = singlePrice
	}
	if doublePrice, ok := requestData["double_price"].(float64); ok {
		room.DoublePrice = doublePrice
	}
	if extraPersonPrice, ok := requestData["extra_person_price"].(float64); ok {
		room.ExtraPersonPrice = extraPersonPrice
	}
	if maxExtraPersons, ok := requestData["max_extra_persons"].(float64); ok {
		room.MaxExtraPersons = int(maxExtraPersons)
	}
	if status, ok := requestData["status"].(float64); ok {
		room.Status = int(status)
	}

	if err := h.roomService.CreateRoom(&room); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to create room",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Room created successfully",
		"data":    room,
	})
}

func (h *RoomHandler) UpdateRoom(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid room ID",
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

	var updates models.Room

	// Map request data to room model (same logic as create)
	if hotelID, ok := requestData["hotel_id"].(float64); ok {
		updates.HotelID = uint(hotelID)
	}
	if categoryID, ok := requestData["room_category_id"].(float64); ok {
		updates.RoomCategoryID = uint(categoryID)
	}
	if roomNumber, ok := requestData["room_number"].(string); ok {
		updates.RoomNumber = roomNumber
	}
	if floor, ok := requestData["floor"].(float64); ok {
		updates.Floor = int(floor)
	}
	if basePrice, ok := requestData["base_price"].(float64); ok {
		updates.BasePrice = basePrice
	}
	if singlePrice, ok := requestData["single_price"].(float64); ok {
		updates.SinglePrice = singlePrice
	}
	if doublePrice, ok := requestData["double_price"].(float64); ok {
		updates.DoublePrice = doublePrice
	}
	if extraPersonPrice, ok := requestData["extra_person_price"].(float64); ok {
		updates.ExtraPersonPrice = extraPersonPrice
	}
	if maxExtraPersons, ok := requestData["max_extra_persons"].(float64); ok {
		updates.MaxExtraPersons = int(maxExtraPersons)
	}
	if status, ok := requestData["status"].(float64); ok {
		updates.Status = int(status)
	}

	room, err := h.roomService.UpdateRoom(uint(id), &updates)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to update room",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Room updated successfully",
		"data":    room,
	})
}

func (h *RoomHandler) DeleteRoom(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid room ID",
		})
		return
	}

	if err := h.roomService.DeleteRoom(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to delete room",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Room deleted successfully",
	})
}

func (h *RoomHandler) GetRoomsByHotel(c *gin.Context) {
	hotelID, err := strconv.ParseUint(c.Param("hotelId"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid hotel ID",
		})
		return
	}

	rooms, err := h.roomService.GetRoomsByHotel(uint(hotelID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to fetch rooms",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Rooms retrieved successfully",
		"data":    rooms,
	})
}
package handlers

import (
	"flyola-services/internal/models"
	"flyola-services/internal/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type RoomCategoryHandler struct {
	roomCategoryService *services.RoomCategoryService
}

func NewRoomCategoryHandler(roomCategoryService *services.RoomCategoryService) *RoomCategoryHandler {
	return &RoomCategoryHandler{roomCategoryService: roomCategoryService}
}

func (h *RoomCategoryHandler) GetRoomCategories(c *gin.Context) {
	categories, err := h.roomCategoryService.GetAllRoomCategories()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to fetch room categories",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": categories,
	})
}

func (h *RoomCategoryHandler) GetRoomCategoryByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid category ID",
		})
		return
	}

	category, err := h.roomCategoryService.GetRoomCategoryByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Room category not found",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": category,
	})
}

func (h *RoomCategoryHandler) CreateRoomCategory(c *gin.Context) {
	var category models.RoomCategory
	if err := c.ShouldBindJSON(&category); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request data",
		})
		return
	}

	if err := h.roomCategoryService.CreateRoomCategory(&category); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to create room category",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"data": category,
	})
}

func (h *RoomCategoryHandler) UpdateRoomCategory(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid category ID",
		})
		return
	}

	var updates models.RoomCategory
	if err := c.ShouldBindJSON(&updates); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request data",
		})
		return
	}

	category, err := h.roomCategoryService.UpdateRoomCategory(uint(id), &updates)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to update room category",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": category,
	})
}

func (h *RoomCategoryHandler) DeleteRoomCategory(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid category ID",
		})
		return
	}

	if err := h.roomCategoryService.DeleteRoomCategory(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to delete room category",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Room category deleted successfully",
	})
}
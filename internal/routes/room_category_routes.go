package routes

import (
	"flyola-services/internal/handlers"

	"github.com/gin-gonic/gin"
)

func SetupRoomCategoryRoutes(rg *gin.RouterGroup, handler *handlers.RoomCategoryHandler) {
	roomCategories := rg.Group("/room-categories")
	{
		roomCategories.GET("", handler.GetRoomCategories)
		roomCategories.GET("/:id", handler.GetRoomCategoryByID)
		roomCategories.POST("", handler.CreateRoomCategory)
		roomCategories.PUT("/:id", handler.UpdateRoomCategory)
		roomCategories.DELETE("/:id", handler.DeleteRoomCategory)
	}
}
package routes

import (
	"flyola-services/internal/handlers"

	"github.com/gin-gonic/gin"
)

func SetupRoomAvailabilityRoutes(rg *gin.RouterGroup, handler *handlers.RoomAvailabilityHandler) {
	roomAvailability := rg.Group("/room-availability")
	{
		roomAvailability.GET("", handler.GetRoomAvailability)
		roomAvailability.GET("/:id", handler.GetRoomAvailabilityByID)
		roomAvailability.POST("", handler.CreateRoomAvailability)
		roomAvailability.PUT("/:id", handler.UpdateRoomAvailability)
		roomAvailability.DELETE("/:id", handler.DeleteRoomAvailability)
	}
}
package routes

import (
	"flyola-services/internal/handlers"

	"github.com/gin-gonic/gin"
)

func SetupRoomRoutes(router *gin.RouterGroup, roomHandler *handlers.RoomHandler) {
	rooms := router.Group("/rooms")
	{
		rooms.GET("", roomHandler.GetRooms)
		rooms.POST("", roomHandler.CreateRoom)
		rooms.GET("/:id", roomHandler.GetRoomByID)
		rooms.PUT("/:id", roomHandler.UpdateRoom)
		rooms.DELETE("/:id", roomHandler.DeleteRoom)
		rooms.GET("/hotel/:hotelId", roomHandler.GetRoomsByHotel)
	}
}
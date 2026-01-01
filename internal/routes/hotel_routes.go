package routes

import (
	"flyola-services/internal/handlers"

	"github.com/gin-gonic/gin"
)

func SetupHotelRoutes(router *gin.RouterGroup, hotelHandler *handlers.HotelHandler) {
	hotels := router.Group("/hotels")
	{
		hotels.GET("", hotelHandler.GetHotels)
		hotels.POST("", hotelHandler.CreateHotel)
		hotels.GET("/:id", hotelHandler.GetHotelByID)
		hotels.PUT("/:id", hotelHandler.UpdateHotel)
		hotels.DELETE("/:id", hotelHandler.DeleteHotel)
		hotels.GET("/city/:cityId", hotelHandler.GetHotelsByCity)
	}
}
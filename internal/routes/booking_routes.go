package routes

import (
	"flyola-services/internal/handlers"

	"github.com/gin-gonic/gin"
)

func SetupBookingRoutes(router *gin.RouterGroup, bookingHandler *handlers.BookingHandler) {
	bookings := router.Group("/bookings")
	{
		bookings.GET("", bookingHandler.GetBookings)
		bookings.POST("", bookingHandler.CreateBooking)
		bookings.GET("/:id", bookingHandler.GetBookingByID)
		bookings.PUT("/:id", bookingHandler.UpdateBooking)
		bookings.DELETE("/:id", bookingHandler.DeleteBooking)
		bookings.PUT("/:id/cancel", bookingHandler.CancelBooking)
	}
}
package routes

import (
	"flyola-services/internal/handlers"

	"github.com/gin-gonic/gin"
)

func SetupHolidayPackageRoutes(router *gin.RouterGroup, holidayPackageHandler *handlers.HolidayPackageHandler) {
	packages := router.Group("/holiday-packages")
	{
		// Package management routes
		packages.GET("", holidayPackageHandler.GetAllPackages)
		packages.GET("/:id", holidayPackageHandler.GetPackageByID)
		packages.GET("/type/:type", holidayPackageHandler.GetPackagesByType)
		packages.GET("/date/:date", holidayPackageHandler.GetPackagesByDate)
		packages.POST("", holidayPackageHandler.CreatePackage) // Admin only
		packages.PUT("/:id", holidayPackageHandler.UpdatePackage) // Admin only
		packages.DELETE("/:id", holidayPackageHandler.DeletePackage) // Admin only

		// Booking routes
		packages.POST("/book", holidayPackageHandler.CreatePackageBooking)
		packages.POST("/book/:id/confirm", holidayPackageHandler.ConfirmPackageBooking)
		
		// Booking management routes
		packages.GET("/bookings/:id", holidayPackageHandler.GetBookingByID)
		packages.GET("/bookings/reference/:reference", holidayPackageHandler.GetBookingByReference)
		packages.DELETE("/bookings/:id", holidayPackageHandler.CancelPackageBooking)
		
		// Admin routes
		packages.GET("/admin/bookings", holidayPackageHandler.GetAllPackageBookings) // Admin only
	}
}
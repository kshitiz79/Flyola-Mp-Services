package handlers

import (
	"flyola-services/internal/models"
	"flyola-services/internal/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type BookingHandler struct {
	bookingService *services.BookingService
}

func NewBookingHandler(bookingService *services.BookingService) *BookingHandler {
	return &BookingHandler{bookingService: bookingService}
}

func (h *BookingHandler) GetBookings(c *gin.Context) {
	// Check if email query parameter is provided
	email := c.Query("email")

	var bookings []models.HotelBooking
	var err error

	if email != "" {
		// Filter by guest email if provided
		bookings, err = h.bookingService.GetBookingsByGuestEmail(email)
	} else {
		// Get all bookings if no email filter
		bookings, err = h.bookingService.GetAllBookings()
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch bookings"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Bookings retrieved successfully", "data": bookings})
}

func (h *BookingHandler) GetBookingByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid booking ID"})
		return
	}

	booking, err := h.bookingService.GetBookingByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Booking not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Booking retrieved successfully", "data": booking})
}

func (h *BookingHandler) CreateBooking(c *gin.Context) {
	var booking models.HotelBooking
	if err := c.ShouldBindJSON(&booking); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data", "details": err.Error()})
		return
	}

	if err := h.bookingService.CreateBooking(&booking); err != nil {
		// Log the actual error for debugging
		println("Error creating booking:", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create booking", "details": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "Booking created successfully", "data": booking})
}

func (h *BookingHandler) UpdateBooking(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid booking ID"})
		return
	}

	var updates models.HotelBooking
	if err := c.ShouldBindJSON(&updates); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
		return
	}

	// Security Check: Prevent manual confirmation/paid status updates via general PUT
	// These should only be updated via VerifyPayment endpoint
	if updates.BookingStatus == "confirmed" || updates.PaymentStatus == "paid" || updates.PaymentStatus == "completed" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Direct status confirmation not allowed. Use payment verification flow."})
		return
	}

	booking, err := h.bookingService.UpdateBooking(uint(id), &updates)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update booking"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Booking updated successfully", "data": booking})
}

func (h *BookingHandler) DeleteBooking(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid booking ID"})
		return
	}

	if err := h.bookingService.DeleteBooking(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete booking"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Booking deleted successfully"})
}

func (h *BookingHandler) CancelBooking(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid booking ID"})
		return
	}

	booking, err := h.bookingService.CancelBooking(uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to cancel booking"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Booking cancelled successfully", "data": booking})
}

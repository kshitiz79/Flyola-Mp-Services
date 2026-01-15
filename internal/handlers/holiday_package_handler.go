package handlers

import (
	"flyola-services/internal/models"
	"flyola-services/internal/services"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type HolidayPackageHandler struct {
	service *services.HolidayPackageService
}

func NewHolidayPackageHandler(service *services.HolidayPackageService) *HolidayPackageHandler {
	return &HolidayPackageHandler{
		service: service,
	}
}

// GetAllPackages handles GET /api/v1/holiday-packages
func (h *HolidayPackageHandler) GetAllPackages(c *gin.Context) {
	packages, err := h.service.GetAllPackages()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Failed to fetch packages: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    packages,
	})
}

// GetPackageByID handles GET /api/v1/holiday-packages/{id}
func (h *HolidayPackageHandler) GetPackageByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Invalid package ID",
		})
		return
	}

	pkg, err := h.service.GetPackageByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"error":   "Package not found: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    pkg,
	})
}

// GetPackagesByType handles GET /api/v1/holiday-packages/type/{type}
func (h *HolidayPackageHandler) GetPackagesByType(c *gin.Context) {
	packageType := c.Param("type")

	packages, err := h.service.GetPackagesByType(packageType)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Failed to fetch packages: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    packages,
	})
}

// GetPackagesByDate handles GET /api/v1/holiday-packages/date/{date}
func (h *HolidayPackageHandler) GetPackagesByDate(c *gin.Context) {
	dateStr := c.Param("date")

	// Validate date format
	_, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Invalid date format, expected YYYY-MM-DD",
		})
		return
	}

	packages, err := h.service.GetPackagesByDate(dateStr)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Failed to fetch packages: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    packages,
	})
}

// CreatePackageBooking handles POST /api/v1/holiday-packages/book
func (h *HolidayPackageHandler) CreatePackageBooking(c *gin.Context) {
	var req struct {
		PackageID       uint                        `json:"package_id"`
		GuestName       string                      `json:"guest_name"`
		GuestEmail      string                      `json:"guest_email"`
		GuestPhone      string                      `json:"guest_phone"`
		TravelDate      string                      `json:"travel_date"`
		SpecialRequests string                      `json:"special_requests"`
		Passengers      []models.PackagePassenger   `json:"passengers"`
		PaymentID       string                      `json:"payment_id"`
		PaymentMethod   string                      `json:"payment_method"`
		PaymentStatus   string                      `json:"payment_status"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Invalid request body: " + err.Error(),
		})
		return
	}

	// Validate required fields
	if req.PackageID == 0 || req.GuestName == "" || req.GuestEmail == "" || req.GuestPhone == "" || req.TravelDate == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Missing required fields",
		})
		return
	}

	if len(req.Passengers) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "At least one passenger is required",
		})
		return
	}

	// Parse travel date
	travelDate, err := time.Parse("2006-01-02", req.TravelDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Invalid travel date format (YYYY-MM-DD)",
		})
		return
	}

	// Get package details to calculate total amount
	pkg, err := h.service.GetPackageByID(req.PackageID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"error":   "Package not found: " + err.Error(),
		})
		return
	}

	// Calculate total amount
	totalAmount := pkg.PricePerPerson * float64(len(req.Passengers))

	// Determine payment and booking status
	paymentStatus := "pending"
	bookingStatus := "pending"
	if req.PaymentID != "" && req.PaymentStatus == "paid" {
		paymentStatus = "paid"
		bookingStatus = "confirmed"
	}

	log.Printf("🔍 Payment Info - PaymentID: %s, PaymentStatus: %s, PaymentMethod: %s", req.PaymentID, req.PaymentStatus, req.PaymentMethod)
	log.Printf("🔍 Setting booking status: %s, payment status: %s", bookingStatus, paymentStatus)

	// Create booking
	booking := &models.PackageBooking{
		PackageID:       req.PackageID,
		GuestName:       req.GuestName,
		GuestEmail:      req.GuestEmail,
		GuestPhone:      req.GuestPhone,
		NumPassengers:   len(req.Passengers),
		TravelDate:      travelDate,
		TotalAmount:     totalAmount,
		SpecialRequests: req.SpecialRequests,
		Passengers:      req.Passengers,
		PaymentID:       req.PaymentID,
		PaymentMethod:   req.PaymentMethod,
		PaymentStatus:   paymentStatus,
		BookingStatus:   bookingStatus,
	}

	log.Printf("🔍 Booking before save - PaymentID: %s, PaymentStatus: %s, BookingStatus: %s", booking.PaymentID, booking.PaymentStatus, booking.BookingStatus)

	if err := h.service.CreatePackageBooking(booking); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Failed to create booking: " + err.Error(),
		})
		return
	}

	// If payment is successful, book individual schedules through Node.js backend
	if booking.PaymentStatus == "paid" && booking.BookingStatus == "confirmed" {
		log.Printf("🎫 Booking individual schedules for package booking ID: %d", booking.ID)
		if err := h.service.BookPackageSchedules(booking.ID); err != nil {
			log.Printf("⚠️ Warning: Failed to book individual schedules: %v", err)
			// Don't fail the entire booking, just log the error
			// The booking is still created, but schedules need to be booked manually
		} else {
			log.Printf("✅ Successfully booked individual schedules for booking ID: %d", booking.ID)
		}
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"data":    booking,
		"message": "Package booking created successfully",
	})
}

// ConfirmPackageBooking handles POST /api/v1/holiday-packages/book/{id}/confirm
func (h *HolidayPackageHandler) ConfirmPackageBooking(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Invalid booking ID",
		})
		return
	}

	var req struct {
		PaymentID     string `json:"payment_id"`
		PaymentMethod string `json:"payment_method"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Invalid request body: " + err.Error(),
		})
		return
	}

	// Update payment status
	if err := h.service.UpdateBookingPaymentStatus(uint(id), req.PaymentID, req.PaymentMethod); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Failed to update payment status: " + err.Error(),
		})
		return
	}

	// Book individual schedules through Node.js backend
	if err := h.service.BookPackageSchedules(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Failed to book package schedules: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Package booking confirmed and schedules booked successfully",
	})
}

// GetBookingByID handles GET /api/v1/holiday-packages/bookings/{id}
func (h *HolidayPackageHandler) GetBookingByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Invalid booking ID",
		})
		return
	}

	booking, err := h.service.GetBookingByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"error":   "Booking not found: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    booking,
	})
}

// GetBookingByReference handles GET /api/v1/holiday-packages/bookings/reference/{reference}
func (h *HolidayPackageHandler) GetBookingByReference(c *gin.Context) {
	reference := c.Param("reference")

	booking, err := h.service.GetBookingByReference(reference)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"error":   "Booking not found: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    booking,
	})
}

// CancelPackageBooking handles DELETE /api/v1/holiday-packages/bookings/{id}
func (h *HolidayPackageHandler) CancelPackageBooking(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Invalid booking ID",
		})
		return
	}

	if err := h.service.CancelPackageBooking(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Failed to cancel booking: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Package booking cancelled successfully",
	})
}

// CreatePackage handles POST /api/v1/holiday-packages (Admin only)
func (h *HolidayPackageHandler) CreatePackage(c *gin.Context) {
	var pkg models.HolidayPackage

	if err := c.ShouldBindJSON(&pkg); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Invalid request body: " + err.Error(),
		})
		return
	}

	// Validate required fields
	if pkg.Title == "" || pkg.PricePerPerson <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Missing required fields: title and price_per_person",
		})
		return
	}

	// Create package in database
	if err := h.service.CreatePackage(&pkg); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Failed to create package: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"data":    pkg,
		"message": "Package created successfully",
	})
}

// UpdatePackage handles PUT /api/v1/holiday-packages/{id} (Admin only)
func (h *HolidayPackageHandler) UpdatePackage(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Invalid package ID",
		})
		return
	}

	var pkg models.HolidayPackage
	if err := c.ShouldBindJSON(&pkg); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Invalid request body: " + err.Error(),
		})
		return
	}

	// Validate required fields
	if pkg.Title == "" || pkg.PricePerPerson <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Missing required fields: title and price_per_person",
		})
		return
	}

	// Set the ID for update
	pkg.ID = uint(id)

	// If status is not provided, preserve the existing status
	if pkg.Status == 0 {
		existingPkg, err := h.service.GetPackageByIDWithoutStatusFilter(uint(id))
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{
				"success": false,
				"error":   "Package not found: " + err.Error(),
			})
			return
		}
		pkg.Status = existingPkg.Status
	}

	// Update package in database
	if err := h.service.UpdatePackage(&pkg); err != nil {
		// Check if it's the specific error about existing bookings
		if err.Error() == "package updated successfully, but schedules cannot be modified as there are existing bookings" {
			c.JSON(http.StatusOK, gin.H{
				"success": true,
				"data":    pkg,
				"message": "Package updated successfully",
				"warning": "Schedules were not modified because there are existing bookings for this package",
			})
			return
		}
		
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Failed to update package: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    pkg,
		"message": "Package updated successfully",
	})
}

// DeletePackage handles DELETE /api/v1/holiday-packages/{id} (Admin only)
func (h *HolidayPackageHandler) DeletePackage(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Invalid package ID",
		})
		return
	}

	if err := h.service.DeletePackage(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Failed to delete package: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Package deleted successfully",
	})
}

// GetAllPackageBookings handles GET /api/v1/holiday-packages/admin/bookings (Admin only)
func (h *HolidayPackageHandler) GetAllPackageBookings(c *gin.Context) {
	bookings, err := h.service.GetAllPackageBookings()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Failed to fetch bookings: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    bookings,
	})
}
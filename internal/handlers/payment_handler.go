package handlers

import (
	"flyola-services/internal/models"
	"flyola-services/internal/services"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type PaymentHandler struct {
	paymentService *services.PaymentService
	bookingService *services.BookingService
	razorpayID     string
	razorpaySecret string
}

func NewPaymentHandler(paymentService *services.PaymentService, bookingService *services.BookingService, razorpayID, razorpaySecret string) *PaymentHandler {
	return &PaymentHandler{
		paymentService: paymentService,
		bookingService: bookingService,
		razorpayID:     razorpayID,
		razorpaySecret: razorpaySecret,
	}
}

func (h *PaymentHandler) CreateOrder(c *gin.Context) {
	var req struct {
		Amount   int64  `json:"amount" binding:"required"`
		Currency string `json:"currency" binding:"required"`
		Receipt  string `json:"receipt" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data", "details": err.Error()})
		return
	}

	// Debug logging
	log.Printf("🔍 Creating Razorpay order - Amount: %d, Currency: %s, Receipt: %s", req.Amount, req.Currency, req.Receipt)
	log.Printf("🔍 Using Razorpay Key ID: %s", h.razorpayID)
	if h.razorpaySecret != "" {
		log.Printf("🔍 Razorpay Secret is set (length: %d)", len(h.razorpaySecret))
	} else {
		log.Printf("🔍 Razorpay Secret is NOT set")
	}

	orderID, err := h.paymentService.CreateRazorpayOrder(req.Amount, req.Currency, req.Receipt, h.razorpayID, h.razorpaySecret)
	if err != nil {
		log.Printf("❌ Razorpay order creation failed: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create Razorpay order", "details": err.Error()})
		return
	}

	log.Printf("✅ Razorpay order created successfully: %s", orderID)
	c.JSON(http.StatusOK, gin.H{
		"order_id": orderID,
	})
}

func (h *PaymentHandler) VerifyPayment(c *gin.Context) {
	var req struct {
		BookingID         uint   `json:"booking_id" binding:"required"`
		RazorpayOrderID   string `json:"razorpay_order_id" binding:"required"`
		RazorpayPaymentID string `json:"razorpay_payment_id" binding:"required"`
		RazorpaySignature string `json:"razorpay_signature" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data", "details": err.Error()})
		return
	}

	// 1. Verify Signature
	err := h.paymentService.VerifyRazorpaySignature(req.RazorpayOrderID, req.RazorpayPaymentID, req.RazorpaySignature, h.razorpaySecret)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Payment verification failed", "details": err.Error()})
		return
	}

	// 2. Update Booking Status
	updates := &models.HotelBooking{
		BookingStatus: "confirmed",
		PaymentStatus: "paid",
		PaymentID:     req.RazorpayPaymentID,
		PaymentMethod: "razorpay",
	}

	booking, err := h.bookingService.UpdateBooking(req.BookingID, updates)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update booking status", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Payment verified and booking confirmed",
		"data":    booking,
	})
}

func (h *PaymentHandler) ProcessPayment(c *gin.Context) {
	var payment models.HotelPayment
	if err := c.ShouldBindJSON(&payment); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
		return
	}

	if err := h.paymentService.ProcessPayment(&payment); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process payment"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Payment processed successfully", "data": payment})
}

func (h *PaymentHandler) GetPaymentByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid payment ID"})
		return
	}

	payment, err := h.paymentService.GetPaymentByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Payment not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Payment retrieved successfully", "data": payment})
}

func (h *PaymentHandler) GetPaymentByBooking(c *gin.Context) {
	bookingID, err := strconv.ParseUint(c.Param("bookingId"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid booking ID"})
		return
	}

	payment, err := h.paymentService.GetPaymentByBooking(uint(bookingID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Payment not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Payment retrieved successfully", "data": payment})
}

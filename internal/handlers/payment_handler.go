package handlers

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
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

// CreateOrder creates a new Razorpay order (for both hotel and package bookings)
func (h *PaymentHandler) CreateOrder(c *gin.Context) {
	var req struct {
		Amount   interface{}            `json:"amount" binding:"required"`
		Currency string                 `json:"currency" binding:"required"`
		Receipt  string                 `json:"receipt" binding:"required"`
		Notes    map[string]interface{} `json:"notes"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data", "details": err.Error()})
		return
	}

	// Convert amount to int64
	var amount int64
	switch v := req.Amount.(type) {
	case float64:
		amount = int64(v)
	case int:
		amount = int64(v)
	case int64:
		amount = v
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid amount type"})
		return
	}

	// Debug logging
	log.Printf("🔍 Creating Razorpay order - Amount: %d, Currency: %s, Receipt: %s", amount, req.Currency, req.Receipt)
	log.Printf("🔍 Using Razorpay Key ID: %s", h.razorpayID)
	if h.razorpaySecret != "" {
		log.Printf("🔍 Razorpay Secret is set (length: %d)", len(h.razorpaySecret))
	} else {
		log.Printf("🔍 Razorpay Secret is NOT set")
	}

	orderID, err := h.paymentService.CreateRazorpayOrder(amount, req.Currency, req.Receipt, h.razorpayID, h.razorpaySecret)
	if err != nil {
		log.Printf("❌ Razorpay order creation failed: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create Razorpay order", "details": err.Error()})
		return
	}

	log.Printf("✅ Razorpay order created successfully: %s", orderID)
	c.JSON(http.StatusOK, gin.H{
		"id":       orderID,
		"entity":   "order",
		"amount":   amount,
		"currency": req.Currency,
		"receipt":  req.Receipt,
		"status":   "created",
		"notes":    req.Notes,
	})
}

// VerifyPayment verifies the Razorpay payment signature (for both hotel and package bookings)
func (h *PaymentHandler) VerifyPayment(c *gin.Context) {
	var req struct {
		BookingID         *uint  `json:"booking_id"` // Optional for hotel bookings
		RazorpayOrderID   string `json:"razorpay_order_id" binding:"required"`
		RazorpayPaymentID string `json:"razorpay_payment_id" binding:"required"`
		RazorpaySignature string `json:"razorpay_signature" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data", "details": err.Error()})
		return
	}

	// Verify signature using HMAC SHA256
	verified := h.verifyRazorpaySignature(req.RazorpayOrderID, req.RazorpayPaymentID, req.RazorpaySignature)
	
	if !verified {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success":  false,
			"verified": false,
			"error":    "Payment verification failed - invalid signature",
		})
		return
	}

	// If booking_id is provided, update the hotel booking status
	if req.BookingID != nil {
		updates := &models.HotelBooking{
			BookingStatus: "confirmed",
			PaymentStatus: "paid",
			PaymentID:     req.RazorpayPaymentID,
			PaymentMethod: "razorpay",
		}

		booking, err := h.bookingService.UpdateBooking(*req.BookingID, updates)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update booking status", "details": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"success":  true,
			"verified": true,
			"message":  "Payment verified and booking confirmed",
			"data":     booking,
		})
		return
	}

	// For package bookings or standalone verification
	c.JSON(http.StatusOK, gin.H{
		"success":  true,
		"verified": true,
		"message":  "Payment verification successful",
	})
}

// verifyRazorpaySignature verifies the Razorpay payment signature
func (h *PaymentHandler) verifyRazorpaySignature(orderID, paymentID, signature string) bool {
	// Create the expected signature
	message := orderID + "|" + paymentID
	
	// Create HMAC SHA256 hash
	hmacHash := hmac.New(sha256.New, []byte(h.razorpaySecret))
	hmacHash.Write([]byte(message))
	expectedSignature := hex.EncodeToString(hmacHash.Sum(nil))

	// Compare signatures
	return hmac.Equal([]byte(signature), []byte(expectedSignature))
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

package routes

import (
	"flyola-services/internal/handlers"

	"github.com/gin-gonic/gin"
)

func SetupPaymentRoutes(router *gin.RouterGroup, paymentHandler *handlers.PaymentHandler) {
	payments := router.Group("/payments")
	{
		payments.POST("/process", paymentHandler.ProcessPayment)
		payments.POST("/create-order", paymentHandler.CreateOrder)
		payments.POST("/verify", paymentHandler.VerifyPayment)
		payments.GET("/booking/:bookingId", paymentHandler.GetPaymentByBooking)
		payments.GET("/:id", paymentHandler.GetPaymentByID)
	}
}

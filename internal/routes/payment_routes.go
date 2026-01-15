package routes

import (
	"flyola-services/internal/handlers"

	"github.com/gin-gonic/gin"
)

func SetupPaymentRoutes(router *gin.RouterGroup, paymentHandler *handlers.PaymentHandler) {
	payments := router.Group("/payments")
	{
		// Payment processing routes
		payments.POST("/create-order", paymentHandler.CreateOrder)
		payments.POST("/verify", paymentHandler.VerifyPayment)
	}
}
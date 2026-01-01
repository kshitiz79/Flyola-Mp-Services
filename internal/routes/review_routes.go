package routes

import (
	"flyola-services/internal/handlers"

	"github.com/gin-gonic/gin"
)

func SetupReviewRoutes(router *gin.RouterGroup, reviewHandler *handlers.ReviewHandler) {
	reviews := router.Group("/reviews")
	{
		reviews.GET("", reviewHandler.GetAllReviews)
		reviews.GET("/hotel/:hotelId", reviewHandler.GetHotelReviews)
		reviews.POST("", reviewHandler.CreateReview)
		reviews.PUT("/:id/status", reviewHandler.UpdateReviewStatus)
		reviews.DELETE("/:id", reviewHandler.DeleteReview)
	}
}
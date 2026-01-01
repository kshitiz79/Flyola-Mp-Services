package handlers

import (
	"flyola-services/internal/models"
	"flyola-services/internal/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ReviewHandler struct {
	reviewService *services.ReviewService
}

func NewReviewHandler(reviewService *services.ReviewService) *ReviewHandler {
	return &ReviewHandler{reviewService: reviewService}
}

func (h *ReviewHandler) GetAllReviews(c *gin.Context) {
	reviews, err := h.reviewService.GetAllReviews()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch reviews"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Reviews retrieved successfully", "data": reviews})
}

func (h *ReviewHandler) GetHotelReviews(c *gin.Context) {
	hotelID, err := strconv.ParseUint(c.Param("hotelId"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid hotel ID"})
		return
	}

	reviews, err := h.reviewService.GetHotelReviews(uint(hotelID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch reviews"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Hotel reviews retrieved successfully", "data": reviews})
}

func (h *ReviewHandler) CreateReview(c *gin.Context) {
	var review models.HotelReview
	if err := c.ShouldBindJSON(&review); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
		return
	}

	if err := h.reviewService.CreateReview(&review); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create review"})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "Review created successfully", "data": review})
}

func (h *ReviewHandler) UpdateReviewStatus(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid review ID"})
		return
	}

	var statusData struct {
		Status int `json:"status"`
	}
	if err := c.ShouldBindJSON(&statusData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
		return
	}

	review, err := h.reviewService.UpdateReviewStatus(uint(id), statusData.Status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update review status"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Review status updated successfully", "data": review})
}

func (h *ReviewHandler) DeleteReview(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid review ID"})
		return
	}

	if err := h.reviewService.DeleteReview(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete review"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Review deleted successfully"})
}
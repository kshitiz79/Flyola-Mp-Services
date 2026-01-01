package handlers

import (
	"flyola-services/internal/models"
	"flyola-services/internal/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type MealPlanHandler struct {
	mealPlanService *services.MealPlanService
}

func NewMealPlanHandler(mealPlanService *services.MealPlanService) *MealPlanHandler {
	return &MealPlanHandler{mealPlanService: mealPlanService}
}

func (h *MealPlanHandler) GetMealPlans(c *gin.Context) {
	mealPlans, err := h.mealPlanService.GetAllMealPlans()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch meal plans"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": mealPlans})
}

func (h *MealPlanHandler) GetMealPlanByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid meal plan ID"})
		return
	}

	mealPlan, err := h.mealPlanService.GetMealPlanByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Meal plan not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": mealPlan})
}

func (h *MealPlanHandler) CreateMealPlan(c *gin.Context) {
	var mealPlan models.MealPlan
	if err := c.ShouldBindJSON(&mealPlan); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
		return
	}

	if err := h.mealPlanService.CreateMealPlan(&mealPlan); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create meal plan"})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"data": mealPlan})
}

func (h *MealPlanHandler) UpdateMealPlan(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid meal plan ID"})
		return
	}

	var updates models.MealPlan
	if err := c.ShouldBindJSON(&updates); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
		return
	}

	mealPlan, err := h.mealPlanService.UpdateMealPlan(uint(id), &updates)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update meal plan"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": mealPlan})
}

func (h *MealPlanHandler) DeleteMealPlan(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid meal plan ID"})
		return
	}

	if err := h.mealPlanService.DeleteMealPlan(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete meal plan"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Meal plan deleted successfully"})
}
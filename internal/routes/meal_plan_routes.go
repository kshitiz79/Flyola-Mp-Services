package routes

import (
	"flyola-services/internal/handlers"

	"github.com/gin-gonic/gin"
)

func SetupMealPlanRoutes(router *gin.RouterGroup, mealPlanHandler *handlers.MealPlanHandler) {
	mealPlans := router.Group("/meal-plans")
	{
		mealPlans.GET("", mealPlanHandler.GetMealPlans)
		mealPlans.POST("", mealPlanHandler.CreateMealPlan)
		mealPlans.GET("/:id", mealPlanHandler.GetMealPlanByID)
		mealPlans.PUT("/:id", mealPlanHandler.UpdateMealPlan)
		mealPlans.DELETE("/:id", mealPlanHandler.DeleteMealPlan)
	}
}
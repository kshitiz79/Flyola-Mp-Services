package routes

import (
	"flyola-services/internal/handlers"

	"github.com/gin-gonic/gin"
)

func SetupCityRoutes(router *gin.RouterGroup, cityHandler *handlers.CityHandler) {
	cities := router.Group("/cities")
	{
		cities.GET("", cityHandler.GetCities)
		cities.POST("", cityHandler.CreateCity)
		cities.GET("/:id", cityHandler.GetCityByID)
		cities.PUT("/:id", cityHandler.UpdateCity)
		cities.DELETE("/:id", cityHandler.DeleteCity)
	}
}
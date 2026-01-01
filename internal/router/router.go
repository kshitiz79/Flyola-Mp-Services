package router

import (
	"flyola-services/internal/config"
	"flyola-services/internal/handlers"
	"flyola-services/internal/middleware"
	"flyola-services/internal/routes"
	"flyola-services/internal/services"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func Initialize(db *gorm.DB, cfg *config.Config) *gin.Engine {
	r := gin.Default()

	// Add middleware
	r.Use(middleware.CORS())
	r.Use(middleware.Logger())

	// Initialize services
	cityService := services.NewCityService(db)
	hotelService := services.NewHotelService(db)
	roomService := services.NewRoomService(db)
	roomCategoryService := services.NewRoomCategoryService(db)
	roomAvailabilityService := services.NewRoomAvailabilityService(db)
	mealPlanService := services.NewMealPlanService(db)
	bookingService := services.NewBookingService(db)
	paymentService := services.NewPaymentService(db)
	reviewService := services.NewReviewService(db)

	// Initialize handlers
	cityHandler := handlers.NewCityHandler(cityService)
	hotelHandler := handlers.NewHotelHandler(hotelService)
	roomHandler := handlers.NewRoomHandler(roomService)
	roomCategoryHandler := handlers.NewRoomCategoryHandler(roomCategoryService)
	roomAvailabilityHandler := handlers.NewRoomAvailabilityHandler(roomAvailabilityService)
	mealPlanHandler := handlers.NewMealPlanHandler(mealPlanService)
	bookingHandler := handlers.NewBookingHandler(bookingService)
	paymentHandler := handlers.NewPaymentHandler(paymentService, bookingService, cfg.RazorpayID, cfg.RazorpaySecret)
	reviewHandler := handlers.NewReviewHandler(reviewService)

	// Health check endpoint
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "ok",
			"message": "Flyola Hotel Services Backend is running",
		})
	})

	// API v1 routes
	v1 := r.Group("/api/v1")
	{
		// Setup all routes using separate route files
		routes.SetupCityRoutes(v1, cityHandler)
		routes.SetupHotelRoutes(v1, hotelHandler)
		routes.SetupRoomRoutes(v1, roomHandler)
		routes.SetupRoomCategoryRoutes(v1, roomCategoryHandler)
		routes.SetupRoomAvailabilityRoutes(v1, roomAvailabilityHandler)
		routes.SetupMealPlanRoutes(v1, mealPlanHandler)
		routes.SetupBookingRoutes(v1, bookingHandler)
		routes.SetupPaymentRoutes(v1, paymentHandler)
		routes.SetupReviewRoutes(v1, reviewHandler)
	}

	return r
}

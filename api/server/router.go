package server

import (
	"database/sql"

	"github.com/kzkei/natureAdvice/api/internal/handlers"
	"github.com/kzkei/natureAdvice/api/internal/repositories"
	"github.com/kzkei/natureAdvice/api/internal/services"

	"github.com/gin-gonic/gin"
)

func SetupRouter(db *sql.DB) *gin.Engine {
	// init repositories
	forecastRepo := repositories.NewForecastRepository(db)
	locationRepo := repositories.NewLocationRepository(db)

	// init services
	scoringService := services.NewScoringService()
	recommendationService := services.NewRecommendationService(forecastRepo, locationRepo, scoringService)

	// init handlers
	handler := handlers.NewHandlers(recommendationService, locationRepo, forecastRepo)

	// setup Gin router
	router := gin.Default() // Gin includes logger and recovery middleware

	// routes
	router.GET("/health", handler.HealthCheck)

	api := router.Group("/api")
	{
		api.POST("/locations", handler.CreateLocation)
		api.GET("/locations", handler.GetLocations)
		api.GET("/locations/:name/forecasts", handler.GetLocationForecast)          // single location forecasts (14 days)
		api.GET("/locations/:name/latest/:date", handler.GetLocationLatestForecast) // single location latest forecast for specified date
		api.GET("/recommendations/:date", handler.GetRecommendations)               // main idea
	}

	return router
}

package handlers

import (
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/kzkei/natureAdvice/api/internal/repositories"

	"github.com/kzkei/natureAdvice/api/internal/services"

	"github.com/gin-gonic/gin"
)

type Handlers struct {
	recService   *services.RecommendationService
	locationRepo *repositories.LocationRepository
	forecastRepo *repositories.ForecastRepository
}

func NewHandlers(
	recService *services.RecommendationService,
	locationRepo *repositories.LocationRepository,
	forecastRepo *repositories.ForecastRepository,
) *Handlers {
	return &Handlers{
		recService:   recService,
		locationRepo: locationRepo,
		forecastRepo: forecastRepo,
	}
}

// Health check
func (h *Handlers) HealthCheck(c *gin.Context) {
	log.Println("HC hit")
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

// POST /api/locations
func (h *Handlers) CreateLocation(c *gin.Context) {

	log.Println("CreateLocation handler hit")
	var req struct {
		Name      string  `json:"name" binding:"required"`
		Latitude  float64 `json:"latitude" binding:"required"`
		Longitude float64 `json:"longitude" binding:"required"`
		Region    string  `json:"region"` // optional
		State     string  `json:"state"`  // optional
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	location, err := h.locationRepo.CreateLocation(req.Name, req.Region, req.State, req.Latitude, req.Longitude)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// log.Printf("created location with ID %d", location.ID)

	c.JSON(http.StatusOK, location)
}

// GET /api/locations
func (h *Handlers) GetLocations(c *gin.Context) {

	log.Println("GetLocations hit")
	locations, err := h.locationRepo.GetLocations()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, locations)
}

// GET /api/locations/:name/forecast
func (h *Handlers) GetLocationForecast(c *gin.Context) {

	log.Println("GetLocationForecast hit")
	nameStr := c.Param("name")

	// check if location exists
	exists, err := h.locationRepo.LocationExistsByName(nameStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid location name"})
		return
	}
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "Location not found"})
		return
	}

	// get forecast for location by name

	forecast, err := h.forecastRepo.GetLocationForecastByName(nameStr)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, forecast)
}

// GET /api/locations/:name/latest/:date
func (h *Handlers) GetLocationLatestForecast(c *gin.Context) {

	log.Println("GetLocationLatestForecast hit")
	nameStr := c.Param("name")

	// checlk if location exists
	exists, err := h.locationRepo.LocationExistsByName(nameStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid location name"})
		return
	}
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "Location not found"})
		return
	}

	// parse date
	dateStr := c.Param("date")
	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid date format (use YYYY-MM-DD)"})
		return
	}

	// get latest location forecast for date
	forecast, err := h.forecastRepo.GetLatestForecastForDate(nameStr, date)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, forecast)
}

// GET /api/recommendations/:date?limit - optional param limit for top N otherwise defaults to top 10
func (h *Handlers) GetRecommendations(c *gin.Context) {

	log.Println("GetRecommendations hit")

	// parse date
	dateStr := c.Param("date")
	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid date format (use YYYY-MM-DD)"})
		return
	}

	// parse limit
	limit, err := strconv.Atoi(c.DefaultQuery("limit", "10"))
	if err != nil || limit < 1 {
		limit = 10 // default to 10
	}

	// get location recommendations for a certain date
	recommendations, err := h.recService.GetLocationRecommendationsForDate(date, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, recommendations)
}

// GET /api/recommendations/:name - to get reccs for a specific location with the full 14-day forecast and location forecast history in mind

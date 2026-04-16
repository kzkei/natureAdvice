package handlers

import (
	"time"

	"github.com/kzkei/natureAdvice/api/internal/models"
)

type locationRepo interface {
	CreateLocation(name, region, state string, latitude, longitude float64) (int64, error)
	GetLocations() ([]*models.Location, error)
	LocationExistsByName(name string) (bool, error)
	GetLocationByName(name string) (*models.Location, error)
	GetLocationNameByID(id int) (string, error)
	GetIDByName(name string) (int, error)
}

type forecastRepo interface {
	GetLocationForecastByName(name string) ([]*models.Forecast, error)
	GetLatestForecastForDate(nameStr string, date time.Time) (*models.Forecast, error)
	GetLocationForecastByID(location_id int) ([]*models.Forecast, error)
	GetForecastsByDate(date time.Time) ([]*models.Forecast, error)
}

type recService interface {
	GetLocationRecommendationsForDate(date time.Time, limit int) ([]models.Recommendation, error)
}

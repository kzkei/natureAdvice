package repositories

import (
	"database/sql"
	"time"

	"github.com/kzkei/natureAdvice/api/internal/models"
)

type ForecastRepository struct {
	db *sql.DB
}

func NewForecastRepository(db *sql.DB) *ForecastRepository {
	return &ForecastRepository{db: db}
}

// Get latest loc forecast for a specific date
func (r *ForecastRepository) GetLatestForecastForDate(location_id int, date time.Time) (*models.Forecast, error) {
	var latest_forecast models.Forecast

	err := r.db.QueryRow(`
        SELECT location_id, forecast_date, temp_high_f, temp_low_f,
               precipitation_chance, wind_speed_mph, uv_index, fetched_at
        FROM latest_forecasts
        WHERE location_id = $1 AND forecast_date = $2
    `, location_id, date).Scan(
		&latest_forecast.LocationID, &latest_forecast.Date, &latest_forecast.TempHigh, &latest_forecast.TempLow,
		&latest_forecast.PrecipChance, &latest_forecast.WindSpeed, &latest_forecast.UVIndex, &latest_forecast.FetchedAt,
	)

	if err != nil {
		return nil, err
	}

	return &latest_forecast, nil
}

// Get 14 day forecast for a location
func (r *ForecastRepository) GetLocationForecastByName(location_name string) ([]*models.Forecast, error) {
	var forecasts []*models.Forecast

	rows, err := r.db.Query(`
        SELECT location_id, forecast_date, temp_high_f, temp_low_f,
               precipitation_chance, wind_speed_mph, uv_index, fetched_at
        FROM weather_forecasts
        WHERE location_name = $1
        ORDER BY forecast_date ASC
    `, location_name)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var f models.Forecast
		err := rows.Scan(
			&f.LocationID, &f.Date, &f.TempHigh, &f.TempLow,
			&f.PrecipChance, &f.WindSpeed, &f.UVIndex, &f.FetchedAt,
		)
		if err != nil {
			return nil, err
		}
		forecasts = append(forecasts, &f)
	}

	return forecasts, nil
}

// Get Location forecasts by id
func (r *ForecastRepository) GetLocationForecastByID(location_id int) ([]*models.Forecast, error) {
	var forecasts []*models.Forecast

	rows, err := r.db.Query(`
		SELECT location_id, forecast_date, temp_high_f, temp_low_f,
			   precipitation_chance, wind_speed_mph, uv_index, fetched_at
		FROM weather_forecasts
		WHERE location_id = $1
		ORDER BY forecast_date ASC
	`, location_id)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var f models.Forecast
		err := rows.Scan(
			&f.LocationID, &f.Date, &f.TempHigh, &f.TempLow,
			&f.PrecipChance, &f.WindSpeed, &f.UVIndex, &f.FetchedAt,
		)
		if err != nil {
			return nil, err
		}
		forecasts = append(forecasts, &f)
	}

	return forecasts, nil

}

// Funcs for scoring/rec services

// Get forecasts for all locations for a specific date
func (r *ForecastRepository) GetForecastsByDate(date time.Time) ([]*models.Forecast, error) {
	var forecasts []*models.Forecast

	rows, err := r.db.Query(`
        SELECT location_id, forecast_date, temp_high_f, temp_low_f,
               precipitation_chance, wind_speed_mph, uv_index, fetched_at
        FROM weather_forecasts
        WHERE forecast_date = $1
        ORDER BY location_id ASC
    `, date)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var f models.Forecast
		err := rows.Scan(
			&f.LocationID, &f.Date, &f.TempHigh, &f.TempLow,
			&f.PrecipChance, &f.WindSpeed, &f.UVIndex, &f.FetchedAt,
		)
		if err != nil {
			return nil, err
		}
		forecasts = append(forecasts, &f)
	}

	return forecasts, nil
}

package models

import "time"

type Forecast struct {
	LocationID int `json:"location_id"`
	// LocationName string    `json:"location_name"` for future use - easier access after migration
	Date         time.Time `json:"forecast_date"`
	TempHigh     float64   `json:"temp_high_f"`
	TempLow      float64   `json:"temp_low_f"`
	PrecipChance int       `json:"precipitation_chance"`
	WindSpeed    float64   `json:"wind_speed_mph"`
	UVIndex      int       `json:"uv_index"`
	FetchedAt    time.Time `json:"fetched_at"`
}

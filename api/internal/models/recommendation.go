package models

// Recommendation for a location on a specific date with a score
type Recommendation struct {
	LocationID   int     `json:"location_id"`
	LocationName string  `json:"location_name"`
	Region       string  `json:"region"`
	Date         string  `json:"date"`
	Score        float64 `json:"score"`
	// Confidence   float64   `json:"confidence"` // for future use - how confident we are in the score based on volatility
	// Volatility   float64   `json:"volatility"` // for future use - how much the score is expected to change based on forecast volatility measured by history
	TempHigh     float64 `json:"temp_high_f"`
	TempLow      float64 `json:"temp_low_f"`
	PrecipChance int     `json:"precip_chance"`
	WindSpeed    float64 `json:"wind_speed_mph"`
	UVIndex      int     `json:"uv_index"`
	// ScoreBreakdown *ScoreBreakdown `json:"score_breakdown,omitempty"`// for future use - added breakdown for transparency
}

// ScoreBreakdown for transparency maybe
type ScoreBreakdown struct {
	TempHighScore float64 `json:"temp_high_score"`
	TempLowScore  float64 `json:"temp_low_score"`
	PrecipScore   float64 `json:"precip_score"`
	WindScore     float64 `json:"wind_score"`
	UVScore       float64 `json:"uv_score"`
}

package services

import (
	"log"

	"github.com/kzkei/natureAdvice/api/internal/models"
)

type ScoringService struct {
}

func NewScoringService() *ScoringService {
	return &ScoringService{}
}

// CalculateScore calculates a score for a given forecast based on forecast factors
func (s *ScoringService) CalculateScore(forecast *models.Forecast) float64 {

	if forecast == nil {
		log.Printf("nil forecast passed to CalculateScore")
		return 0.0
	}

	// simple scoring out of 30 total possible points for basic ideal conditions
	score := 0.0

	// forecast factors
	if forecast.TempHigh >= 50 && forecast.TempHigh <= 80 {
		score += 10.0 // ideal temperature range
	}
	if forecast.TempLow >= 40 && forecast.TempLow <= 70 {
		score += 10.0 // ideal low temperature range
	}
	if forecast.PrecipChance == 0 {
		score += 5.0
	}
	if forecast.WindSpeed < 10 {
		score += 3.0 // low wind is preferable
	}
	if forecast.UVIndex < 5 {
		score += 2.0 // moderate UV index
	}

	return score
}

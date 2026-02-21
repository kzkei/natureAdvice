package services

import (
	"log"
	"time"

	"github.com/kzkei/natureAdvice/api/internal/models"
	"github.com/kzkei/natureAdvice/api/internal/repositories"
)

type RecommendationService struct {
	forecastRepo   *repositories.ForecastRepository
	locationRepo   *repositories.LocationRepository
	scoringService *ScoringService
}

func NewRecommendationService(
	forecastRepo *repositories.ForecastRepository,
	locationRepo *repositories.LocationRepository,
	scoringService *ScoringService,
) *RecommendationService {
	return &RecommendationService{
		forecastRepo:   forecastRepo,
		locationRepo:   locationRepo,
		scoringService: scoringService,
	}
}

// GetLocationRecommendationsForDate returns the top N location recommendations for a given date
func (s *RecommendationService) GetLocationRecommendationsForDate(date time.Time, limit int) ([]models.Recommendation, error) {

	log.Printf("getting location recs for date %s with limit %d", date.Format("2006-01-02"), limit)
	// get forecasts for all locations for the date
	forecasts, err := s.forecastRepo.GetForecastsByDate(date)
	if err != nil {
		return nil, err
	}

	log.Printf("found %d forecasts for date %s", len(forecasts), date.Format("2006-01-02"))

	// score each forecast and create recommendations
	var recommendations []models.Recommendation

	for _, forecast := range forecasts {
		log.Printf("scoring location %d on date %s", forecast.LocationID, date.Format("2006-01-02"))
		// score
		score := s.scoringService.CalculateScore(forecast)

		// create recommendation
		recommendations = append(recommendations, models.Recommendation{
			LocationID:   forecast.LocationID,
			LocationName: forecast.LocationName,
			Region:       forecast.Region,
			Date:         date.Format("2006-01-02"),
			Score:        score,
			TempHigh:     forecast.TempHigh,
			TempLow:      forecast.TempLow,
			PrecipChance: forecast.PrecipChance,
			WindSpeed:    forecast.WindSpeed,
			UVIndex:      forecast.UVIndex,
			// ScoreBreakdown: &models.ScoreBreakdown{}, // for future use - added breakdown for transparency
		})
	}

	// sort recommendations by score and return top N
	recommendations = SortAndLimitRecommendations(recommendations, limit)

	log.Printf("returning %d recs for date %s", len(recommendations), date.Format("2006-01-02"))

	return recommendations, nil
}

func SortAndLimitRecommendations(recommendations []models.Recommendation, limit int) []models.Recommendation {
	var sorted []models.Recommendation

	log.Printf("sorting recs and finding top %d", limit)

	// sort recommendations by score in descending order
	for _, rec := range recommendations {
		inserted := false
		for j, sortedRec := range sorted {
			if rec.Score > sortedRec.Score {
				sorted = append(sorted[:j], append([]models.Recommendation{rec}, sorted[j:]...)...)
				inserted = true
				break
			}
		}
		if !inserted {
			sorted = append(sorted, rec)
		}
	}

	log.Printf("sorted recs, total %d", len(sorted))

	// limit to the top N recommendations
	if len(sorted) > limit {
		sorted = sorted[:limit]
	}

	log.Printf("limited recs to top %d, returning", len(sorted))

	return sorted
}

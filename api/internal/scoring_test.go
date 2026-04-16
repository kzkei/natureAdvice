package internal

import (
	"testing"

	"github.com/go-playground/assert/v2"
	"github.com/kzkei/natureAdvice/api/internal/models"
	"github.com/kzkei/natureAdvice/api/internal/services"
)

func TestHighScoreCalculation(t *testing.T) {
	// create a mock forecast with known values
	forecast := models.Forecast{
		TempHigh:     80, // +10
		TempLow:      60, // +10
		PrecipChance: 20, // +0
		WindSpeed:    10, // +0
		UVIndex:      5,  // +0
	}

	scoringService := services.NewScoringService()

	score := scoringService.CalculateScore(&forecast)

	expectedScore := 20.0

	if score != expectedScore {
		t.Errorf("Expected score %f, got %f", expectedScore, score)
	}

	assert.Equal(t, expectedScore, score)
}

func TestLowScoreCalculation(t *testing.T) {
	// create a mock forecast with known values
	forecast := models.Forecast{
		TempHigh:     50, // +10
		TempLow:      30, // +0
		PrecipChance: 80, // +0
		WindSpeed:    7,  // +3
		UVIndex:      3,  // +2
	}

	scoringService := services.NewScoringService()

	score := scoringService.CalculateScore(&forecast)

	expectedScore := 15.0

	if score != expectedScore {
		t.Errorf("Expected score %f, got %f", expectedScore, score)
	}

	assert.Equal(t, expectedScore, score)
}

func TestNilForecast(t *testing.T) {
	scoringService := services.NewScoringService()

	score := scoringService.CalculateScore(nil)

	expectedScore := 0.0

	if score != expectedScore {
		t.Errorf("Expected score %f, got %f", expectedScore, score)
	}

	assert.Equal(t, expectedScore, score)
}

func TestExtremeConditions(t *testing.T) {
	// create a mock forecast with extreme values
	forecast := models.Forecast{
		TempHigh:     100, // +0
		TempLow:      20,  // +0
		PrecipChance: 100, // +0
		WindSpeed:    20,  // +0
		UVIndex:      10,  // +0
	}

	scoringService := services.NewScoringService()

	score := scoringService.CalculateScore(&forecast)

	expectedScore := 0.0

	if score != expectedScore {
		t.Errorf("Expected score %f, got %f", expectedScore, score)
	}

	assert.Equal(t, expectedScore, score)
}

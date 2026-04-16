package handlers_test

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/kzkei/natureAdvice/api/internal/handlers"
	"github.com/kzkei/natureAdvice/api/internal/repositories"
	"github.com/kzkei/natureAdvice/api/internal/services"
	"github.com/kzkei/natureAdvice/api/internal/testhelpers"
)

// setupIntegrationRouter wires real dependencies against the test DB
func setupIntegrationRouter(db *sql.DB) *gin.Engine {
	gin.SetMode(gin.TestMode)

	locationRepo := repositories.NewLocationRepository(db)
	forecastRepo := repositories.NewForecastRepository(db)
	scoringService := services.NewScoringService()
	recService := services.NewRecommendationService(forecastRepo, locationRepo, scoringService)

	h := handlers.NewHandlers(recService, locationRepo, forecastRepo)

	router := gin.New()
	router.GET("/api/locations", h.GetLocations)
	router.GET("/api/locations/:name/forecasts", h.GetLocationForecast)
	router.GET("/api/locations/:name/latest/:date", h.GetLocationLatestForecast)
	router.GET("/api/recommendations/:date", h.GetRecommendations)

	return router
}

// GetLocations tests

func TestGetLocations_ReturnsSeededLocations(t *testing.T) {
	db, cleanup := testhelpers.SetupTestDB(t)
	defer cleanup()

	testhelpers.SeedLocation(t, db, "Yosemite National Park", "YosemiteNP")
	testhelpers.SeedLocation(t, db, "Yellowstone National Park", "YNP")

	router := setupIntegrationRouter(db)

	req := httptest.NewRequest("GET", "/api/locations", nil)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}

	var locations []map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &locations); err != nil {
		t.Fatalf("failed to parse response body: %v", err)
	}

	if len(locations) != 2 {
		t.Errorf("expected 2 locations, got %d", len(locations))
	}
}

func TestGetLocations_EmptyDB_ReturnsEmptyOrOK(t *testing.T) {
	db, cleanup := testhelpers.SetupTestDB(t)
	defer cleanup()

	// no seeding — empty db
	router := setupIntegrationRouter(db)

	req := httptest.NewRequest("GET", "/api/locations", nil)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200 on empty db, got %d", w.Code)
	}
}

// GetLocationForecast tests

func TestGetLocationForecast_ValidLocation_Returns200(t *testing.T) {
	db, cleanup := testhelpers.SetupTestDB(t)
	defer cleanup()

	id := testhelpers.SeedLocation(t, db, "Yosemite National Park", "YosemiteNP")
	testhelpers.SeedWeatherForecast(t, db, id, "Yosemite National Park", "West", time.Now().Truncate(24*time.Hour))

	router := setupIntegrationRouter(db)

	name := url.PathEscape("Yosemite National Park")
	req := httptest.NewRequest("GET", "/api/locations/"+name+"/forecasts", nil)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d, body: %s", w.Code, w.Body.String())
	}

}

func TestGetLocationForecast_UnknownLocation_Returns404(t *testing.T) {
	db, cleanup := testhelpers.SetupTestDB(t)
	defer cleanup()

	router := setupIntegrationRouter(db)

	name := url.PathEscape("Nonexistent Park")
	req := httptest.NewRequest("GET", "/api/locations/"+name+"/forecasts", nil)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("expected 404, got %d", w.Code)
	}
}

// GetLocationLatestForecast tests

func TestGetLocationLatestForecast_ValidLocationAndDate_Returns200(t *testing.T) {
	db, cleanup := testhelpers.SetupTestDB(t)
	defer cleanup()

	date := time.Now().Truncate(24 * time.Hour)

	id := testhelpers.SeedLocation(t, db, "Yosemite National Park", "YosemiteNP")
	testhelpers.SeedLatestForecast(t, db, id, "Yosemite National Park", "West", date)

	router := setupIntegrationRouter(db)

	dateStr := date.Format("2006-01-02")
	name := url.PathEscape("Yosemite National Park")
	req := httptest.NewRequest("GET", "/api/locations/"+name+"/latest/"+dateStr, nil)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d, body: %s", w.Code, w.Body.String())
	}
}

func TestGetLocationLatestForecast_BadDateFormat_Returns400(t *testing.T) {
	db, cleanup := testhelpers.SetupTestDB(t)
	defer cleanup()

	testhelpers.SeedLocation(t, db, "Yosemite National Park", "YosemiteNP")

	router := setupIntegrationRouter(db)

	name := url.PathEscape("Yosemite National Park")
	req := httptest.NewRequest("GET", "/api/locations/"+name+"/latest/not-a-date", nil)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
}

func TestGetLocationLatestForecast_UnknownLocation_Returns404(t *testing.T) {
	db, cleanup := testhelpers.SetupTestDB(t)
	defer cleanup()

	router := setupIntegrationRouter(db)

	name := url.PathEscape("Ghost Park")
	req := httptest.NewRequest("GET", "/api/locations/"+name+"/latest/2026-06-01", nil)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("expected 404, got %d", w.Code)
	}
}

// GetRecs tests

func TestGetRecommendations_ValidDate_Returns200(t *testing.T) {
	db, cleanup := testhelpers.SetupTestDB(t)
	defer cleanup()

	date := time.Now().Truncate(24 * time.Hour)
	id := testhelpers.SeedLocation(t, db, "Yosemite National Park", "YosemiteNP")
	testhelpers.SeedWeatherForecast(t, db, id, "Yosemite National Park", "West", date)

	router := setupIntegrationRouter(db)

	dateStr := date.Format("2006-01-02")
	req := httptest.NewRequest("GET", "/api/recommendations/"+dateStr, nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d, body: %s", w.Code, w.Body.String())
	}
}

func TestGetRecommendations_BadDateFormat_Returns400(t *testing.T) {
	db, cleanup := testhelpers.SetupTestDB(t)
	defer cleanup()

	router := setupIntegrationRouter(db)

	req := httptest.NewRequest("GET", "/api/recommendations/not-a-date", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
}

func TestGetRecommendations_LimitParam_RespectsLimit(t *testing.T) {
	db, cleanup := testhelpers.SetupTestDB(t)
	defer cleanup()

	date := time.Now().Truncate(24 * time.Hour)

	// seed 3 locations with forecasts
	for _, park := range []struct{ name, code string }{
		{"Yosemite National Park", "YosemiteNP"},
		{"Yellowstone National Park", "YNP"},
		{"Zion National Park", "ZNP"},
	} {
		id := testhelpers.SeedLocation(t, db, park.name, park.code)
		testhelpers.SeedWeatherForecast(t, db, id, park.name, "West", date)
	}

	router := setupIntegrationRouter(db)

	dateStr := date.Format("2006-01-02")
	req := httptest.NewRequest("GET", "/api/recommendations/"+dateStr+"?limit=2", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}

	var recs []map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &recs); err != nil {
		t.Fatalf("failed to parse response: %v", err)
	}

	if len(recs) != 2 {
		t.Errorf("expected 2 recommendations with limit=2, got %d", len(recs))
	}
}

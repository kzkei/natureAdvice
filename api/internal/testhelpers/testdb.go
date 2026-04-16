package testhelpers

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	"time"

	"github.com/kzkei/natureAdvice/api/config"
	"github.com/kzkei/natureAdvice/api/database"
)

// TestConfig returns config pointed at test db container
func TestConfig() *config.Config {
	return &config.Config{
		DBHost:     "localhost",
		DBPort:     "5434",
		DBUser:     "natureadvice_test",
		DBPassword: "natureadvice_test",
		DBName:     "natureadvice_test",
		APIPort:    "8001",
	}
}

// SetupTestDB connects to test db and runs migrations, returns the db connection and a cleanup function to call after after each test
func SetupTestDB(t *testing.T) (*sql.DB, func()) {
	t.Helper()

	cfg := TestConfig()

	log.Printf("setting up test db with config: %+v", cfg)

	db, err := database.Connect(cfg)
	if err != nil {
		t.Fatalf("failed to connect to test database: %v", err)
	}

	// run migrations on test db
	runMigrations(t, db)

	// clean up db after tests
	cleanup := func() {
		tearDownTables(t, db)
		db.Close()
	}

	return db, cleanup
}

// runMigrations runs migration files (only files with 'up') on the test db via connection
func runMigrations(t *testing.T, db *sql.DB) {
	t.Helper()

	// walk up from this file and find migrations folder
	_, filename, _, _ := runtime.Caller(0)
	migrationsPath := filepath.Join(filepath.Dir(filename), "..", "..", "..", "db", "migrations")

	entries, err := os.ReadDir(migrationsPath)
	if err != nil {
		t.Fatalf("failed to read migrations directory: %v", err)
	}

	// run migrations in order, only 'up' and skip directories
	for _, entry := range entries {

		// skip directories
		if entry.IsDir() {
			continue
		}

		// only run 'up' migrations
		if !strings.HasSuffix(entry.Name(), "up.sql") {
			continue
		}

		migrationFile := filepath.Join(migrationsPath, entry.Name())
		content, err := os.ReadFile(migrationFile)
		if err != nil {
			t.Fatalf("failed to read migration file %s: %v", entry.Name(), err)
		}

		// run migration on db
		if _, err := db.Exec(string(content)); err != nil {
			// IF EXISTS on alters means re-running is safe, log and continue
			log.Printf("migration %s note: %v", entry.Name(), err)
		}

		log.Printf("ran migration: %s", entry.Name())
	}
}

// tearDownTables truncates all tables in test db to clean up after tests
func tearDownTables(t *testing.T, db *sql.DB) {
	t.Helper()

	// order matters, remove tables with fks first
	tables := []string{"latest_forecasts", "weather_forecasts", "locations"}
	for _, table := range tables {
		if _, err := db.Exec(fmt.Sprintf("TRUNCATE TABLE %s CASCADE", table)); err != nil {
			t.Fatalf("failed to truncate %s: %v", table, err)
		}
		log.Printf("truncated table: %s", table)
	}
}

// seed funcs

// SeedLocations seeds the location table with test data, returns location id
// dont care abut accuracy of data
func SeedLocation(t *testing.T, db *sql.DB, name, parkCode string) int64 {
	t.Helper()

	var id int64
	err := db.QueryRow(
		`INSERT INTO locations 
			(name, park_code, latitude, longitude, state, region, timezone, elevation_ft)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		 RETURNING id`,
		name, parkCode, 37.8651, -119.5383, "CA", "West", "America/Los_Angeles", 4000,
	).Scan(&id)

	if err != nil {
		t.Fatalf("failed to seed location %s: %v", name, err)
	}

	return id
}

// SeedWeatherForecast inserts a weather forecast for a location
// dont care about accuracy of data
func SeedWeatherForecast(t *testing.T, db *sql.DB, locationID int64, locationName, region string, date time.Time) {
	t.Helper()

	_, err := db.Exec(
		`INSERT INTO weather_forecasts
			(location_id, forecast_date, temp_high_f, temp_low_f,
			 precipitation_chance, wind_speed_mph, uv_index,
			 location_name, region)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`,
		locationID, date, 75.0, 55.0, 10, 5.0, 3, locationName, region,
	)

	if err != nil {
		t.Fatalf("failed to seed weather forecast for location %d: %v", locationID, err)
	}
}

// SeedLatestForecast inserts a latest_forecast entry for a location
// dont care about accuracy
func SeedLatestForecast(t *testing.T, db *sql.DB, locationID int64, locationName, region string, date time.Time) {
	t.Helper()

	_, err := db.Exec(
		`INSERT INTO latest_forecasts
			(location_id, forecast_date, temp_high_f, temp_low_f,
			 precipitation_chance, wind_speed_mph, uv_index,
			 location_name, region)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`,
		locationID, date, 75.0, 55.0, 10, 5.0, 3, locationName, region,
	)

	if err != nil {
		t.Fatalf("failed to seed latest forecast for location %d: %v", locationID, err)
	}
}

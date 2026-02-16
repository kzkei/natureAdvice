package repositories

import (
	"database/sql"

	"github.com/kzkei/natureAdvice/api/internal/models"
)

type LocationRepository struct {
	db *sql.DB
}

func NewLocationRepository(db *sql.DB) *LocationRepository {
	return &LocationRepository{db: db}
}

// CreateLocation creates a new location entity in the database
func (r *LocationRepository) CreateLocation(name string, latitude, longitude float64) (int64, error) {
	var id int64
	err := r.db.QueryRow(`
		INSERT INTO locations (name, latitude, longitude)
		VALUES ($1, $2, $3)
		RETURNING id
	`, name, latitude, longitude).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

// GetLocations retrieves all location entites from the database
func (r *LocationRepository) GetLocations() ([]*models.Location, error) {
	rows, err := r.db.Query(`SELECT id, name, latitude, longitude FROM locations`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var locations []*models.Location
	for rows.Next() {
		var loc models.Location
		if err := rows.Scan(&loc.ID, &loc.Name, &loc.Latitude, &loc.Longitude); err != nil {
			return nil, err
		}
		locations = append(locations, &loc)
	}
	return locations, nil
}

// LocationExistsByName returns true if a location with the given name exists
func (r *LocationRepository) LocationExistsByName(name string) (bool, error) {
	var exists bool
	err := r.db.QueryRow(`
		SELECT EXISTS(SELECT 1 FROM locations WHERE LOWER(name) = LOWER($1))
	`, name).Scan(&exists) // scan overwrites exists var (nil)
	if err != nil {
		return false, err
	}
	return exists, nil
}

// GetLocationByName retrieves a location entity by its name
func (r *LocationRepository) GetLocationByName(name string) (*models.Location, error) {
	var loc models.Location
	err := r.db.QueryRow(`
		SELECT id, name, latitude, longitude
		FROM locations
		WHERE LOWER(name) = LOWER($1)
	`, name).Scan(&loc.ID, &loc.Name, &loc.Latitude, &loc.Longitude)
	if err != nil {
		return nil, err
	}
	return &loc, nil
}

// GetLocationNameByID retrieves a location name by its ID
func (r *LocationRepository) GetLocationNameByID(id int) (string, error) {
	var name string
	err := r.db.QueryRow(`
		SELECT name
		FROM locations
		WHERE id = $1
	`, id).Scan(&name)
	if err != nil {
		return "", err
	}
	return name, nil
}

// change forecast table to include name and use name for queries instead of id
func (r *LocationRepository) GetIDByName(name string) (int, error) {
	var id int
	err := r.db.QueryRow(`
		SELECT id
		FROM locations
		WHERE LOWER(name) = LOWER($1)
	`, name).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

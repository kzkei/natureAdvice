CREATE TABLE locations (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    park_code VARCHAR(10) UNIQUE,
    latitude DOUBLE PRECISION NOT NULL,
    longitude DOUBLE PRECISION NOT NULL,
    state VARCHAR(2),
    region VARCHAR(50),
    elevation_ft INTEGER,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE weather_forecasts (
    id SERIAL PRIMARY KEY,
    location_id INTEGER NOT NULL REFERENCES locations(id),
    forecast_date DATE NOT NULL,
    temp_high_f DECIMAL(5,2),
    temp_low_f DECIMAL(5,2),
    precipitation_chance INTEGER,
    wind_speed_mph DOUBLE PRECISION,
    uv_index INTEGER,
    fetched_at TIMESTAMP DEFAULT NOW(),
    UNIQUE(location_id, forecast_date)
);

CREATE INDEX idx_location_date ON weather_forecasts(location_id, forecast_date);
CREATE INDEX idx_locations_park_code ON locations(park_code);
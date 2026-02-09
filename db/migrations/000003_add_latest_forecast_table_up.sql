CREATE TABLE latest_forecasts (
    location_id INTEGER NOT NULL,
    forecast_date DATE NOT NULL,
    temp_high_f DECIMAL(5,2),
    temp_low_f DECIMAL(5,2),
    precipitation_chance INTEGER,
    wind_speed_mph DOUBLE PRECISION,
    uv_index INTEGER,
    fetched_at TIMESTAMP DEFAULT NOW(),
    PRIMARY KEY (location_id, forecast_date),
    FOREIGN KEY (location_id) REFERENCES locations(id)
);

CREATE INDEX idx_latest_forecasts_date 
    ON latest_forecasts(forecast_date);
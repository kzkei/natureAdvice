-- add loc_name and region cols to forecasts for easier access without extra querying
ALTER TABLE weather_forecasts ADD COLUMN IF NOT EXISTS location_name VARCHAR(50);
ALTER TABLE weather_forecasts ADD COLUMN IF NOT EXISTS region VARCHAR(50);

ALTER TABLE latest_forecasts ADD COLUMN IF NOT EXISTS location_name VARCHAR(50);

ALTER TABLE latest_forecasts ADD COLUMN IF NOT EXISTS region VARCHAR(50);
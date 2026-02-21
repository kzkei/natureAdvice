ALTER TABLE weather_forecasts DROP COLUMN IF EXISTS location_name;
ALTER TABLE weather_forecasts DROP COLUMN IF EXISTS region;
ALTER TABLE latest_forecasts DROP COLUMN IF EXISTS location_name;
ALTER TABLE latest_forecasts DROP COLUMN IF EXISTS region;
from extract.locations import extract_location_data
from extract.weather import extract_weather_data
from transform.normalize import transform_weather
from load.load import load_forecasts
import logging

### This main is now void as the ETL process is now orchestrated in a DAG within Airflow

# logging setup for main ETL entry point, can use in entry points elsewhere
from utils.log_config import setup_logging
setup_logging()

logger = logging.getLogger(__name__)

# main orchestrates the ETL flow and acts as the main ETL entry point
# separates extraction flow from processing flow in a two step process
def main():

    logger.info("Starting ETL main")

    # 1 - extract
    locations = extract_location_data()
    logger.debug(f"extracting and transforming weather data for {len(locations)} locations")

    # store successfully extracted forecast data
    raw_forecasts = []
    for location in locations:
        try:
            weather_raw = extract_weather_data(
            location['latitude'],
            location['longitude'],
            location['timezone']
            )

            raw_forecasts.append(
            {'location_id': location['id'],
            'weather_raw': weather_raw})

        except Exception as e:
            logger.error(f"failed to fetch weather for {location['name']}: {e}")
            # continue on failed API calls to avoid pipeline crash
        
    # 2 - process extracted data (successful extractions only)
    for entry in raw_forecasts:
        try:
            # transform raw location weather data
            transformed = transform_weather(entry['weather_raw'])
            
            # load normalized forecast
            load_forecasts(location_id=entry['location_id'], forecasts=transformed)

        except Exception as e:
            logger.error(f"failed to transform/load entry {entry['location_id']}: {e}")

    logger.info("ETL complete")

if __name__ == "__main__":
    main()
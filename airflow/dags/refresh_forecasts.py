from airflow.decorators import dag, task
from datetime import timedelta
import logging

from etl.extract.locations import extract_location_data
from etl.extract.weather import extract_weather_data
from etl.transform.normalize import transform_weather
from etl.load.load import load_forecasts

# logging setup for ETL entry point
from etl.utils.log_config import setup_logging
setup_logging()

logger = logging.getLogger(__name__)

default_args = {
    'owner': 'airflow',
    'depends_on_past': False,
    'email_on_failure': False,
    'email_on_retry': False,
    'retries': 2,
    'retry_delay': timedelta(minutes=5),
}

@dag(default_args=default_args, schedule='0 */6 * * *', catchup=False) # fetch every 6 hours, likely to change
def forecast_pipeline():

    # separate extraction and normalize/load steps

    @task
    def get_locations():
        return extract_location_data()
    
    @task
    def get_forecasts(locations):
        raw_forecasts = []
        for loc in locations:
            try:
                weather = extract_weather_data(
                    loc['latitude'],
                    loc['longitude'],
                    loc['timezone'])
                
                raw_forecasts.append(
                    {'location_id': loc['id'],
                    'weather_raw': weather})

            except Exception as e:
                logger.error(f"failed to fetch weather for {loc['name']}: {e}")

        return raw_forecasts

    
    @task
    def process_locations(raw_forecasts):
        results = []
        for entry in raw_forecasts:
            try:
                # transform raw location weather data
                transformed = transform_weather(entry['weather_raw'])

                # load normalized forecasts into forecasts all time and latest tables (via load.py)
                load_forecasts(location_id=entry['location_id'], forecasts=transformed)
                results.append(entry['location_id'])

            except Exception as e:
                logger.error(f"failed to transform/load entry {entry['location_id']}: {e}")

        return len(results)
    
    # define dag flow of sequential steps
    locations = get_locations()
    forecasts = get_forecasts(locations=locations)
    process_locations(raw_forecasts=forecasts)

forecast_pipeline()
import requests
from utils import conn_ops
import logging

# setup logging
logger = logging.getLogger(__name__)

# EXTRACT location weather from OpenMeteo
# returns raw weather data json
def extract_weather_data(lat, lon, timezone):
    logger.info("Starting api weather extraction")


    url = "https://api.open-meteo.com/v1/forecast"

    params = {
        'forecast_days': 14,
        'latitude': lat,
        'longitude': lon,
        'daily': ",".join([
            'temperature_2m_max',
            'temperature_2m_min',
            'precipitation_probability_max',
            'wind_speed_10m_max',
            'uv_index_max'
        ]),
        'temperature_unit': 'fahrenheit',
        'wind_speed_unit': 'mph',
        'timezone': timezone,
    }

    logger.debug(f"lat and lon are {lat} and {lon}")

    response = requests.get(url, params=params, timeout=15)
    response.raise_for_status()

    logger.info("weather extraction complete")

    return response.json()

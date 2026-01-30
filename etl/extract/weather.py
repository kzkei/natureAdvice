import requests
from utils import conn_ops

# EXTRACT location weather
def extract_weather_data(lat, lon, timezone):
    # TODO logging, error handling

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

    response = requests.get(url, params=params, timeout=15)
    response.raise_for_status()

    return response.json()

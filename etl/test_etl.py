import requests
import psycopg2
from datetime import datetime

# Db config to connect
DB_CONFIG = {
    'host': 'localhost',
    'port': 5433,
    'user': 'natureadvice',
    'password': 'natureadvice23',
    'database': 'natureadvice'
}

# EXTRACT locations from locations table
def extract_location_metadata():
    # pyscopg2 is used for postgres adapted connection
    conn = psycopg2.connect(**DB_CONFIG)
    cursor = conn.cursor()

    # execute fetch SQL into locations
    cursor.execute("SELECT * FROM locations;")
    locations = cursor.fetchall()

    cursor.close()
    conn.close()

# return converted tuples in list of dict
    return [
        {
            'id': row[0],
            'name': row[1],
            'park_code': row[2],
            'latitude': row[3],
            'longitude': row[4],
            'state': row[5],
            'region': row[6],
            'elevation_ft': row[7],
            'timezone': row[10]
        }
        for row in locations
    ]

# EXTRACT location weather
def extract_weather_metadata(lat, lon, tz):

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
        'timezone': tz,
    }

    response = requests.get(url, params=params)
    # print(response.url)
    # print(response.text)
    response.raise_for_status()

    # print(response)

    return response.json()

# transform weather data into schema
def transform_weather(raw):
    daily = raw['daily']

    return

def main():

    return

if __name__ == "__main__":
    main()
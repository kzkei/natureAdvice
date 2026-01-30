import requests
import psycopg2
from utils import conn_ops
from datetime import datetime

# EXTRACT locations from locations table
def extract_location_metadata():
    # TODO logging, error handling

    # connect
    conn, cursor = conn_ops.open()

    # execute fetch SQL into locations
    cursor.execute("SELECT * FROM locations;")
    locations = cursor.fetchall()

    conn_ops.close(conn=conn, cursor=cursor)

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
    response.raise_for_status()

    return response.json()

# transform weather data into defined schema
def transform_weather(raw):
    daily = raw['daily']

    # append each respective raw i data to forecasts list of dict
    forecasts = []
    for i in range(len(daily['time'])):
        forecasts.append({
            'forecast_date': daily['time'][i],
            'temp_high_f': daily['temperature_2m_max'][i],
            'temp_low_f': daily['temperature_2m_min'][i],
            'precipitation_chance': daily['precipitation_probability_max'][i],
            'wind_speed_mph': daily['wind_speed_10m_max'][i],
            'uv_index': daily['uv_index_max'][i]
        })

    # print(forecasts)
    return forecasts

# leave load for L file

# to test ETL flow
def main():

    print("starting test ETL")

    locations = extract_location_metadata()

    print(f"extracting and transforming weather data for {len(locations)} locations")

    forecasts = []
    for location in locations:
        weather_raw = extract_weather_metadata(
        location['latitude'],
        location['longitude'],
        location['timezone']
        )
        # transform every locations raw weather metadata
        transformed = transform_weather(weather_raw)
        
        # append transformed to forecasts list
        forecasts.append({
            'location_id': location['id'],
            'forecasts': transformed
        })

    print(forecasts)
    print("L not in test pipeline")

    return

if __name__ == "__main__":
    main()
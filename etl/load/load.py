import logging, psycopg2
from utils import conn_ops

# setup logging
logger = logging.getLogger(__name__)

# LOAD forecasts into db
def load_forecasts(location_id, forecasts):
    logger.info("Starting forecast load")

    conn, cursor = conn_ops.open()

    for forecast in forecasts:
        logger.debug(f"loading for location_id: {location_id}")
        cursor.execute("""
            INSERT INTO weather_forecasts 
            (location_id, forecast_date, temp_high_f, temp_low_f, 
             precipitation_chance, wind_speed_mph, uv_index)
            VALUES (%s, %s, %s, %s, %s, %s, %s)
            ON CONFLICT (location_id, forecast_date) 
            DO UPDATE SET
                temp_high_f = EXCLUDED.temp_high_f,
                temp_low_f = EXCLUDED.temp_low_f,
                precipitation_chance = EXCLUDED.precipitation_chance,
                wind_speed_mph = EXCLUDED.wind_speed_mph,
                uv_index = EXCLUDED.uv_index,
                fetched_at = NOW()
        """, (
            location_id,
            forecast['forecast_date'],
            forecast['temp_high_f'],
            forecast['temp_low_f'],
            forecast['precipitation_chance'],
            forecast['wind_speed_mph'],
            forecast['uv_index']
        ))

    conn.commit()
    conn_ops.close(conn=conn,cursor=cursor)

    logger.info("Loading forecasts completed")
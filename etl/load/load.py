import logging, psycopg2
from utils import conn_ops

# setup logging
logger = logging.getLogger(__name__)

# LOAD forecasts into db
def load_forecasts(location_id, forecasts):
    logger.info("Starting forecast load")

    conn, cursor = conn_ops.open()
    try:
        for forecast in forecasts:

            # append to weather_forecasts (for history)
            # on conflict updates fetched entry if task is retried in DAG
            logger.debug(f"loading weather_forecasts for location_id: {location_id}")
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

            # upsert latest forecast for each unique
            logger.debug(f"loading latest_forecast for location_id: {location_id}")
            cursor.execute("""
                INSERT INTO latest_forecasts 
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
        logger.info(f"loaded {len(forecasts)} forecasts for location {location_id}")

    except Exception as e:
        conn.rollback()
        logger.error(f"failed to load forecasts: {e}")
        raise

    finally:
        conn_ops.close(conn=conn,cursor=cursor)
        logger.info("Loading forecasts completed")
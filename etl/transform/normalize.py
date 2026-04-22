import logging

# setup logging
logger = logging.getLogger(__name__)

# normalize weather data for defined schema
# returns normalized weather data
def transform_weather(raw):
    logger.info("Starting raw weather transformation")

    # skip on missing critical keys for fault tolerance
    if not raw.get('daily') or not raw['daily'].get('time'):
        raise ValueError("invalid API response: missing critical keys in raw")

    daily = raw['daily']

    logger.debug(f"raw data passed: {raw}")

    # append each respective raw i data to forecasts list of dict
    forecasts = []
    for i in range(len(daily['time'])):

        # check each forecast exists, skip whole entry if not
        if daily['time'][i] is None:
            logger.warning(f"missing forecast date at index {i}, skipping")
            continue

        forecasts.append({
            'forecast_date': daily['time'][i],
            'temp_high_f': default_if_none(daily['temperature_2m_max'][i], 57.0),
            'temp_low_f': default_if_none(daily['temperature_2m_min'][i], 57.0),
            'precipitation_chance': default_if_none(daily['precipitation_probability_max'][i], 0),
            'wind_speed_mph': default_if_none(daily['wind_speed_10m_max'][i], 0),
            'uv_index': default_if_none(daily['uv_index_max'][i], 0)
        })

    logger.info("raw data normalizing complete")

    return forecasts

# helper to handle nulls
def default_if_none(value, default):
    if value is not None:
        return value
    else:
        return default
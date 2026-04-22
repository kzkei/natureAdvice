# pipeline transform/normalize tests with fixture

import pytest
import transform.normalize as transform
import datetime

@pytest.fixture
def raw_instance():
    """Return jsonified open-meteo API response"""
    return {
  "latitude": 52.52,
  "longitude": 13.419998,
  "generationtime_ms": 0.28693675994873047,
  "utc_offset_seconds": 7200,
  "timezone": "Europe/Berlin",
  "timezone_abbreviation": "GMT+2",
  "elevation": 38.0,
  "daily_units": {
    "time": "iso8601",
    "temperature_2m_max": "°F",
    "temperature_2m_min": "°F",
    "uv_index_max": "",
    "precipitation_probability_max": "%",
    "wind_speed_10m_max": "mp/h"
  },
  "daily": {
    "time": [
      "2026-04-17",
      "2026-04-18",
      "2026-04-19",
      "2026-04-20",
      "2026-04-21",
      "2026-04-22",
      "2026-04-23",
      "2026-04-24",
      "2026-04-25",
      "2026-04-26",
      "2026-04-27",
      "2026-04-28",
      "2026-04-29",
      "2026-04-30"
    ],
    "temperature_2m_max": [67.5, 68.8, 60.3, 49.2, 58.1, 62.2, 56.2, 56.2, 60.8, 62.6, 55.1, 53.1, 55.0, 51.5],
    "temperature_2m_min": [47.7, 49.9, 45.5, 40.5, 35.2, 36.6, 40.1, 35.7, 40.6, 48.1, 48.3, 44.6, 38.5, 38.5],
    "uv_index_max": [5.35, 5.35, 4.95, 3.2, 5.4, 4.25, 4.7, 5.45, 4.9, 3.2, 3.6, 3.9, 5.45, 4.3],
    "precipitation_probability_max": [0, 10, 68, 51, 33, 6, 13, 11, 23, 35, 26, 16, 20, 16],
    "wind_speed_10m_max": [7.3, 4.9, 9.8, 10.4, 6.5, 12.2, 13.7, 6.7, 6.2, 9.1, 13.8, 19.2, 12.8, 10.9]
  }
}

def test_happypath(raw_instance):
    """Test normalizing data success - happy path"""

    # trasnform data in fixture reference
    transformed = transform.transform_weather(raw_instance)

    # assert satisfaction
    # count of normalized entries (forecasts)
    assert len(transformed) == 14

    # types
    assert type(transformed[0]['forecast_date']) is str # given as str, is date in table
    assert type(transformed[0]['temp_high_f']) == float
    assert type(transformed[0]['temp_low_f']) == float
    assert type(transformed[0]['precipitation_chance']) == int
    assert type(transformed[0]['wind_speed_mph']) == float
    assert type(transformed[0]['uv_index']) == float # persists as float, SQL insertion rounds to nearest int

    # value persistence
    assert transformed[0]['temp_high_f'] == 67.5

def test_null_time(raw_instance):
    """Test null time field is rejected"""

    raw_instance['daily']['time'] = None

    with pytest.raises(ValueError):
        transform.transform_weather(raw_instance)

def test_null_single_forecast(raw_instance):
    """Test null time for one entry is skipped"""

    raw_instance['daily']['time'][0] = None

    transformed = transform.transform_weather(raw_instance)
    assert len(transformed) == 13 # one missing forecast date gets skipped
    

def test_null_temp(raw_instance):
    """Test null temp fields default correctly"""

    raw_instance['daily']['temperature_2m_max'][0] = None

    transformed = transform.transform_weather(raw_instance)
    assert transformed[0]['temp_high_f'] == 57.0

def test_null_uv(raw_instance):
    """Test null uv field defaults correctly"""

    raw_instance['daily']['uv_index_max'][0] = None

    transformed = transform.transform_weather(raw_instance)
    assert transformed[0]['uv_index'] == 0.0 # is float until LOAD occurs

def test_zero_values(raw_instance):
    """Test zeros produce successful transformation"""

    raw_instance['daily']['temperature_2m_max'][0] = 0.0
    raw_instance['daily']['precipitation_probability_max'][0] = 0
    raw_instance['daily']['wind_speed_10m_max'][0] = 0.0

    transformed = transform.transform_weather(raw_instance)
    assert transformed[0]['temp_high_f'] == 0.0
    assert transformed[0]['precipitation_chance'] == 0
    assert transformed[0]['wind_speed_mph'] == 0.0

def test_negative_vals(raw_instance):
    """Test negatives produce successful transformation"""

    raw_instance['daily']['temperature_2m_max'][0] = -500
    raw_instance['daily']['precipitation_probability_max'][0] = -600
    raw_instance['daily']['wind_speed_10m_max'][0] = -601

    transformed = transform.transform_weather(raw_instance)
    assert transformed[0]['temp_high_f'] == -500.0
    assert transformed[0]['precipitation_chance'] == -600
    assert transformed[0]['wind_speed_mph'] == -601.0
# normalize weather data for defined schema
def transform_weather(raw):
    # TODO logging, error handling
    
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

    return forecasts
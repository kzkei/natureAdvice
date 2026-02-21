# natureAdvice
The purpose of this project is to dynamically recommend points of interest, initially national parks in the US, with 14-day forecast data in mind to aid in planning a visit!

This project is: end-to-end data engineering project where Airflow DAG schedules ETL pipeline (Extract: Open-Meteo API -> Transform: data normalization -> Load: postgreSQL). Go REST API serves dynamic park recs by computing visit-quality scores from stored weather forecasts. System processes (initially 63) 63+ locations, refreshes every 6 hours, maintains 14-day forecast horizon

This repository includes:
- the backend system which defines, and handles, available data as well as data to serve to endpoint requests
- an end-to-end Python ETL pipepline that gathers, cleans and loads forecast data every 6 hours with Airflow/DAG integration
- a dynamic location/date recommendaton service with optional region and "top N" limit parameters
- a basic, scalable scoring service that grades locations for ideal visiting conditions

Skills demonstrated: 
- Python ETL pipeline design
- Go REST API design
- Airflow/DAG integration
- PostgreSQL schema & migration design
- Clean Architecture
- Logging & error handling

To use:
- Go, Docker, Python, PostgreSQL
- setup Airflow:
    cd airflow
    docker-compose up -d
    UI at http://localhost:8080

- create .env file at root with contents:
    DB_HOST=localhost
    DB_PORT=5433
    DB_USER=natureadvice
    DB_PASSWORD=your_password
    DB_NAME=natureadvice
    API_PORT=8000

- run migrations in /db/migrations
- pip install -r requirements.txt (in venv if preferred)
- run python etl/main.py
- cd api, go run main.go
- request at:
    curl "http://localhost:8000/api/...": 
        GET /locations      // get all locations
        GET /locations/:name/forecasts      // single location forecast (14 days)
        GET /locations/:name/latest/:date       // single location latest forecast for specified date (within 14 days of current day)
        GET /recommendations/:date?limit=       // main idea, optional limit param (int <= num of available locations)

    Note: date param in url is of 3000-01-31 format


// to add: POST location endpoint implementation, weather history, confidence&volitality, optional region param for rec endpoint, defensive programming wherever see fit 

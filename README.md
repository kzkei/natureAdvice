# natureAdvice ![Tests](https://github.com/kzkei/natureAdvice/actions/workflows/test.yml/badge.svg)

A personal data engineering project that recommends Us national parsks based on upcoming weather conditions 

This project is an end-to-end data engineering project where Airflow DAG schedules ETL pipeline (Extract: Open-Meteo API -> Transform: data normalization -> Load: postgreSQL). Go REST API serves dynamic park recs by computing visit-quality scores from stored weather forecasts. System processes 63 locations, refreshes every 6 hours, maintains 14-day forecast horizon

## Architecture

Open-Metoeo API -> Python ETL -> PostgreSQL -> Go REST API
(Airflow schedules every 6 hours)

## Stack

Pipeline: Python, Airflow
Storage: PostgreSQL
API: Go
Infrastructure: Docker

## Running locally

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

## Endpoints 

- request at:
    curl "http://localhost:8000/api/...": 

        GET /locations                          // get all locations
        GET /locations/:name/forecasts          // single location forecast (14 days)
        GET /locations/:name/latest/:date       // single location latest forecast for specified date (within 14 days of current day)
        GET /recommendations/:date?limit=       // main idea, optional limit param (int <= num of available locations)

    Note: date param in url is of YYYY-MM-DD format

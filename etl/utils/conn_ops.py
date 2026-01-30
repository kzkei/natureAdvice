import psycopg2
from dotenv import load_dotenv
import os

# load dotenv vars
load_dotenv()

# TODO logg, error handling setup

# load Db config from dotenv to connect
DB_CONFIG = {
    'host': os.getenv('DB_HOST', 'localhost'),
    'port': int(os.getenv('DB_PORT', 5433)),
    'user': os.getenv('DB_USER', 'natureadvice'),
    'password': os.getenv('DB_PASSWORD'),
    'database': os.getenv('DB_NAME', 'natureadvice')
}

# validate DB_CONFIG password
if not DB_CONFIG['password']:
    raise ValueError("DB_PASSWORD env var required")

def open():
    conn = psycopg2.connect(**DB_CONFIG)
    cursor = conn.cursor()

    return conn, cursor

def close(conn, cursor):
    cursor.close()
    conn.close()

    return
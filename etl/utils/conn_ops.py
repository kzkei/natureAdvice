import psycopg2
from dotenv import load_dotenv, find_dotenv
import os, logging

# setup logging
logger = logging.getLogger(__name__)

# load dotenv vars using find .env
load_dotenv(find_dotenv())

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

# opens db connection
# returns conn, cursor
def open():
    logger.info("opening db connection")
    conn = psycopg2.connect(**DB_CONFIG)
    cursor = conn.cursor()

    logger.info("db connected")
    return conn, cursor

# closes db connection with defined conn, cursor
def close(conn, cursor):
    logger.info("closing connection")

    try:
        cursor.close()
        conn.close()

        logger.info("connection closed")
        return
    
    except Exception as e:
        logger.error(f"failed to close connection: {e}")
        raise
        